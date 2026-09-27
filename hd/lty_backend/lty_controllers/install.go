package lty_controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"

	"lty_backend/lty_config"
	"lty_backend/lty_utils"
)

type InstallReq struct {
	DBHost    string `json:"db_host"`
	DBPort    string `json:"db_port"`
	DBUser    string `json:"db_user"`
	DBPass    string `json:"db_pass"`
	DBName    string `json:"db_name"`
	AdminUser string `json:"admin_user"`
	AdminPass string `json:"admin_pass"`
}

type CheckDBReq struct {
	DBHost string `json:"db_host"`
	DBPort string `json:"db_port"`
	DBUser string `json:"db_user"`
	DBPass string `json:"db_pass"`
	DBName string `json:"db_name"`
}

func CheckDBConnection(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req CheckDBReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "参数错误", "details": err.Error()})
		return
	}

	if req.DBHost == "" || req.DBPort == "" || req.DBUser == "" || req.DBName == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "缺少必填参数"})
		return
	}

	// 尝试连接到指定的数据库
	dsnWithDB := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", req.DBUser, req.DBPass, req.DBHost, req.DBPort, req.DBName)
	dbTest, err := sql.Open("mysql", dsnWithDB)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "数据库连接失败: " + err.Error()})
		return
	}
	defer dbTest.Close()

	if err := dbTest.Ping(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("无法连接到数据库 '%s'。请检查数据库是否存在且账号密码正确！(原因: %v)", req.DBName, err)})
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "数据库连接成功",
	})
}

func CheckInstall(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"installed": lty_config.IsInstalled,
	})
}

func Install(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if lty_config.IsInstalled {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "系统已经安装过，请勿重复安装"})
		return
	}

	var req InstallReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "参数错误", "details": err.Error()})
		return
	}

	if req.DBHost == "" || req.DBPort == "" || req.DBUser == "" || req.DBName == "" || req.AdminUser == "" || req.AdminPass == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "缺少必填参数"})
		return
	}

	// 1. 尝试直接连接到指定数据库
	dsnWithDB := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", req.DBUser, req.DBPass, req.DBHost, req.DBPort, req.DBName)
	dbTest, err := sql.Open("mysql", dsnWithDB)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "连接数据库失败: " + err.Error()})
		return
	}
	defer dbTest.Close()

	if err := dbTest.Ping(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("无法连接到数据库 '%s'。请确保数据库已存在且账号密码正确！(错误信息: %v)", req.DBName, err)})
		return
	}

	// 2. 保存配置到文件
	cfg := lty_config.AppConfig{
		DBHost: req.DBHost,
		DBPort: req.DBPort,
		DBUser: req.DBUser,
		DBPass: req.DBPass,
		DBName: req.DBName,
	}
	if err := lty_config.SaveConfig(cfg); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "保存配置文件失败: " + err.Error()})
		return
	}

	// 4. 初始化全局数据库连接并执行迁移
	if err := lty_config.InitDB(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "初始化数据库连接失败: " + err.Error()})
		return
	}
	if err := lty_config.AutoMigrate(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "迁移数据表失败: " + err.Error()})
		return
	}

	// 5. 创建或更新初始管理员账号
	hashedPassword := lty_utils.HashPassword(req.AdminPass)

	// 检查是否已经存在管理员
	var firstAdminId int
	err = lty_config.DB.QueryRow("SELECT id FROM lty_admins ORDER BY id ASC LIMIT 1").Scan(&firstAdminId)
	if err == nil {
		// 存在管理员，更新第一个管理员的信息
		updateQuery := "UPDATE lty_admins SET username = ?, password = ? WHERE id = ?"
		if _, err := lty_config.DB.Exec(updateQuery, req.AdminUser, hashedPassword, firstAdminId); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "更新管理员账号失败: " + err.Error()})
			return
		}
	} else {
		// 不存在管理员，插入新的管理员
		insertQuery := "INSERT INTO lty_admins (username, password) VALUES (?, ?)"
		if _, err := lty_config.DB.Exec(insertQuery, req.AdminUser, hashedPassword); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "创建管理员账号失败: " + err.Error()})
			return
		}
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "安装成功",
	})
}

package lty_controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"lty_backend/lty_config"
	"lty_backend/lty_models"
	"lty_backend/lty_utils"
)

type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type CreateAdminReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UpdateAdminPasswordReq struct {
	ID          uint   `json:"id"`
	NewPassword string `json:"new_password"`
}

// UpdatePassword 修改管理员密码
func UpdatePassword(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req UpdateAdminPasswordReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "参数错误"})
		return
	}

	if req.ID == 0 || req.NewPassword == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "ID和新密码不能为空"})
		return
	}

	// 可选：检查要修改的管理员是否存在
	var count int
	err := lty_config.DB.QueryRow("SELECT COUNT(*) FROM lty_admins WHERE id = ?", req.ID).Scan(&count)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "数据库查询错误"})
		return
	}
	if count == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "要修改的管理员不存在"})
		return
	}

	hashedPassword := lty_utils.HashPassword(req.NewPassword)
	updateQuery := "UPDATE lty_admins SET password = ? WHERE id = ?"
	if _, err := lty_config.DB.Exec(updateQuery, hashedPassword, req.ID); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "修改密码失败"})
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "密码修改成功",
	})
}

// CreateAdmin 新增管理员
func CreateAdmin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req CreateAdminReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "参数错误"})
		return
	}

	if req.Username == "" || req.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "用户名和密码不能为空"})
		return
	}

	// 检查用户名是否已存在
	var count int
	err := lty_config.DB.QueryRow("SELECT COUNT(*) FROM lty_admins WHERE username = ?", req.Username).Scan(&count)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "数据库查询错误"})
		return
	}
	if count > 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "用户名已存在"})
		return
	}

	hashedPassword := lty_utils.HashPassword(req.Password)
	insertQuery := "INSERT INTO lty_admins (username, password) VALUES (?, ?)"
	result, err := lty_config.DB.Exec(insertQuery, req.Username, hashedPassword)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "创建管理员失败"})
		return
	}

	id, _ := result.LastInsertId()
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "管理员创建成功",
		"id":      id,
	})
}

// GetAdmins 获取管理员列表
func GetAdmins(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	query := "SELECT id, username, created_at, updated_at FROM lty_admins ORDER BY id DESC"
	rows, err := lty_config.DB.Query(query)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "数据库查询错误"})
		return
	}
	defer rows.Close()

	var admins []map[string]interface{}
	for rows.Next() {
		var admin lty_models.Admin
		if err := rows.Scan(&admin.ID, &admin.Username, &admin.CreatedAt, &admin.UpdatedAt); err != nil {
			continue
		}

		admins = append(admins, map[string]interface{}{
			"id":         admin.ID,
			"username":   admin.Username,
			"created_at": admin.CreatedAt,
			"updated_at": admin.UpdatedAt,
		})
	}

	if admins == nil {
		admins = make([]map[string]interface{}, 0)
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"list": admins,
	})
}

// DeleteAdmin 删除管理员
func DeleteAdmin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req struct {
		ID uint `json:"id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "参数错误"})
		return
	}

	// 不允许删除当前登录的管理员
	adminIDVal := r.Context().Value("admin_id")
	if adminIDVal != nil && adminIDVal.(uint) == req.ID {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "不能删除当前登录的管理员账号"})
		return
	}

	query := "DELETE FROM lty_admins WHERE id = ?"
	_, err := lty_config.DB.Exec(query, req.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "删除管理员失败"})
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "管理员删除成功",
	})
}

// Login 管理员登录
func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req LoginReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "参数错误"})
		return
	}

	if req.Username == "" || req.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "参数错误"})
		return
	}

	var admin lty_models.Admin
	query := "SELECT id, username, password, created_at, updated_at FROM lty_admins WHERE username = ?"
	err := lty_config.DB.QueryRow(query, req.Username).Scan(&admin.ID, &admin.Username, &admin.Password, &admin.CreatedAt, &admin.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "用户名或密码错误"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "数据库查询错误"})
		return
	}

	if admin.Password != lty_utils.HashPassword(req.Password) {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "用户名或密码错误"})
		return
	}

	token, err := lty_utils.GenerateToken(admin.ID, admin.Username, "admin")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "生成Token失败"})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"token": token,
		"admin": map[string]interface{}{
			"id":       admin.ID,
			"username": admin.Username,
		},
	})
}

// GetUserInfo 获取当前登录管理员信息
func GetUserInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	adminIDVal := r.Context().Value("admin_id")
	if adminIDVal == nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "未授权"})
		return
	}

	adminID := adminIDVal.(uint)

	var admin lty_models.Admin
	query := "SELECT id, username, created_at, updated_at FROM lty_admins WHERE id = ?"
	err := lty_config.DB.QueryRow(query, adminID).Scan(&admin.ID, &admin.Username, &admin.CreatedAt, &admin.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "用户不存在"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "数据库查询错误"})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":         admin.ID,
		"username":   admin.Username,
		"created_at": admin.CreatedAt,
	})
}

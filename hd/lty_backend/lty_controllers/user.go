package lty_controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"lty_backend/lty_config"
	"lty_backend/lty_models"
	"lty_backend/lty_utils"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
var phoneRegex = regexp.MustCompile(`^1[3-9]\d{9}$`)

type CreateUserReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Status   int     `json:"status"`
	Points   float64 `json:"points"`
}

type UserLoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserRegisterReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateUserReq struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Status   int     `json:"status"`
	Points   float64 `json:"points"`
	Title    string  `json:"title"`
	Reason   string  `json:"reason"`
}

// CreateUser 新增用户
func CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req CreateUserReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "参数错误"})
		return
	}

	if req.Password == "" || req.Email == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "密码和邮箱不能为空"})
		return
	}

	if !emailRegex.MatchString(req.Email) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "邮箱格式不正确"})
		return
	}

	if req.Phone != "" && !phoneRegex.MatchString(req.Phone) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "手机号格式不正确"})
		return
	}

	var count int
	if req.Username != "" {
		err := lty_config.DB.QueryRow("SELECT COUNT(*) FROM lty_users WHERE username = ?", req.Username).Scan(&count)
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
	}

	err := lty_config.DB.QueryRow("SELECT COUNT(*) FROM lty_users WHERE email = ?", req.Email).Scan(&count)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "数据库查询错误"})
		return
	}
	if count > 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "邮箱已存在"})
		return
	}

	if req.Phone != "" {
		err = lty_config.DB.QueryRow("SELECT COUNT(*) FROM lty_users WHERE phone = ?", req.Phone).Scan(&count)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "数据库查询错误"})
			return
		}
		if count > 0 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "手机号已存在"})
			return
		}
	}

	hashedPassword := lty_utils.HashPassword(req.Password)
	username := req.Username
	if username == "" {
		username = req.Email
	}
	insertQuery := "INSERT INTO lty_users (username, password, email, phone, status, points) VALUES (?, ?, ?, ?, ?, ?)"
	result, err := lty_config.DB.Exec(insertQuery, username, hashedPassword, req.Email, req.Phone, req.Status, req.Points)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "创建用户失败"})
		return
	}

	id, _ := result.LastInsertId()
	
	if req.Points != 0 {
		lty_config.DB.Exec("INSERT INTO lty_point_records (user_id, points_change, title, reason) VALUES (?, ?, ?, ?)", id, req.Points, "官方发放", "管理员添加用户初始积分")
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "用户创建成功",
		"id":      id,
	})
}

// GetUsers 获取用户列表
func GetUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("page_size")
	search := r.URL.Query().Get("search")

	page := 1
	if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
		page = p
	}

	pageSize := 10
	if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 {
		pageSize = ps
	}

	offset := (page - 1) * pageSize

	whereClause := "1=1"
	var args []interface{}

	if search != "" {
		whereClause += " AND (username LIKE ? OR email LIKE ? OR phone LIKE ?)"
		searchParam := "%" + search + "%"
		args = append(args, searchParam, searchParam, searchParam)
	}

	var total int
	countQuery := "SELECT COUNT(*) FROM lty_users WHERE " + whereClause
	if err := lty_config.DB.QueryRow(countQuery, args...).Scan(&total); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "数据库查询错误"})
		return
	}

	query := "SELECT id, username, email, phone, status, points, created_at, updated_at FROM lty_users WHERE " + whereClause + " ORDER BY id DESC LIMIT ? OFFSET ?"
	args = append(args, pageSize, offset)

	rows, err := lty_config.DB.Query(query, args...)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "数据库查询错误"})
		return
	}
	defer rows.Close()

	var users []map[string]interface{}
	for rows.Next() {
		var user lty_models.User
		if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Phone, &user.Status, &user.Points, &user.CreatedAt, &user.UpdatedAt); err != nil {
			continue
		}

		users = append(users, map[string]interface{}{
			"id":         user.ID,
			"username":   user.Username,
			"email":      user.Email,
			"phone":      user.Phone,
			"status":     user.Status,
			"points":     user.Points,
			"created_at": user.CreatedAt,
			"updated_at": user.UpdatedAt,
		})
	}

	if users == nil {
		users = make([]map[string]interface{}, 0)
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"list":  users,
		"total": total,
	})
}

// UpdateUser 更新用户信息
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req UpdateUserReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "参数错误"})
		return
	}

	if req.ID == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "ID不能为空"})
		return
	}

	if req.Email == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "邮箱不能为空"})
		return
	}

	if !emailRegex.MatchString(req.Email) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "邮箱格式不正确"})
		return
	}

	if req.Phone != "" && !phoneRegex.MatchString(req.Phone) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "手机号格式不正确"})
		return
	}

	var count int
	var currentPoints float64
	err := lty_config.DB.QueryRow("SELECT COUNT(*), IFNULL(points, 0) FROM lty_users WHERE id = ?", req.ID).Scan(&count, &currentPoints)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "数据库查询错误"})
		return
	}
	if count == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "用户不存在"})
		return
	}

	err = lty_config.DB.QueryRow("SELECT COUNT(*) FROM lty_users WHERE email = ? AND id != ?", req.Email, req.ID).Scan(&count)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "数据库查询错误"})
		return
	}
	if count > 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "邮箱已存在"})
		return
	}

	if req.Phone != "" {
		err = lty_config.DB.QueryRow("SELECT COUNT(*) FROM lty_users WHERE phone = ? AND id != ?", req.Phone, req.ID).Scan(&count)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "数据库查询错误"})
			return
		}
		if count > 0 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "手机号已存在"})
			return
		}
	}

	if req.Password != "" {
		hashedPassword := lty_utils.HashPassword(req.Password)
		updateQuery := "UPDATE lty_users SET username = ?, password = ?, email = ?, phone = ?, status = ?, points = ? WHERE id = ?"
		if _, err := lty_config.DB.Exec(updateQuery, req.Username, hashedPassword, req.Email, req.Phone, req.Status, req.Points, req.ID); err != nil {
			fmt.Println("Error updating user (with password):", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "更新用户信息失败"})
			return
		}
	} else {
		updateQuery := "UPDATE lty_users SET username = ?, email = ?, phone = ?, status = ?, points = ? WHERE id = ?"
		if _, err := lty_config.DB.Exec(updateQuery, req.Username, req.Email, req.Phone, req.Status, req.Points, req.ID); err != nil {
			fmt.Println("Error updating user (without password):", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "更新用户信息失败"})
			return
		}
	}

	if req.Points != currentPoints {
		pointsChange := req.Points - currentPoints
		reason := req.Reason
		if reason == "" {
			reason = "管理员修改"
		}
		title := req.Title
		if title == "" {
			title = "官方发放"
		}
		lty_config.DB.Exec("INSERT INTO lty_point_records (user_id, points_change, title, reason) VALUES (?, ?, ?, ?)", req.ID, pointsChange, title, reason)
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "用户信息更新成功",
	})
}

// DeleteUser 删除用户
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req struct {
		ID uint `json:"id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "参数错误"})
		return
	}

	query := "DELETE FROM lty_users WHERE id = ?"
	_, err := lty_config.DB.Exec(query, req.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "删除用户失败"})
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "用户删除成功",
	})
}

// UserLogin 普通用户登录
func UserLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req UserLoginReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "参数错误"})
		return
	}

	if req.Email == "" || req.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "参数错误"})
		return
	}

	var user lty_models.User
	query := "SELECT id, username, email, password, status, points FROM lty_users WHERE email = ?"
	err := lty_config.DB.QueryRow(query, req.Email).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Status, &user.Points)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "邮箱或密码错误"})
		return
	}

	if user.Status == 0 {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]string{"error": "账号已被禁用"})
		return
	}

	if user.Password != lty_utils.HashPassword(req.Password) {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "邮箱或密码错误"})
		return
	}

	// 使用 Email 作为 Username 字段存储在 Token 中，方便兼容已有逻辑
	usernameToStore := user.Username
	if usernameToStore == "" {
		usernameToStore = user.Email
	}
	token, err := lty_utils.GenerateToken(user.ID, usernameToStore, "user")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "生成Token失败"})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"token": token,
		"user": map[string]interface{}{
			"id":       user.ID,
			"username": usernameToStore,
			"email":    user.Email,
			"points":   user.Points,
		},
	})
}

// GetUserInfoForFrontend 获取前端当前登录用户信息
func GetUserInfoForFrontend(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userID := r.Context().Value("user_id").(uint)

	var user lty_models.User
	query := "SELECT id, username, email, phone, status, points FROM lty_users WHERE id = ?"
	err := lty_config.DB.QueryRow(query, userID).Scan(&user.ID, &user.Username, &user.Email, &user.Phone, &user.Status, &user.Points)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "用户不存在"})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
		"points":   user.Points,
	})
}

type GetPointRecordsReq struct {
	UserID   uint `json:"user_id"` // For admin, to query specific user's points
	Page     int  `json:"page"`
	PageSize int  `json:"page_size"`
}

type PointRecord struct {
	ID           uint      `json:"id"`
	UserID       uint      `json:"user_id"`
	PointsChange float64   `json:"points_change"`
	Title        string    `json:"title"`
	Reason       string    `json:"reason"`
	CreatedAt    time.Time `json:"created_at"`
}

// GetPointRecords 获取积分记录
func GetPointRecords(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var userID uint
	isAdmin := false

	// 先检查是否是管理员
	adminIDVal := r.Context().Value("admin_id")
	if adminIDVal != nil {
		isAdmin = true
	} else {
		// 不是管理员，检查是否是普通用户
		userIDVal := r.Context().Value("user_id")
		if userIDVal == nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "未授权"})
			return
		}
		userID = userIDVal.(uint)
	}

	queryUserID := r.URL.Query().Get("user_id")
	if queryUserID != "" && isAdmin {
		if uid, err := strconv.Atoi(queryUserID); err == nil {
			userID = uint(uid)
		}
	} else if isAdmin && queryUserID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "管理员查询需指定用户ID"})
		return
	}

	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("page_size")

	page := 1
	if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
		page = p
	}

	pageSize := 20
	if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 {
		pageSize = ps
	}

	offset := (page - 1) * pageSize

	var total int
	err := lty_config.DB.QueryRow("SELECT COUNT(*) FROM lty_point_records WHERE user_id = ?", userID).Scan(&total)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "数据库查询错误"})
		return
	}

	query := "SELECT id, user_id, points_change, IFNULL(title, ''), reason, created_at FROM lty_point_records WHERE user_id = ? ORDER BY id DESC LIMIT ? OFFSET ?"
	rows, err := lty_config.DB.Query(query, userID, pageSize, offset)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "数据库查询错误"})
		return
	}
	defer rows.Close()

	var records []PointRecord
	for rows.Next() {
		var record PointRecord
		if err := rows.Scan(&record.ID, &record.UserID, &record.PointsChange, &record.Title, &record.Reason, &record.CreatedAt); err != nil {
			continue
		}
		records = append(records, record)
	}

	if records == nil {
		records = make([]PointRecord, 0)
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"list":  records,
		"total": total,
	})
}

func UserRegister(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req UserRegisterReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "参数错误"})
		return
	}

	if req.Email == "" || req.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "邮箱和密码不能为空"})
		return
	}

	if !emailRegex.MatchString(req.Email) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "邮箱格式不正确"})
		return
	}

	var count int
	err := lty_config.DB.QueryRow("SELECT COUNT(*) FROM lty_users WHERE email = ?", req.Email).Scan(&count)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "数据库查询错误"})
		return
	}
	if count > 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "邮箱已存在"})
		return
	}

	hashedPassword := lty_utils.HashPassword(req.Password)
	// 如果前端不传username，默认使用邮箱作为username
	username := req.Email
	
	// 获取新用户默认积分配置
	defaultPoints := 5.0
	setting, err := lty_models.GetSetting("default_register_points")
	if err == nil && setting != nil && setting.KeyValue != "" {
		if val, err := strconv.ParseFloat(setting.KeyValue, 64); err == nil {
			defaultPoints = val
		}
	}

	insertQuery := "INSERT INTO lty_users (username, password, email, phone, status, points) VALUES (?, ?, ?, ?, ?, ?)"
	result, err := lty_config.DB.Exec(insertQuery, username, hashedPassword, req.Email, "", 1, defaultPoints) // 默认启用状态为 1, 使用配置积分
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "注册失败"})
		return
	}

	id, _ := result.LastInsertId()
	
	if defaultPoints > 0 {
		lty_config.DB.Exec("INSERT INTO lty_point_records (user_id, points_change, title, reason) VALUES (?, ?, ?, ?)", id, defaultPoints, "官方发放", "新用户注册赠送")
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "注册成功",
		"id":      id,
	})
}

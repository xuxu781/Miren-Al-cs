package lty_controllers

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"lty_backend/lty_config"
	"lty_backend/lty_models"
)

type GenerateCdkeyReq struct {
	Count     int        `json:"count"`
	Points    float64    `json:"points"`
	MaxUses   int        `json:"max_uses"`
	ExpiresAt *time.Time `json:"expires_at"`
	Prefix    string     `json:"prefix"`
}

func generateRandomCdkey(length int) string {
	b := make([]byte, length)
	rand.Read(b)
	return fmt.Sprintf("%x", b)[:length]
}

func GenerateCdkeys(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req GenerateCdkeyReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "参数错误"})
		return
	}

	if req.Count <= 0 || req.Count > 1000 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "生成数量必须在1-1000之间"})
		return
	}

	if req.Points <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "积分必须大于0"})
		return
	}

	maxUses := req.MaxUses
	if maxUses < 0 {
		maxUses = 1
	}

	prefix := "Mirenai-"
	if req.Prefix != "" {
		prefix = req.Prefix + "-"
	}

	successCount := 0
	var generatedKeys []string
	for i := 0; i < req.Count; i++ {
		key := prefix + generateRandomCdkey(16)
		_, err := lty_config.DB.Exec("INSERT INTO lty_cdkeys (cdkey, points, status, max_uses, current_uses, expires_at) VALUES (?, ?, 0, ?, 0, ?)", key, req.Points, maxUses, req.ExpiresAt)
		if err == nil {
			successCount++
			generatedKeys = append(generatedKeys, key)
		}
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": fmt.Sprintf("成功生成 %d 个卡密", successCount),
		"keys":    generatedKeys,
	})
}

func GetCdkeys(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("page_size")
	search := r.URL.Query().Get("search")
	statusStr := r.URL.Query().Get("status")
	userEmailStr := r.URL.Query().Get("user_email")
	pointsStr := r.URL.Query().Get("points")

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
		whereClause += " AND c.cdkey LIKE ?"
		args = append(args, "%"+search+"%")
	}
	
	if statusStr != "" {
		if status, err := strconv.Atoi(statusStr); err == nil {
			if status == 3 { // 已过期
				whereClause += " AND c.status = 0 AND c.expires_at IS NOT NULL AND c.expires_at < NOW()"
			} else if status == 0 { // 未使用（且未过期）
				whereClause += " AND c.status = 0 AND (c.expires_at IS NULL OR c.expires_at >= NOW())"
			} else {
				whereClause += " AND c.status = ?"
				args = append(args, status)
			}
		}
	}

	if userEmailStr != "" {
		whereClause += ` AND (
			c.used_by_user_id IN (SELECT id FROM lty_users WHERE email = ?) 
			OR c.id IN (
				SELECT cdkey_id FROM lty_cdkey_usages cu 
				JOIN lty_users u ON cu.user_id = u.id 
				WHERE u.email = ?
			)
		)`
		args = append(args, userEmailStr, userEmailStr)
	}

	if pointsStr != "" {
		if points, err := strconv.ParseFloat(pointsStr, 64); err == nil {
			whereClause += " AND c.points = ?"
			args = append(args, points)
		}
	}

	var total int
	if err := lty_config.DB.QueryRow("SELECT COUNT(*) FROM lty_cdkeys c WHERE "+whereClause, args...).Scan(&total); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "数据库查询错误"})
		return
	}

	isExport := r.URL.Query().Get("export") == "1"

	var query string
	if isExport {
		query = `
			SELECT c.id, c.cdkey, c.points, c.status, c.max_uses, c.current_uses, 
			       c.used_by_user_id, u.email as used_by_user_email,
			       c.expires_at, c.created_at, c.updated_at, c.used_at 
			FROM lty_cdkeys c
			LEFT JOIN lty_users u ON c.used_by_user_id = u.id
			WHERE ` + whereClause + `
			ORDER BY c.id DESC
		`
	} else {
		query = `
			SELECT c.id, c.cdkey, c.points, c.status, c.max_uses, c.current_uses, 
			       c.used_by_user_id, u.email as used_by_user_email,
			       c.expires_at, c.created_at, c.updated_at, c.used_at 
			FROM lty_cdkeys c
			LEFT JOIN lty_users u ON c.used_by_user_id = u.id
			WHERE ` + whereClause + `
			ORDER BY c.id DESC LIMIT ? OFFSET ?
		`
		args = append(args, pageSize, offset)
	}

	rows, err := lty_config.DB.Query(query, args...)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "数据库查询错误"})
		return
	}
	defer rows.Close()

	var cdkeys []map[string]interface{}
	for rows.Next() {
		var c lty_models.Cdkey
		var usedBy *uint
		var usedByEmail *string
		var usedAt *time.Time
		var expiresAt *time.Time
		if err := rows.Scan(&c.ID, &c.Cdkey, &c.Points, &c.Status, &c.MaxUses, &c.CurrentUses, &usedBy, &usedByEmail, &expiresAt, &c.CreatedAt, &c.UpdatedAt, &usedAt); err != nil {
			continue
		}
		c.UsedByUserID = usedBy
		c.ExpiresAt = expiresAt
		c.UsedAt = usedAt
		
		cdkeyMap := map[string]interface{}{
			"id": c.ID,
			"cdkey": c.Cdkey,
			"points": c.Points,
			"status": c.Status,
			"max_uses": c.MaxUses,
			"current_uses": c.CurrentUses,
			"used_by_user_id": c.UsedByUserID,
			"used_by_user_email": usedByEmail,
			"expires_at": c.ExpiresAt,
			"created_at": c.CreatedAt,
			"updated_at": c.UpdatedAt,
			"used_at": c.UsedAt,
		}
		cdkeys = append(cdkeys, cdkeyMap)
	}

	if cdkeys == nil {
		cdkeys = make([]map[string]interface{}, 0)
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"list":  cdkeys,
		"total": total,
	})
}

func GetCdkeyUsages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	cdkeyIDStr := r.URL.Query().Get("cdkey_id")
	if cdkeyIDStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "参数错误"})
		return
	}

	cdkeyID, err := strconv.Atoi(cdkeyIDStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "参数错误"})
		return
	}

	query := `
		SELECT u.email, cu.created_at
		FROM lty_cdkey_usages cu
		JOIN lty_users u ON cu.user_id = u.id
		WHERE cu.cdkey_id = ?
		ORDER BY cu.created_at DESC
	`

	rows, err := lty_config.DB.Query(query, cdkeyID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "数据库查询错误"})
		return
	}
	defer rows.Close()

	var usages []map[string]interface{}
	for rows.Next() {
		var email string
		var createdAt time.Time
		if err := rows.Scan(&email, &createdAt); err != nil {
			continue
		}
		usages = append(usages, map[string]interface{}{
			"email":      email,
			"created_at": createdAt,
		})
	}

	if usages == nil {
		usages = make([]map[string]interface{}, 0)
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"list": usages,
	})
}

func DeleteCdkey(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req struct {
		ID uint `json:"id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "参数错误"})
		return
	}

	_, err := lty_config.DB.Exec("DELETE FROM lty_cdkeys WHERE id = ?", req.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "删除失败"})
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "删除成功"})
}

type BatchDeleteCdkeyReq struct {
	IDs []uint `json:"ids"`
}

func BatchDeleteCdkey(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req BatchDeleteCdkeyReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "参数错误"})
		return
	}

	if len(req.IDs) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "未选择任何记录"})
		return
	}

	query := "DELETE FROM lty_cdkeys WHERE id IN ("
	args := make([]interface{}, len(req.IDs))
	for i, id := range req.IDs {
		if i > 0 {
			query += ", "
		}
		query += "?"
		args[i] = id
	}
	query += ")"

	_, err := lty_config.DB.Exec(query, args...)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "批量删除失败"})
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "批量删除成功"})
}

type UpdateCdkeyReq struct {
	ID        uint       `json:"id"`
	Points    float64    `json:"points"`
	MaxUses   int        `json:"max_uses"`
	ExpiresAt *time.Time `json:"expires_at"`
}

func UpdateCdkey(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req UpdateCdkeyReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "参数错误"})
		return
	}

	if req.Points <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "积分必须大于0"})
		return
	}

	maxUses := req.MaxUses
	if maxUses < 0 {
		maxUses = 1
	}

	tx, err := lty_config.DB.Begin()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "服务器内部错误"})
		return
	}
	defer tx.Rollback()

	var currentUses, status int
	var usedByUserID *uint
	err = tx.QueryRow("SELECT current_uses, status, used_by_user_id FROM lty_cdkeys WHERE id = ? FOR UPDATE", req.ID).Scan(&currentUses, &status, &usedByUserID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "找不到该卡密"})
		return
	}

	// 如果当前不是已作废状态，则根据 max_uses 重新计算状态
	if status != 2 {
		if maxUses > 0 && currentUses >= maxUses {
			status = 1 // 已耗尽
		} else {
			status = 0 // 可用
		}
	}

	// 执行更新
	if status == 0 && maxUses != 1 {
		// 如果重置为了可用状态，并且变成了多用户（maxUses!=1），
		// 可以选择将 used_by_user_id 置空或者保留（这里我们选择保留原样以作记录），
		// 最重要的是 status 变为 0 后，其他用户通过 /api/user/cdkeys/use 兑换时，
		// 就不再被 `if c.Status != 0` 拦截。
	}

	_, err = tx.Exec("UPDATE lty_cdkeys SET points = ?, max_uses = ?, expires_at = ?, status = ? WHERE id = ?", req.Points, maxUses, req.ExpiresAt, status, req.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "更新失败"})
		return
	}

	if err := tx.Commit(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "事务提交失败"})
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "更新成功"})
}

func VoidCdkey(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req struct {
		ID uint `json:"id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "参数错误"})
		return
	}

	_, err := lty_config.DB.Exec("UPDATE lty_cdkeys SET status = 2 WHERE id = ?", req.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "作废失败"})
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "作废成功"})
}

func EnableCdkey(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req struct {
		ID uint `json:"id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "参数错误"})
		return
	}

	tx, err := lty_config.DB.Begin()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "服务器内部错误"})
		return
	}
	defer tx.Rollback()

	var currentUses, maxUses, status int
	err = tx.QueryRow("SELECT current_uses, max_uses, status FROM lty_cdkeys WHERE id = ? FOR UPDATE", req.ID).Scan(&currentUses, &maxUses, &status)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "找不到该卡密"})
		return
	}

	if status != 2 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "该卡密不是作废状态"})
		return
	}

	// 恢复状态：如果耗尽则为 1，否则为 0
	if maxUses > 0 && currentUses >= maxUses {
		status = 1
	} else {
		status = 0
	}

	_, err = tx.Exec("UPDATE lty_cdkeys SET status = ? WHERE id = ?", status, req.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "启用失败"})
		return
	}

	if err := tx.Commit(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "事务提交失败"})
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "启用成功"})
}

type UseCdkeyReq struct {
	Cdkey string `json:"cdkey"`
}

func UseCdkey(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	userIDVal := r.Context().Value("user_id")
	if userIDVal == nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "未登录"})
		return
	}
	userID := userIDVal.(uint)

	var req UseCdkeyReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "参数错误"})
		return
	}

	if req.Cdkey == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "卡密不能为空"})
		return
	}

	tx, err := lty_config.DB.Begin()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "服务器内部错误"})
		return
	}
	defer tx.Rollback()

	var c lty_models.Cdkey
	var expiresAt *time.Time
	err = tx.QueryRow("SELECT id, points, status, max_uses, current_uses, expires_at FROM lty_cdkeys WHERE cdkey = ? FOR UPDATE", req.Cdkey).Scan(&c.ID, &c.Points, &c.Status, &c.MaxUses, &c.CurrentUses, &expiresAt)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "无效的卡密"})
		return
	}
	c.ExpiresAt = expiresAt

	if c.Status != 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "卡密已被使用或已耗尽"})
		return
	}

	if c.ExpiresAt != nil && time.Now().After(*c.ExpiresAt) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "卡密已过期"})
		return
	}

	var usageCount int
	err = tx.QueryRow("SELECT COUNT(*) FROM lty_cdkey_usages WHERE cdkey_id = ? AND user_id = ?", c.ID, userID).Scan(&usageCount)
	if usageCount > 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "您已经使用过该卡密"})
		return
	}

	c.CurrentUses++
	if c.MaxUses > 0 && c.CurrentUses >= c.MaxUses {
		c.Status = 1
	}

	// Update cdkey status (记录最后使用人，状态根据耗尽情况决定)
	_, err = tx.Exec("UPDATE lty_cdkeys SET status = ?, current_uses = ?, used_by_user_id = ?, used_at = ? WHERE id = ?", c.Status, c.CurrentUses, userID, time.Now(), c.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "使用卡密失败"})
		return
	}

	// Insert usage record
	_, err = tx.Exec("INSERT INTO lty_cdkey_usages (cdkey_id, user_id) VALUES (?, ?)", c.ID, userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "记录使用历史失败"})
		return
	}

	// Update user points
	_, err = tx.Exec("UPDATE lty_users SET points = points + ? WHERE id = ?", c.Points, userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "增加积分失败"})
		return
	}

	// Add point record
	reason := fmt.Sprintf("使用卡密: %s", req.Cdkey)
	_, err = tx.Exec("INSERT INTO lty_point_records (user_id, points_change, title, reason) VALUES (?, ?, ?, ?)", userID, c.Points, "卡密兑换", reason)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "记录积分日志失败"})
		return
	}

	if err := tx.Commit(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "事务提交失败"})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": fmt.Sprintf("成功兑换 %.2f 积分", c.Points),
		"points":  c.Points,
	})
}

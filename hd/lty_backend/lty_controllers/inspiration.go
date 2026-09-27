package lty_controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"lty_backend/lty_config"
	"lty_backend/lty_models"
)

type CreateInspirationReq struct {
	Content      string `json:"content"`
	ImageURL     string `json:"image_url"`
	MainCategory string `json:"main_category"`
	SubCategory  string `json:"sub_category"`
}

func CreateInspiration(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req CreateInspirationReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "参数错误"})
		return
	}

	if req.Content == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "内容不能为空"})
		return
	}

	_, err := lty_config.DB.Exec("INSERT INTO lty_inspirations (content, image_url, main_category, sub_category) VALUES (?, ?, ?, ?)", req.Content, req.ImageURL, req.MainCategory, req.SubCategory)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "创建失败: " + err.Error()})
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "创建成功"})
}

type UpdateInspirationReq struct {
	ID           uint   `json:"id"`
	Content      string `json:"content"`
	ImageURL     string `json:"image_url"`
	MainCategory string `json:"main_category"`
	SubCategory  string `json:"sub_category"`
}

func UpdateInspiration(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req UpdateInspirationReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "参数错误"})
		return
	}

	if req.Content == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "内容不能为空"})
		return
	}

	_, err := lty_config.DB.Exec("UPDATE lty_inspirations SET content = ?, image_url = ?, main_category = ?, sub_category = ? WHERE id = ?", req.Content, req.ImageURL, req.MainCategory, req.SubCategory, req.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "更新失败: " + err.Error()})
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "更新成功"})
}

func DeleteInspiration(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req struct {
		ID uint `json:"id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "参数错误"})
		return
	}

	_, err := lty_config.DB.Exec("DELETE FROM lty_inspirations WHERE id = ?", req.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "删除失败"})
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "删除成功"})
}

type BatchDeleteInspirationReq struct {
	IDs []uint `json:"ids"`
}

func BatchDeleteInspiration(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req BatchDeleteInspirationReq
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

	query := "DELETE FROM lty_inspirations WHERE id IN ("
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

func GetInspirations(w http.ResponseWriter, r *http.Request) {
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
		whereClause += " AND content LIKE ?"
		args = append(args, "%"+search+"%")
	}

	var total int
	if err := lty_config.DB.QueryRow("SELECT COUNT(*) FROM lty_inspirations WHERE "+whereClause, args...).Scan(&total); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "数据库查询错误: " + err.Error()})
		return
	}

	query := "SELECT id, content, IFNULL(image_url, ''), IFNULL(main_category, ''), IFNULL(sub_category, ''), created_at, updated_at FROM lty_inspirations WHERE " + whereClause + " ORDER BY id DESC LIMIT ? OFFSET ?"
	args = append(args, pageSize, offset)

	rows, err := lty_config.DB.Query(query, args...)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "数据库查询错误: " + err.Error()})
		return
	}
	defer rows.Close()

	var list []lty_models.Inspiration
	for rows.Next() {
		var item lty_models.Inspiration
		if err := rows.Scan(&item.ID, &item.Content, &item.ImageURL, &item.MainCategory, &item.SubCategory, &item.CreatedAt, &item.UpdatedAt); err != nil {
			continue
		}
		list = append(list, item)
	}

	if list == nil {
		list = make([]lty_models.Inspiration, 0)
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"list":  list,
		"total": total,
	})
}

// GetPublicInspirations is used for the frontend without auth
func GetPublicInspirations(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	mainCategory := r.URL.Query().Get("main_category")
	subCategory := r.URL.Query().Get("sub_category")

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

	whereClause := "1=1"
	var args []interface{}

	if mainCategory != "" {
		whereClause += " AND main_category = ?"
		args = append(args, mainCategory)
	}
	if subCategory != "" {
		whereClause += " AND sub_category = ?"
		args = append(args, subCategory)
	}

	var total int
	if err := lty_config.DB.QueryRow("SELECT COUNT(*) FROM lty_inspirations WHERE "+whereClause, args...).Scan(&total); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "数据库查询错误: " + err.Error()})
		return
	}

	query := "SELECT id, content, IFNULL(image_url, ''), IFNULL(main_category, ''), IFNULL(sub_category, ''), created_at, updated_at FROM lty_inspirations WHERE " + whereClause + " ORDER BY id DESC LIMIT ? OFFSET ?"
	args = append(args, pageSize, offset)
	
	rows, err := lty_config.DB.Query(query, args...)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "数据库查询错误: " + err.Error()})
		return
	}
	defer rows.Close()

	var list []lty_models.Inspiration
	for rows.Next() {
		var item lty_models.Inspiration
		if err := rows.Scan(&item.ID, &item.Content, &item.ImageURL, &item.MainCategory, &item.SubCategory, &item.CreatedAt, &item.UpdatedAt); err != nil {
			continue
		}
		list = append(list, item)
	}

	if list == nil {
		list = make([]lty_models.Inspiration, 0)
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"list": list,
		"total": total,
		"page": page,
		"page_size": pageSize,
	})
}

func GetPublicInspirationCategories(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var keyValue string
	err := lty_config.DB.QueryRow("SELECT key_value FROM lty_settings WHERE key_name = ?", "inspiration_categories").Scan(&keyValue)
	
	if err != nil || keyValue == "" {
		// Default fallback if not set
		json.NewEncoder(w).Encode(map[string]interface{}{
			"categories": []map[string]interface{}{
				{"name": "人物", "sub": []string{"其他"}},
				{"name": "风景", "sub": []string{"其他"}},
				{"name": "建筑", "sub": []string{"其他"}},
				{"name": "科幻", "sub": []string{"其他"}},
				{"name": "二次元", "sub": []string{"其他"}},
				{"name": "其他", "sub": []string{"其他"}},
			},
		})
		return
	}

	var parsedCategories interface{}
	if err := json.Unmarshal([]byte(keyValue), &parsedCategories); err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"categories": []interface{}{},
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"categories": parsedCategories,
	})
}

func GetInspirationCategories(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	mainRows, err := lty_config.DB.Query("SELECT DISTINCT main_category FROM lty_inspirations WHERE main_category != '' AND main_category IS NOT NULL")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "数据库查询错误: " + err.Error()})
		return
	}
	defer mainRows.Close()

	var mainCategories []string
	for mainRows.Next() {
		var cat string
		if err := mainRows.Scan(&cat); err == nil {
			mainCategories = append(mainCategories, cat)
		}
	}

	subRows, err := lty_config.DB.Query("SELECT DISTINCT sub_category FROM lty_inspirations WHERE sub_category != '' AND sub_category IS NOT NULL")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "数据库查询错误: " + err.Error()})
		return
	}
	defer subRows.Close()

	var subCategories []string
	for subRows.Next() {
		var cat string
		if err := subRows.Scan(&cat); err == nil {
			subCategories = append(subCategories, cat)
		}
	}

	if mainCategories == nil {
		mainCategories = []string{}
	}
	if subCategories == nil {
		subCategories = []string{}
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"main_categories": mainCategories,
		"sub_categories":  subCategories,
	})
}

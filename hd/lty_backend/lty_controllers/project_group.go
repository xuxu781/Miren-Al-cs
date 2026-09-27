package lty_controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"unicode/utf8"

	"lty_backend/lty_models"
)

// CreateProjectGroup 创建项目组
func CreateProjectGroup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userID := r.Context().Value("user_id").(uint)

	var req struct {
		Name string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": 400, "message": "无效的请求参数"})
		return
	}

	if utf8.RuneCountInString(req.Name) > 10 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": 400, "message": "项目组名称不能超过10个字符"})
		return
	}

	group := lty_models.ProjectGroup{
		UserID: userID,
		Name:   req.Name,
	}

	if err := lty_models.CreateProjectGroup(&group); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": 500, "message": "创建项目组失败"})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"data": group,
	})
}

// UpdateProjectGroup 更新项目组名称
func UpdateProjectGroup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userID := r.Context().Value("user_id").(uint)

	var req struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.ID == 0 || req.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": 400, "message": "无效的请求参数"})
		return
	}

	if utf8.RuneCountInString(req.Name) > 10 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": 400, "message": "项目组名称不能超过10个字符"})
		return
	}

	if err := lty_models.UpdateProjectGroup(userID, req.ID, req.Name); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": 500, "message": "更新项目组失败"})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{"code": 200, "message": "更新成功"})
}

// DeleteProjectGroup 删除项目组
func DeleteProjectGroup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userID := r.Context().Value("user_id").(uint)

	var req struct {
		ID uint `json:"id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.ID == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": 400, "message": "无效的请求参数"})
		return
	}

	if err := lty_models.DeleteProjectGroup(userID, req.ID); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": 500, "message": "删除项目组失败"})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{"code": 200, "message": "删除成功"})
}

// GetProjectGroupDetail 获取单个项目组信息
func GetProjectGroupDetail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userID := r.Context().Value("user_id").(uint)
	
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		json.NewEncoder(w).Encode(map[string]interface{}{"code": 400, "message": "id is required"})
		return
	}
	
	id, err := strconv.Atoi(idStr)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{"code": 400, "message": "invalid id format"})
		return
	}
	
	group, err := lty_models.GetProjectGroupByID(userID, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": 500, "message": "获取项目组详情失败"})
		return
	}
	
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"data": group,
	})
}

// GetProjectGroups 获取用户的项目组列表
func GetProjectGroups(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userID := r.Context().Value("user_id").(uint)

	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("page_size")

	page := 1
	if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
		page = p
	}

	pageSize := 5 // default to 5
	if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 {
		pageSize = ps
	}

	groups, total, err := lty_models.GetProjectGroupsPaginated(userID, page, pageSize)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": 500, "message": "获取项目组列表失败"})
		return
	}

	if groups == nil {
		groups = make([]lty_models.ProjectGroup, 0)
	}

	// 获取正在生成的项目组ID列表，用于前端全局统计
	generatingGroupIDs, _ := lty_models.GetGeneratingProjectGroupIDs(userID)
	if generatingGroupIDs == nil {
		generatingGroupIDs = make([]uint, 0)
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{
			"list":                 groups,
			"total":                total,
			"page":                 page,
			"generating_group_ids": generatingGroupIDs,
		},
	})
}

// MoveSessionToGroup 移动会话到项目组
func MoveSessionToGroup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userID := r.Context().Value("user_id").(uint)

	var req struct {
		SessionID string `json:"session_id"`
		ProjectID uint   `json:"project_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.SessionID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": 400, "message": "无效的请求参数"})
		return
	}

	if err := lty_models.MoveSession(userID, req.SessionID, req.ProjectID); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": 500, "message": "移动会话失败"})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{"code": 200, "message": "移动成功"})
}

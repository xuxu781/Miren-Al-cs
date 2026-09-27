package lty_models

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"lty_backend/lty_config"
)

// formatUserError 根据原始日志提取对用户友好的错误信息
func formatUserError(rawLog string) string {
	lowerLog := strings.ToLower(rawLog)
	if strings.Contains(lowerLog, "policy") || strings.Contains(lowerLog, "safety") || strings.Contains(lowerLog, "violation") || strings.Contains(lowerLog, "违规") || strings.Contains(lowerLog, "sensitive") {
		return "提示词或参考图可能包含违规/敏感内容，请修改后重试"
	}
	if strings.Contains(lowerLog, "insufficient_quota") || strings.Contains(lowerLog, "balance") || strings.Contains(lowerLog, "额度") || strings.Contains(lowerLog, "fund") {
		return "服务额度不足，请稍后重试或联系管理员"
	}
	if strings.Contains(lowerLog, "积分不足") {
		return "积分不足，请前往充值页面进行充值"
	}
	if strings.Contains(lowerLog, "timeout") || strings.Contains(lowerLog, "deadline") || strings.Contains(lowerLog, "超时") || strings.Contains(lowerLog, "context canceled") {
		return "请求超时，请稍后重试"
	}
	if strings.Contains(lowerLog, "connection") || strings.Contains(lowerLog, "network") || strings.Contains(lowerLog, "网络") || strings.Contains(lowerLog, "refused") {
		return "网络连接异常，请稍后重试"
	}
	if strings.Contains(lowerLog, "invalid_request_error") || strings.Contains(lowerLog, "bad request") || strings.Contains(lowerLog, "400") || strings.Contains(lowerLog, "invalid") {
		return "无效的请求，请检查提示词或参数设置"
	}
	// 默认兜底报错，不暴露API字眼
	return "服务繁忙或遇到未知问题，请稍后重试"
}

// ExtractNumImages 提取生成的张数
func ExtractNumImages(logContent string) int {
	numImages := 1
	// 1. 新版并发生成：匹配 (Task X/Y)
	reTask := regexp.MustCompile(`\(Task \d+/(\d+)\)`)
	if matches := reTask.FindStringSubmatch(logContent); len(matches) > 1 {
		if num, err := strconv.Atoi(matches[1]); err == nil {
			numImages = num
		}
	} else {
		// 2. 旧版 JSON 请求体：匹配 "n": 3 或 "n":3
		reJSON := regexp.MustCompile(`"n"\s*:\s*(\d+)`)
		if matches := reJSON.FindStringSubmatch(logContent); len(matches) > 1 {
			if num, err := strconv.Atoi(matches[1]); err == nil {
				numImages = num
			}
		} else {
			// 3. 旧版 Multipart 请求体：匹配 n=3
			reMulti := regexp.MustCompile(`\bn=(\d+)`)
			if matches := reMulti.FindStringSubmatch(logContent); len(matches) > 1 {
				if num, err := strconv.Atoi(matches[1]); err == nil {
					numImages = num
				}
			} else {
				// 4. Gemini 原生协议 candidateCount 提取
				reGemini := regexp.MustCompile(`"candidateCount"\s*:\s*(\d+)`)
				if matches := reGemini.FindStringSubmatch(logContent); len(matches) > 1 {
					if num, err := strconv.Atoi(matches[1]); err == nil {
						numImages = num
					}
				}
			}
		}
	}
	return numImages
}

type Task struct {
	ID         uint      `json:"id"`
	UserID     uint      `json:"user_id"`
	Username   string    `json:"username"`
	SessionID  string    `json:"session_id"`
	SessionName string   `json:"session_name"`
	ProjectID  uint      `json:"project_id"`
	SeriesID   string    `json:"series_id"`
	Model      string    `json:"model"`
	ChannelName string   `json:"channel_name"` // 渠道名称
	Prompt     string    `json:"prompt"`
	Size       string    `json:"size"`
	Resolution string    `json:"resolution"`
	NumImages  int       `json:"num_images"`
	ReferenceImage string    `json:"reference_image"`
	ImageURL       string    `json:"image_url"`
	Status     int       `json:"status"` // 0: pending, 1: success, 2: failed, 3: partial success, 4: deleted
	LogContent string    `json:"log_content,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func CreateTask(task *Task) error {
	now := time.Now()
	if task.CreatedAt.IsZero() {
		task.CreatedAt = now
	}
	if task.UpdatedAt.IsZero() {
		task.UpdatedAt = now
	}
	
	var maxSortOrder int
	lty_config.DB.QueryRow("SELECT IFNULL(MAX(sort_order), 0) FROM lty_tasks WHERE user_id = ?", task.UserID).Scan(&maxSortOrder)
	
	query := "INSERT INTO lty_tasks (user_id, username, session_id, session_name, project_id, series_id, model, prompt, size, resolution, reference_image, status, log_content, created_at, updated_at, sort_order, num_images) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	result, err := lty_config.DB.Exec(query, task.UserID, task.Username, task.SessionID, task.SessionName, task.ProjectID, task.SeriesID, task.Model, task.Prompt, task.Size, task.Resolution, task.ReferenceImage, task.Status, task.LogContent, task.CreatedAt, task.UpdatedAt, maxSortOrder + 1, task.NumImages)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	task.ID = uint(id)
	return nil
}

func UpdateTask(task *Task) error {
	query := "UPDATE lty_tasks SET image_url = ?, status = ?, log_content = ?, model = ? WHERE id = ?"
	_, err := lty_config.DB.Exec(query, task.ImageURL, task.Status, task.LogContent, task.Model, task.ID)
	return err
}

func GetTasks(status *int, keyword string, limit, offset int) ([]Task, int, error) {
	var total int
	whereClause := "WHERE t.status != 4"
	var countArgs []interface{}

	if status != nil {
		whereClause += " AND t.status = ?"
		countArgs = append(countArgs, *status)
	}
	if keyword != "" {
		whereClause += " AND (t.prompt LIKE ? OR t.username LIKE ?)"
		countArgs = append(countArgs, "%"+keyword+"%", "%"+keyword+"%")
	}

	countQuery := "SELECT COUNT(*) FROM lty_tasks t " + whereClause
	err := lty_config.DB.QueryRow(countQuery, countArgs...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `
		SELECT 
			t.id, t.user_id, t.username, IFNULL(t.session_id, ''), IFNULL(t.session_name, ''), IFNULL(t.project_id, 0), IFNULL(t.series_id, ''), IFNULL(t.model, ''), 
			IFNULL((SELECT channel_name FROM lty_model_upstreams WHERE series_id = t.series_id AND logical_model = t.model LIMIT 1), '') as channel_name, 
			t.prompt, t.size, IFNULL(t.resolution, ''), IFNULL(t.reference_image, ''), IFNULL(t.image_url, ''), t.status, t.created_at, t.updated_at, IFNULL(t.log_content, '') as log_content, IFNULL(t.num_images, 1) as num_images
		FROM lty_tasks t
		` + whereClause + `
		ORDER BY t.created_at DESC 
		LIMIT ? OFFSET ?
	`
	
	args := append(countArgs, limit, offset)
	rows, err := lty_config.DB.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var t Task
		var resolution string
		var referenceImage string
		var imageUrl string
		var logContent string
		var model string
		var seriesId string
		var sessionName string
		var projectId uint
		var numImages int
		var channelName string
		if err := rows.Scan(&t.ID, &t.UserID, &t.Username, &t.SessionID, &sessionName, &projectId, &seriesId, &model, &channelName, &t.Prompt, &t.Size, &resolution, &referenceImage, &imageUrl, &t.Status, &t.CreatedAt, &t.UpdatedAt, &logContent, &numImages); err != nil {
			return nil, 0, err
		}
		
		t.SessionName = sessionName
		t.ProjectID = projectId
		t.SeriesID = seriesId
		t.Model = model
		t.ChannelName = channelName
		t.Resolution = resolution
		t.ReferenceImage = referenceImage
		t.ImageURL = imageUrl
		t.LogContent = logContent
		t.NumImages = numImages
		
		tasks = append(tasks, t)
	}

	return tasks, total, nil
}

// CleanupTimeoutTasks 自动清理超时卡住的任务
func CleanupTimeoutTasks() {
	rows, err := lty_config.DB.Query("SELECT id, updated_at, IFNULL(model, ''), IFNULL(image_url, ''), IFNULL(num_images, 1), user_id, IFNULL(series_id, ''), IFNULL(resolution, '') FROM lty_tasks WHERE status = 0")
	if err != nil {
		return
	}
	defer rows.Close()

	type pendingTask struct {
		ID         uint
		UpdatedAt  time.Time
		Model      string
		ImageURL   string
		NumImages  int
		UserID     uint
		SeriesID   string
		Resolution string
	}
	var tasks []pendingTask
	for rows.Next() {
		var t pendingTask
		if err := rows.Scan(&t.ID, &t.UpdatedAt, &t.Model, &t.ImageURL, &t.NumImages, &t.UserID, &t.SeriesID, &t.Resolution); err == nil {
			tasks = append(tasks, t)
		}
	}
	rows.Close()

	now := time.Now()
	for _, t := range tasks {
		timeoutSeconds := 300 // 默认 5 分钟
		if t.Model != "" {
			var ts int
			err := lty_config.DB.QueryRow("SELECT IFNULL(timeout_seconds, 0) FROM lty_model_upstreams WHERE logical_model = ? AND status = 'active' ORDER BY is_primary DESC, id DESC LIMIT 1", t.Model).Scan(&ts)
			if err == nil && ts > 0 {
				timeoutSeconds = ts
			}
		}

		// 增加 15 秒的缓冲时间，避免与正常的超时机制产生冲突
		if now.Sub(t.UpdatedAt) > time.Duration(timeoutSeconds+15)*time.Second {
			msg := fmt.Sprintf("任务执行超时(超过%d秒)", timeoutSeconds)
			CleanupSingleStuckTask(t.ID, t.ImageURL, t.NumImages, t.UserID, t.SeriesID, t.Model, t.Resolution, msg)
		}
	}
}

// CleanupAllStuckTasks 用于系统启动时清理所有僵尸任务
func CleanupAllStuckTasks() {
	rows, err := lty_config.DB.Query("SELECT id, IFNULL(image_url, ''), IFNULL(num_images, 1), user_id, IFNULL(series_id, ''), IFNULL(model, ''), IFNULL(resolution, '') FROM lty_tasks WHERE status = 0")
	if err != nil {
		return
	}
	defer rows.Close()

	type pendingTask struct {
		ID         uint
		ImageURL   string
		NumImages  int
		UserID     uint
		SeriesID   string
		Model      string
		Resolution string
	}
	var tasks []pendingTask
	for rows.Next() {
		var t pendingTask
		if err := rows.Scan(&t.ID, &t.ImageURL, &t.NumImages, &t.UserID, &t.SeriesID, &t.Model, &t.Resolution); err == nil {
			tasks = append(tasks, t)
		}
	}
	rows.Close()

	for _, t := range tasks {
		CleanupSingleStuckTask(t.ID, t.ImageURL, t.NumImages, t.UserID, t.SeriesID, t.Model, t.Resolution, "任务超时(系统重启或意外中断)")
	}
}

// CleanupUnusedImages 自动清理未被任何任务引用的历史图片
func CleanupUnusedImages() {
	// 获取清理设置
	enableSetting, _ := GetSetting("enable_auto_clean_images")
	if enableSetting == nil || enableSetting.KeyValue != "true" {
		return
	}

	timeSetting, _ := GetSetting("auto_clean_images_time")
	unitSetting, _ := GetSetting("auto_clean_images_unit")
	
	timeValue := 7
	if timeSetting != nil {
		if v, err := strconv.Atoi(timeSetting.KeyValue); err == nil && v >= 0 {
			timeValue = v
		}
	}
	
	unitValue := "days"
	if unitSetting != nil && unitSetting.KeyValue != "" {
		unitValue = unitSetting.KeyValue
	}

	uploadDirs := []string{"./uploads/images", "./uploads/system"}
	for _, dir := range uploadDirs {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			continue
		}
	}

	// 查找所有被引用的图片URL
	query := "SELECT IFNULL(image_url, ''), IFNULL(reference_image, '') FROM lty_tasks"
	rows, err := lty_config.DB.Query(query)
	if err != nil {
		log.Printf("Failed to query used images: %v", err)
		return
	}
	defer rows.Close()

	usedFiles := make(map[string]bool)
	for rows.Next() {
		var imgURL, refImg string
		if err := rows.Scan(&imgURL, &refImg); err == nil {
			for _, file := range extractFilesToDelete(imgURL) {
				usedFiles[filepath.ToSlash(filepath.Clean(file))] = true
			}
			for _, file := range extractFilesToDelete(refImg) {
				usedFiles[filepath.ToSlash(filepath.Clean(file))] = true
			}
		}
	}

	// 查找所有被引用的系统文件
	settingsQuery := "SELECT IFNULL(key_value, '') FROM lty_settings"
	settingsRows, err := lty_config.DB.Query(settingsQuery)
	if err == nil {
		defer settingsRows.Close()
		for settingsRows.Next() {
			var val string
			if err := settingsRows.Scan(&val); err == nil {
				for _, file := range extractFilesToDelete(val) {
					usedFiles[filepath.ToSlash(filepath.Clean(file))] = true
				}
			}
		}
	}

	var cutoffTime time.Time
	now := time.Now()
	switch unitValue {
	case "minutes":
		cutoffTime = now.Add(-time.Duration(timeValue) * time.Minute)
	case "hours":
		cutoffTime = now.Add(-time.Duration(timeValue) * time.Hour)
	case "days":
		cutoffTime = now.AddDate(0, 0, -timeValue)
	default:
		cutoffTime = now.AddDate(0, 0, -timeValue)
	}

	deletedCount := 0
	// 遍历上传目录，删除未被引用且过期的文件
	for _, uploadDir := range uploadDirs {
		err = filepath.Walk(uploadDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				// 判断文件是否在被引用的列表中
				relativePath := filepath.ToSlash(filepath.Clean("/" + path))
				
				// 对于Windows路径转换处理
				relativePath = strings.ReplaceAll(relativePath, "lty_backend/uploads", "uploads")
				relativePath = strings.ReplaceAll(relativePath, "/hd/", "/")
				relativePath = strings.ReplaceAll(relativePath, "//", "/")

				// 只要不被引用且满足过期时间要求，就删除
				if !usedFiles[relativePath] && info.ModTime().Before(cutoffTime) {
					if err := os.Remove(path); err == nil {
						deletedCount++
						log.Printf("Auto cleaning unused image: %s", path)
					}
				}
			}
			return nil
		})
		if err != nil {
			log.Printf("Error walking through uploads dir %s: %v", uploadDir, err)
		}
	}

	// 如果有文件被清理，记录日志到 lty_settings 表中，方便前端查看
	if deletedCount > 0 {
		logEntry := fmt.Sprintf("[%s] 成功清理了 %d 张未被引用且过期的历史图片", time.Now().Format("2006-01-02 15:04:05"), deletedCount)
		
		existingLog, _ := GetSetting("auto_clean_images_log")
		var newLog string
		if existingLog != nil && existingLog.KeyValue != "" {
			// 限制日志条数，保留最新 7 条（即旧日志保留 6 条，加上本次 1 条）
			logs := strings.Split(existingLog.KeyValue, "\n")
			if len(logs) > 6 {
				logs = logs[:6]
			}
			newLog = logEntry + "\n" + strings.Join(logs, "\n")
		} else {
			newLog = logEntry
		}
		SaveSetting("auto_clean_images_log", newLog, "自动清理历史图片日志")
	}
}

// CleanupSingleStuckTask 处理单个卡住的任务，包含状态更新和积分退还
func CleanupSingleStuckTask(taskID uint, imageURL string, numImages int, userID uint, seriesID, model, resolution, reasonMsg string) {
	var successCount int
	if imageURL != "" && imageURL != "[]" {
		var urls []string
		if strings.HasPrefix(imageURL, "[") {
			if err := json.Unmarshal([]byte(imageURL), &urls); err == nil {
				successCount = len(urls)
			}
		} else {
			successCount = 1
		}
	}

	failedCount := numImages - successCount
	if failedCount < 0 {
		failedCount = 0
	}

	status := 2 // 默认失败
	var logMsg string
	timeStr := time.Now().Format("2006-01-02 15:04:05")

	if successCount > 0 {
		status = 3 // 部分成功
		logMsg = fmt.Sprintf("\n%s | %s\nPartially generated %d images. %d failed.", timeStr, reasonMsg, successCount, failedCount)
	} else {
		logMsg = fmt.Sprintf("\n%s | %s", timeStr, reasonMsg)
	}

	// 1. 更新任务状态
	lty_config.DB.Exec("UPDATE lty_tasks SET status = ?, log_content = IFNULL(log_content, '') || ? WHERE id = ?", status, logMsg, taskID)

	// 2. 如果有失败张数，退还积分
	if failedCount > 0 {
		// 查找模型配置以获取单张积分
		targetId := seriesID
		if targetId == "" {
			targetId = model
		}
		
		creditsPerImage := 1
		if targetId != "" {
			var resolutionConfigsStr string
			err := lty_config.DB.QueryRow("SELECT IFNULL(resolution_configs, '') FROM lty_model_upstreams WHERE (series_id = ? OR logical_model = ? OR name = ?) AND status = 'active' ORDER BY is_primary DESC, id DESC LIMIT 1", targetId, targetId, targetId).Scan(&resolutionConfigsStr)
			if err == nil && resolutionConfigsStr != "" {
				var configs map[string]map[string]interface{}
				if json.Unmarshal([]byte(resolutionConfigsStr), &configs) == nil {
					if resConf, ok := configs[resolution]; ok {
						if val, exists := resConf["credits_per_image"]; exists {
							if floatVal, ok := val.(float64); ok {
								creditsPerImage = int(floatVal)
							}
						}
					}
				}
			}
		}

		refundPoints := creditsPerImage * failedCount
		if refundPoints > 0 {
			lty_config.DB.Exec("UPDATE lty_users SET points = points + ? WHERE id = ?", refundPoints, userID)
			reasonStr := fmt.Sprintf("生图失败退还 %d 张 - 模型: %s, 分辨率: %s", failedCount, targetId, resolution)
			lty_config.DB.Exec("INSERT INTO lty_point_records (user_id, points_change, title, reason) VALUES (?, ?, ?, ?)", userID, refundPoints, "系统退回", reasonStr)
		}
	}
}

type Session struct {
	SessionID string    `json:"session_id"`
	Title     string    `json:"title"`
	ImageURL  string    `json:"image_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Status    int       `json:"status"`
	ProjectID uint      `json:"project_id"`
	SortOrder int       `json:"sort_order"`
	GeneratingCount int `json:"generating_count"`
}

func GetUserSessions(userID uint, projectID int, limit, offset int) ([]Session, int, error) {
	var total int
	
	countQuery := "SELECT COUNT(DISTINCT session_id) FROM lty_tasks WHERE user_id = ? AND session_id != ''"
	args := []interface{}{userID}
	if projectID != -1 {
		countQuery += " AND IFNULL(project_id, 0) = ?"
		args = append(args, projectID)
	}
	
	err := lty_config.DB.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `
		SELECT t_latest.session_id, IF(IFNULL(t_latest.session_name, '') = '', t_latest.prompt, t_latest.session_name) as title, IFNULL(t_latest.image_url, '') as image_url, t_latest.created_at, t_latest.updated_at, t_latest.status, IFNULL(t_latest.project_id, 0), IFNULL(t_latest.sort_order, 0), IFNULL(tc.generating_count, 0) as generating_count
		FROM lty_tasks t_latest
		INNER JOIN (
			SELECT session_id, MAX(id) as last_id
			FROM lty_tasks 
			WHERE user_id = ? AND session_id != '' `
	
	if projectID != -1 {
		query += ` AND IFNULL(project_id, 0) = ? `
	}
			
	query += `
			GROUP BY session_id
		) t2 ON t_latest.id = t2.last_id
		LEFT JOIN (
			SELECT session_id, COUNT(id) as generating_count
			FROM lty_tasks
			WHERE user_id = ? AND status = 0
			GROUP BY session_id
		) tc ON t_latest.session_id = tc.session_id
		ORDER BY t_latest.sort_order DESC, t_latest.id DESC 
		LIMIT ? OFFSET ?
	`
	
	queryArgs := []interface{}{userID}
	if projectID != -1 {
		queryArgs = append(queryArgs, projectID)
	}
	queryArgs = append(queryArgs, userID, limit, offset)

	rows, err := lty_config.DB.Query(query, queryArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var sessions []Session
	for rows.Next() {
		var s Session
		if err := rows.Scan(&s.SessionID, &s.Title, &s.ImageURL, &s.CreatedAt, &s.UpdatedAt, &s.Status, &s.ProjectID, &s.SortOrder, &s.GeneratingCount); err != nil {
			return nil, 0, err
		}
		sessions = append(sessions, s)
	}
	
	return sessions, total, nil
}

// GetSessionByID 获取单个会话信息
func GetSessionByID(userID uint, sessionID string) (*Session, error) {
	query := `
		SELECT t_latest.session_id, IF(IFNULL(t_latest.session_name, '') = '', t_latest.prompt, t_latest.session_name) as title, IFNULL(t_latest.image_url, '') as image_url, t_latest.created_at, t_latest.updated_at, t_latest.status, IFNULL(t_latest.project_id, 0), IFNULL(t_latest.sort_order, 0), IFNULL(tc.generating_count, 0) as generating_count
		FROM lty_tasks t_latest
		INNER JOIN (
			SELECT session_id, MAX(id) as last_id
			FROM lty_tasks 
			WHERE user_id = ? AND session_id = ?
			GROUP BY session_id
		) t2 ON t_latest.id = t2.last_id
		LEFT JOIN (
			SELECT session_id, COUNT(id) as generating_count
			FROM lty_tasks
			WHERE user_id = ? AND status = 0
			GROUP BY session_id
		) tc ON t_latest.session_id = tc.session_id
	`
	var s Session
	err := lty_config.DB.QueryRow(query, userID, sessionID, userID).Scan(&s.SessionID, &s.Title, &s.ImageURL, &s.CreatedAt, &s.UpdatedAt, &s.Status, &s.ProjectID, &s.SortOrder, &s.GeneratingCount)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func GetTasksBySession(userID uint, sessionID string, limit, offset int) ([]Task, error) {
	query := `
		SELECT 
			id, user_id, username, IFNULL(session_id, ''), IFNULL(session_name, ''), IFNULL(project_id, 0), IFNULL(series_id, ''), IFNULL(model, ''), prompt, size, IFNULL(resolution, ''), IFNULL(reference_image, ''), IFNULL(image_url, ''), status, created_at, updated_at, IFNULL(log_content, '') as log_content, IFNULL(num_images, 1) as num_images
		FROM lty_tasks 
		WHERE user_id = ? AND session_id = ? AND status != 4
		ORDER BY created_at DESC, id DESC
		LIMIT ? OFFSET ?
	`
	rows, err := lty_config.DB.Query(query, userID, sessionID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var t Task
		var resolution string
		var referenceImage string
		var imageUrl string
		var logContent string
		var model string
		var seriesId string
		var sessionName string
		var projectId uint
		var numImages int
		if err := rows.Scan(&t.ID, &t.UserID, &t.Username, &t.SessionID, &sessionName, &projectId, &seriesId, &model, &t.Prompt, &t.Size, &resolution, &referenceImage, &imageUrl, &t.Status, &t.CreatedAt, &t.UpdatedAt, &logContent, &numImages); err != nil {
			return nil, err
		}
		
		t.SessionName = sessionName
		t.ProjectID = projectId
		t.SeriesID = seriesId
		t.Model = model
		t.Resolution = resolution
		t.ReferenceImage = referenceImage
		t.ImageURL = imageUrl
		t.LogContent = logContent
		t.NumImages = numImages
		
		tasks = append(tasks, t)
	}

	// 逆序以保证时间正序
	for i, j := 0, len(tasks)-1; i < j; i, j = i+1, j-1 {
		tasks[i], tasks[j] = tasks[j], tasks[i]
	}

	return tasks, nil
}

// RenameSession 重命名会话
func RenameSession(userID uint, sessionID string, newName string) error {
	query := "UPDATE lty_tasks SET session_name = ? WHERE user_id = ? AND session_id = ?"
	_, err := lty_config.DB.Exec(query, newName, userID, sessionID)
	return err
}

// DeleteSession 删除会话
func DeleteSession(userID uint, sessionID string) error {
	// 删除数据库记录，图片保留，交由自动清理任务处理
	delQuery := "DELETE FROM lty_tasks WHERE user_id = ? AND session_id = ?"
	_, err := lty_config.DB.Exec(delQuery, userID, sessionID)
	return err
}

// DeleteTask 删除单个任务
func DeleteTask(userID uint, taskID uint64) error {
	// 标记为已删除(status=4)，清空图片字段交由自动清理任务处理物理文件
	// 这样可以保留空会话
	delQuery := "UPDATE lty_tasks SET status = 4, image_url = '', reference_image = '' WHERE id = ? AND user_id = ?"
	_, err := lty_config.DB.Exec(delQuery, taskID, userID)
	return err
}

func extractFilesToDelete(jsonStr string) []string {
	var files []string
	if jsonStr == "" {
		return files
	}
	var urls []string
	if strings.HasPrefix(jsonStr, "[") {
		json.Unmarshal([]byte(jsonStr), &urls)
	} else {
		urls = append(urls, jsonStr)
	}
	for _, u := range urls {
		if parsed, err := url.Parse(u); err == nil {
			if strings.HasPrefix(parsed.Path, "/uploads/") {
				files = append(files, parsed.Path)
			} else if strings.HasPrefix(parsed.Path, "uploads/") {
				files = append(files, "/"+parsed.Path)
			}
		} else if strings.HasPrefix(u, "/uploads/") {
			files = append(files, u)
		} else if strings.HasPrefix(u, "uploads/") {
			files = append(files, "/"+u)
		}
	}
	return files
}

// MoveSession 移动会话到指定项目组
func MoveSession(userID uint, sessionID string, projectID uint) error {
	query := "UPDATE lty_tasks SET project_id = ? WHERE user_id = ? AND session_id = ?"
	_, err := lty_config.DB.Exec(query, projectID, userID, sessionID)
	return err
}

// ReorderSessions 对会话进行排序
func ReorderSessions(userID uint, sessionIDs []string) error {
	if len(sessionIDs) == 0 {
		return nil
	}

	tx, err := lty_config.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback() // Will be ignored if committed

	// 1. Get current sort_order values for these sessions
	placeholders := make([]string, len(sessionIDs))
	args := make([]interface{}, len(sessionIDs)+1)
	args[0] = userID
	for i, id := range sessionIDs {
		placeholders[i] = "?"
		args[i+1] = id
	}

	querySelect := "SELECT session_id, MAX(sort_order), MAX(id) FROM lty_tasks WHERE user_id = ? AND session_id IN (" + strings.Join(placeholders, ",") + ") GROUP BY session_id"
	rows, err := tx.Query(querySelect, args...)
	if err != nil {
		return err
	}
	
	var currentOrders []int
	sessionToLastID := make(map[string]int)
	for rows.Next() {
		var sID string
		var sOrder int
		var lastID int
		if err := rows.Scan(&sID, &sOrder, &lastID); err == nil {
			currentOrders = append(currentOrders, sOrder)
			sessionToLastID[sID] = lastID
		}
	}
	rows.Close()

	// 2. Sort the current sort_orders descending
	sort.Slice(currentOrders, func(i, j int) bool {
		return currentOrders[i] > currentOrders[j]
	})

	// If for some reason we didn't get enough sort_orders, fallback
	if len(currentOrders) < len(sessionIDs) {
		var maxSort int
		tx.QueryRow("SELECT IFNULL(MAX(sort_order), 0) FROM lty_tasks WHERE user_id = ?", userID).Scan(&maxSort)
		for len(currentOrders) < len(sessionIDs) {
			maxSort++
			currentOrders = append([]int{maxSort}, currentOrders...)
		}
		sort.Slice(currentOrders, func(i, j int) bool {
			return currentOrders[i] > currentOrders[j]
		})
	}

	// Ensure strictly decreasing from bottom to top to handle duplicates
	if len(currentOrders) > 0 {
		for i := len(currentOrders) - 2; i >= 0; i-- {
			if currentOrders[i] <= currentOrders[i+1] {
				currentOrders[i] = currentOrders[i+1] + 1
			}
		}
	}

	// 3. 仅更新每个 session 的最后一条记录（MAX(id)），大幅提升更新性能，避免锁表和延迟
	var caseSql strings.Builder
	caseSql.WriteString("UPDATE lty_tasks SET sort_order = CASE id ")
	
	updateArgs := make([]interface{}, 0, len(sessionIDs)*2)
	idsToUpdate := make([]int, 0, len(sessionIDs))

	for i, sessionID := range sessionIDs {
		lastID, ok := sessionToLastID[sessionID]
		if !ok {
			continue
		}
		newOrder := currentOrders[i]
		caseSql.WriteString("WHEN ? THEN ? ")
		updateArgs = append(updateArgs, lastID, newOrder)
		idsToUpdate = append(idsToUpdate, lastID)
	}

	if len(idsToUpdate) > 0 {
		caseSql.WriteString("END WHERE user_id = ? AND id IN (")
		placeholders := make([]string, len(idsToUpdate))
		for i := range placeholders {
			placeholders[i] = "?"
		}
		caseSql.WriteString(strings.Join(placeholders, ",") + ")")

		finalArgs := append(updateArgs, userID)
		for _, id := range idsToUpdate {
			finalArgs = append(finalArgs, id)
		}

		_, err = tx.Exec(caseSql.String(), finalArgs...)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func GetUserTasks(userID uint, limit, offset int) ([]Task, int, error) {
	var total int
	err := lty_config.DB.QueryRow("SELECT COUNT(*) FROM lty_tasks WHERE user_id = ?", userID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `
		SELECT 
			id, user_id, username, IFNULL(session_id, ''), IFNULL(session_name, ''), IFNULL(project_id, 0), IFNULL(series_id, ''), IFNULL(model, ''), prompt, size, IFNULL(resolution, ''), IFNULL(reference_image, ''), IFNULL(image_url, ''), status, created_at, updated_at, IFNULL(log_content, '') as log_content
		FROM lty_tasks 
		WHERE user_id = ? 
		ORDER BY created_at DESC 
		LIMIT ? OFFSET ?
	`
	rows, err := lty_config.DB.Query(query, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var t Task
		var resolution string
		var referenceImage string
		var imageUrl string
		var logContent string
		var model string
		var seriesId string
		var sessionName string
		var projectId uint
		if err := rows.Scan(&t.ID, &t.UserID, &t.Username, &t.SessionID, &sessionName, &projectId, &seriesId, &model, &t.Prompt, &t.Size, &resolution, &referenceImage, &imageUrl, &t.Status, &t.CreatedAt, &t.UpdatedAt, &logContent); err != nil {
			return nil, 0, err
		}
		
		t.SessionName = sessionName
		t.ProjectID = projectId
		t.SeriesID = seriesId
		t.Model = model
		t.Resolution = resolution
		t.ReferenceImage = referenceImage
		t.ImageURL = imageUrl
		if t.Status == 2 {
			logContent = formatUserError(logContent)
		}
		t.LogContent = logContent
		
		// 从 log_content 提取生成的张数
		t.NumImages = ExtractNumImages(logContent)
		
		tasks = append(tasks, t)
	}

	return tasks, total, nil
}


package lty_controllers

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
	"database/sql"

	"lty_backend/lty_config"
	"lty_backend/lty_models"
	"lty_backend/lty_controllers/providers/image"
)

// GenerateImage 用户前台生图请求
func GenerateImage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID := r.Context().Value("user_id").(uint)
	username := r.Context().Value("username").(string)

	var req struct {
		Prompt         string `json:"prompt"`
		Resolution     string `json:"resolution"`
		AspectRatio    string `json:"aspect_ratio"`
		NumImages      int    `json:"num_images"`
		ReferenceImage string `json:"reference_image"`
		SeriesID       string `json:"series_id"`
		SessionID      string `json:"session_id"`
		ProjectID      uint   `json:"project_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": 400, "message": "无效的请求参数"})
		return
	}

	// Fetch upstream configuration to check credits needed
	var creditsPerImage float64 = 1 // Default to 1
	var resolutionConfigs sql.NullString
	err := lty_config.DB.QueryRow(`
		SELECT resolution_configs
		FROM lty_model_upstreams 
		WHERE (series_id = ? OR logical_model = ?) AND status = 'active'
		ORDER BY is_primary DESC, id DESC LIMIT 1
	`, req.SeriesID, req.SeriesID).Scan(&resolutionConfigs)

	if err == nil && resolutionConfigs.Valid && resolutionConfigs.String != "" {
		var configs map[string]interface{}
		if err := json.Unmarshal([]byte(resolutionConfigs.String), &configs); err == nil {
			if tierConfig, ok := configs[req.Resolution].(map[string]interface{}); ok {
				if cpi, ok := tierConfig["credits_per_image"].(float64); ok {
					creditsPerImage = cpi
				}
			}
		}
	}

	totalCreditsNeeded := creditsPerImage
	if req.NumImages > 1 {
		totalCreditsNeeded = creditsPerImage * float64(req.NumImages)
	}

	// 不再拼接，保持和上游接口一致
	size := req.AspectRatio
	if size == "" {
		size = "auto"
	}

	task := lty_models.Task{
		UserID:         userID,
		Username:       username,
		SessionID:      req.SessionID,
		ProjectID:      req.ProjectID,
		SeriesID:       req.SeriesID,
		Prompt:         req.Prompt,
		Size:           size,
		Resolution:     req.Resolution,
		NumImages:      req.NumImages, // 存储实际选择的张数
		ReferenceImage: req.ReferenceImage,
		Status:         0, // Pending
		LogContent:     time.Now().Format("2006-01-02 15:04:05") + " | Task created",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := lty_models.CreateTask(&task); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": 500, "message": "创建任务失败"})
		return
	}

	// Check user points
	var currentPoints float64
	err = lty_config.DB.QueryRow("SELECT points FROM lty_users WHERE id = ?", userID).Scan(&currentPoints)
	if err != nil {
		task.Status = 2
		task.LogContent = task.LogContent + "\n" + time.Now().Format("2006-01-02 15:04:05") + " | 获取用户信息失败"
		lty_models.UpdateTask(&task)
		
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": 500, "message": "获取用户信息失败"})
		return
	}

	if currentPoints < totalCreditsNeeded {
		task.Status = 2
		task.LogContent = task.LogContent + "\n" + time.Now().Format("2006-01-02 15:04:05") + " | 积分不足"
		lty_models.UpdateTask(&task)
		
		w.WriteHeader(http.StatusBadRequest)
		// 即使是 400 也可以返回 task_id，以便前端更新临时任务的 ID
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code": 400, 
			"message": "积分不足",
			"data": map[string]interface{}{
				"task_id": task.ID,
			},
		})
		return
	}

	// Deduct points
	_, err = lty_config.DB.Exec("UPDATE lty_users SET points = points - ? WHERE id = ? AND points >= ?", totalCreditsNeeded, userID, totalCreditsNeeded)
	if err != nil {
		task.Status = 2
		task.LogContent = task.LogContent + "\n" + time.Now().Format("2006-01-02 15:04:05") + " | 扣除积分失败"
		lty_models.UpdateTask(&task)
		
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": 500, "message": "扣除积分失败"})
		return
	}

	// 记录积分消耗
	reasonStr := fmt.Sprintf("生成图片 %d 张 - 模型: %s, 比例: %s, 分辨率: %s", req.NumImages, req.SeriesID, size, req.Resolution)
	lty_config.DB.Exec("INSERT INTO lty_point_records (user_id, points_change, title, reason) VALUES (?, ?, ?, ?)", userID, -totalCreditsNeeded, "图片生成", reasonStr)
	// Trigger async task processing here...
	go processTask(task.ID, req.SeriesID, req.Prompt, req.Resolution, req.AspectRatio, req.NumImages, req.ReferenceImage)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{
			"task_id": task.ID,
		},
	})
}

func processTask(taskID uint, seriesID string, prompt string, resolution string, aspectRatio string, numImages int, referenceImage string) {
	// Fetch task
	var task lty_models.Task
	err := lty_config.DB.QueryRow("SELECT id, user_id, status, log_content, IFNULL(model, '') FROM lty_tasks WHERE id = ?", taskID).Scan(&task.ID, &task.UserID, &task.Status, &task.LogContent, &task.Model)
	if err != nil {
		return
	}

	updateLog := func(msg string) {
		timeStr := time.Now().Format("2006-01-02 15:04:05")
		// 每次获取最新任务状态，防止并发覆盖
		var currentLog string
		lty_config.DB.QueryRow("SELECT log_content FROM lty_tasks WHERE id = ?", task.ID).Scan(&currentLog)

		task.LogContent = currentLog + "\n" + timeStr + " | " + msg
		lty_models.UpdateTask(&task)
	}

	updateLog("Starting to process task...")

	// 1. Fetch upstream configuration
	var upstream lty_models.ModelUpstream
	var resolutionConfigs sql.NullString
	creditsPerImage := 1.0 // Default to 1.0
	err = lty_config.DB.QueryRow(`
		SELECT id, logical_model, provider, connection_url, api_key, IFNULL(timeout_seconds, 0), resolution_configs
		FROM lty_model_upstreams 
		WHERE (series_id = ? OR logical_model = ?) AND status = 'active'
		ORDER BY is_primary DESC, id DESC LIMIT 1
	`, seriesID, seriesID).Scan(&upstream.ID, &upstream.LogicalModel, &upstream.Provider, &upstream.ConnectionURL, &upstream.APIKey, &upstream.TimeoutSeconds, &resolutionConfigs)

	if err == nil && resolutionConfigs.Valid && resolutionConfigs.String != "" {
		var configs map[string]interface{}
		if parseErr := json.Unmarshal([]byte(resolutionConfigs.String), &configs); parseErr == nil {
			if tierConfig, ok := configs[resolution].(map[string]interface{}); ok {
				if cpi, ok := tierConfig["credits_per_image"].(float64); ok {
					creditsPerImage = cpi
				}
				if model, ok := tierConfig["model"].(string); ok && model != "" {
					upstream.LogicalModel = model
				}
				if sizeParam, ok := tierConfig["size_parameter"].(string); ok && sizeParam != "" {
					upstream.SizeParameter = sizeParam
				}
			}
		}
	}

	numImagesRefund := numImages
	if numImagesRefund <= 0 {
		numImagesRefund = 1
	}

	if err != nil {
		updateLog(fmt.Sprintf("Failed to find active model upstream for %s", seriesID))
		task.Status = 2 // Failed
		lty_models.UpdateTask(&task)
		// Refund points
		lty_config.DB.Exec("UPDATE lty_users SET points = points + ? WHERE id = ?", creditsPerImage*float64(numImagesRefund), task.UserID)
		lty_config.DB.Exec("INSERT INTO lty_point_records (user_id, points_change, title, reason) VALUES (?, ?, ?, ?)", task.UserID, creditsPerImage*float64(numImagesRefund), "系统退回", "生图失败退还")
		return
	}

	task.Model = upstream.LogicalModel
	updateLog(fmt.Sprintf("Using upstream model: %s (Provider: %s)", upstream.LogicalModel, upstream.Provider))

	// 2. Prepare request to upstream API
	reqURL := upstream.ConnectionURL

	// Check setting mode
	useURL := true
	setting, err := lty_models.GetSetting("enable_reference_image")
	if err == nil && setting != nil && setting.KeyValue == "false" {
		useURL = false
	}

	// 根据渠道不同，独立处理 URL 补全逻辑
	if upstream.Provider == "gemini" || upstream.Provider == "gemini_image" {
		// 1. Gemini 原生协议处理
		if !strings.HasSuffix(reqURL, ":predict") && !strings.HasSuffix(reqURL, ":generateContent") {
			// 如果 URL 看起来不像原生的 gemini 结尾，但我们又不想强行拼装
			// 这里仅仅处理最基础的根域名情况，如果你填了具体的路径，就原样发送
			parsedURL, parseErr := url.Parse(reqURL)
			if parseErr == nil && (parsedURL.Path == "" || parsedURL.Path == "/") {
				// 针对带有参考图的情况（修图/垫图），由于原生 Gemini 绘图接口对这部分支持非常特殊
				// 如果检测到用户想要上传参考图，且没有填写具体路径，我们将其智能补全为 OpenAI 的 edits 接口，让中转站去处理
				if !useURL && referenceImage != "" {
					if strings.HasSuffix(reqURL, "/") {
						reqURL = reqURL + "v1/images/edits"
					} else {
						reqURL = reqURL + "/v1/images/edits"
					}
				} else {
					if !strings.HasSuffix(reqURL, "/") {
						reqURL = reqURL + "/"
					}
					reqURL = reqURL + "v1beta/models/" + upstream.LogicalModel + ":generateContent"
				}
			}
		}
	} else if upstream.Provider == "openai" || upstream.Provider == "openai_image" || upstream.Provider == "midjourney_image" || upstream.Provider == "sd_image" {
		// 2. OpenAI 标准协议处理
		parsedURL, parseErr := url.Parse(reqURL)
		if parseErr == nil && (parsedURL.Path == "" || parsedURL.Path == "/" || parsedURL.Path == "/v1" || parsedURL.Path == "/v1/") {
			suffix := "images/generations"
			if !useURL && referenceImage != "" {
				suffix = "images/edits"
			}
			if strings.HasSuffix(reqURL, "/v1") || strings.HasSuffix(reqURL, "/v1/") {
				if strings.HasSuffix(reqURL, "/") {
					reqURL = reqURL + suffix
				} else {
					reqURL = reqURL + "/" + suffix
				}
			} else {
				if strings.HasSuffix(reqURL, "/") {
					reqURL = reqURL + "v1/" + suffix
				} else {
					reqURL = reqURL + "/v1/" + suffix
				}
			}
		}
	}

	// 按照上游接口要求，size 传比例，resolution 传档位
	sizeForAPI := aspectRatio
	if sizeForAPI == "" {
		sizeForAPI = "auto"
	}

	// 处理尺寸参数转换
	if upstream.SizeParameter == "pixel_resolution" || upstream.SizeParameter == "pixel_size" {
		if sizeForAPI == "auto" {
			sizeForAPI = "1:1" // 如果是 auto，在需要转像素时默认当做 1:1 处理
		}
		// 解析比例，例如 "16:9"
		parts := strings.Split(sizeForAPI, ":")
		if len(parts) == 2 {
			wRatio, errW := strconv.ParseFloat(parts[0], 64)
			hRatio, errH := strconv.ParseFloat(parts[1], 64)
			
			if errW == nil && errH == nil && wRatio > 0 && hRatio > 0 {
				// 确定基准像素，例如 1K = 1024, 2K = 2048, 4K = 4096
				basePixels := 1024.0 // 默认 1K
				if resolution == "2K" {
					basePixels = 2048.0
				} else if resolution == "4K" {
					basePixels = 4096.0
				} else if resolution == "8K" {
					basePixels = 8192.0
				}

				// 计算宽高
				// 保持基准像素为长边
				var width, height float64
				if wRatio >= hRatio {
					width = basePixels
					height = basePixels * (hRatio / wRatio)
				} else {
					height = basePixels
					width = basePixels * (wRatio / hRatio)
				}

				// 取整，并确保是 64 的倍数（很多模型要求宽高是64的倍数）
				finalWidth := int(math.Round(width/64.0)) * 64
				finalHeight := int(math.Round(height/64.0)) * 64

				sizeForAPI = fmt.Sprintf("%dx%d", finalWidth, finalHeight)
			}
		}
	}

	resolutionForAPI := resolution

	var allImageUrls []string
	var mu sync.Mutex
	var wg sync.WaitGroup
	var failedCount int

	totalImages := numImages
	if totalImages <= 0 {
		totalImages = 1
	}

	for i := 0; i < totalImages; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()

			var currentGenerated []string
			var err error

			if upstream.Provider == "gemini" || upstream.Provider == "gemini_image" {
				currentGenerated, err = image.ProcessGeminiImageTask(&upstream, prompt, sizeForAPI, resolutionForAPI, referenceImage, useURL, updateLog, index, totalImages, DownloadImageToLocal, GetBaseURL)
			} else if upstream.Provider == "openai" || upstream.Provider == "openai_image" || upstream.Provider == "midjourney_image" || upstream.Provider == "sd_image" {
				currentGenerated, err = image.ProcessOpenAIImageTask(&upstream, prompt, sizeForAPI, resolutionForAPI, referenceImage, useURL, updateLog, index, totalImages, DownloadImageToLocal, GetBaseURL)
			} else {
				updateLog("Error: Unsupported provider format configuration.")
				err = fmt.Errorf("unsupported provider")
			}

			if err != nil || len(currentGenerated) == 0 {
				mu.Lock()
				failedCount++
				mu.Unlock()
				return
			}

			if len(currentGenerated) > 0 {
				mu.Lock()
				allImageUrls = append(allImageUrls, currentGenerated...)
				imageUrlsJson, _ := json.Marshal(allImageUrls)
				task.ImageURL = string(imageUrlsJson)
				// Save partial progress
				lty_models.UpdateTask(&task)
				updateLog(fmt.Sprintf("Task %d completed. Current total: %d/%d", index+1, len(allImageUrls), totalImages))
				mu.Unlock()
			} else {
				mu.Lock()
				failedCount++
				mu.Unlock()
			}
		}(i)
	}

	wg.Wait()

	if len(allImageUrls) == 0 {
		updateLog("All requests failed")
		task.Status = 2
		lty_models.UpdateTask(&task)
		// Refund all points
		lty_config.DB.Exec("UPDATE lty_users SET points = points + ? WHERE id = ?", creditsPerImage*float64(totalImages), task.UserID)
		reasonStr := fmt.Sprintf("生图失败退还 %d 张 - 模型: %s, 比例: %s, 分辨率: %s", totalImages, seriesID, aspectRatio, resolution)
		lty_config.DB.Exec("INSERT INTO lty_point_records (user_id, points_change, title, reason) VALUES (?, ?, ?, ?)", task.UserID, creditsPerImage*float64(totalImages), "系统退回", reasonStr)
		return
	}

	if failedCount > 0 {
		updateLog(fmt.Sprintf("Partially generated %d images. %d failed.", len(allImageUrls), failedCount))
		// Refund points for failed count
		lty_config.DB.Exec("UPDATE lty_users SET points = points + ? WHERE id = ?", creditsPerImage*float64(failedCount), task.UserID)
		reasonStr := fmt.Sprintf("生图失败退还 %d 张 - 模型: %s, 比例: %s, 分辨率: %s", failedCount, seriesID, aspectRatio, resolution)
		lty_config.DB.Exec("INSERT INTO lty_point_records (user_id, points_change, title, reason) VALUES (?, ?, ?, ?)", task.UserID, creditsPerImage*float64(failedCount), "系统退回", reasonStr)
		
		// If some succeeded and some failed, mark status as 3 (Partial Success)
		task.Status = 3
	} else {
		updateLog(fmt.Sprintf("Successfully generated all %d image(s)", len(allImageUrls)))
		task.Status = 1 // Success
	}

	lty_models.UpdateTask(&task)

	// Final log save
	timeStr := time.Now().Format("2006-01-02 15:04:05")
	var currentLog string
	lty_config.DB.QueryRow("SELECT log_content FROM lty_tasks WHERE id = ?", task.ID).Scan(&currentLog)
	task.LogContent = currentLog + "\n" + timeStr + " | Task completed successfully."
	lty_models.UpdateTask(&task)
}

// ReplaceImageDomain 替换图片URL的域名部分为当前 baseURL
func ReplaceImageDomain(imageURL string, baseURL string) string {
	if imageURL == "" || baseURL == "" {
		return imageURL
	}
	
	baseURL = strings.TrimRight(baseURL, "/")

	var urls []string
	if strings.HasPrefix(imageURL, "[") {
		if err := json.Unmarshal([]byte(imageURL), &urls); err == nil {
			for i, u := range urls {
				if parsed, err := url.Parse(u); err == nil {
					if strings.HasPrefix(parsed.Path, "/uploads/") {
						urls[i] = baseURL + parsed.Path
					}
				} else if strings.HasPrefix(u, "/uploads/") {
					urls[i] = baseURL + u
				}
			}
			if b, err := json.Marshal(urls); err == nil {
				return string(b)
			}
		}
	} else {
		if parsed, err := url.Parse(imageURL); err == nil {
			if strings.HasPrefix(parsed.Path, "/uploads/") {
				return baseURL + parsed.Path
			}
		} else if strings.HasPrefix(imageURL, "/uploads/") {
			return baseURL + imageURL
		}
	}
	return imageURL
}

// GetTasks 获取任务列表（带分页）
func GetTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("page_size")
	statusStr := r.URL.Query().Get("status")
	keyword := r.URL.Query().Get("keyword")

	page := 1
	if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
		page = p
	}

	pageSize := 10
	if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 {
		pageSize = ps
	}

	var status *int
	if statusStr != "" {
		if s, err := strconv.Atoi(statusStr); err == nil {
			status = &s
		}
	}

	offset := (page - 1) * pageSize

	tasks, total, err := lty_models.GetTasks(status, keyword, pageSize, offset)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "获取任务列表失败"})
		return
	}

	if tasks == nil {
		tasks = make([]lty_models.Task, 0)
	}

	baseURL := GetBaseURL(r)
	for i := range tasks {
		tasks[i].ImageURL = ReplaceImageDomain(tasks[i].ImageURL, baseURL)
		tasks[i].ReferenceImage = ReplaceImageDomain(tasks[i].ReferenceImage, baseURL)
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"list":  tasks,
		"total": total,
	})
}

// GetTaskLogs 获取任务日志
func GetTaskLogs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	taskIDStr := r.URL.Query().Get("task_id")
	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil || taskID <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "无效的任务ID"})
		return
	}

	var task lty_models.Task
	err = lty_config.DB.QueryRow("SELECT id, user_id, username, IFNULL(session_id, ''), IFNULL(session_name, ''), IFNULL(project_id, 0), IFNULL(series_id, ''), prompt, size, IFNULL(resolution, ''), IFNULL(reference_image, ''), status, IFNULL(image_url, ''), IFNULL(log_content, ''), created_at, updated_at, IFNULL(model, ''), IFNULL(num_images, 1) FROM lty_tasks WHERE id = ?", taskID).Scan(
		&task.ID, &task.UserID, &task.Username, &task.SessionID, &task.SessionName, &task.ProjectID, &task.SeriesID, &task.Prompt, &task.Size, &task.Resolution, &task.ReferenceImage, &task.Status, &task.ImageURL, &task.LogContent, &task.CreatedAt, &task.UpdatedAt, &task.Model, &task.NumImages,
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "获取任务失败"})
		return
	}
	
	baseURL := GetBaseURL(r)
	task.ImageURL = ReplaceImageDomain(task.ImageURL, baseURL)
	task.ReferenceImage = ReplaceImageDomain(task.ReferenceImage, baseURL)

	// 不再从日志提取，直接使用数据库字段
	logContent := task.LogContent
	createdAt := task.CreatedAt

	logs := []map[string]interface{}{}
	if logContent != "" {
		lines := strings.Split(logContent, "\n")
		var currentLogMsg string
		var currentLogTime time.Time
		var hasCurrent bool

		for _, line := range lines {
			line = strings.TrimRight(line, "\r")
			if strings.TrimSpace(line) == "" {
				continue
			}

			isNewLog := false
			if len(line) >= 22 {
				parts := strings.SplitN(line, " | ", 2)
				if len(parts) == 2 {
					if parsedTime, err := time.ParseInLocation("2006-01-02 15:04:05", parts[0], time.Local); err == nil {
						if hasCurrent {
							logs = append(logs, map[string]interface{}{
								"id":          len(logs) + 1,
								"created_at":  currentLogTime,
								"log_content": strings.TrimSpace(currentLogMsg),
							})
						}
						currentLogTime = parsedTime
						currentLogMsg = parts[1]
						hasCurrent = true
						isNewLog = true
					}
				}
			}

			if !isNewLog {
				if hasCurrent {
					currentLogMsg += "\n" + line
				} else {
					currentLogTime = createdAt
					currentLogMsg = line
					hasCurrent = true
				}
			}
		}
		if hasCurrent {
			logs = append(logs, map[string]interface{}{
				"id":          len(logs) + 1,
				"created_at":  currentLogTime,
				"log_content": strings.TrimSpace(currentLogMsg),
			})
		}
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"list": logs,
		"raw":  logContent,
		"task": task,
	})
}

// GetUserTasks 获取当前用户的任务列表（带分页）
func GetUserTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID := r.Context().Value("user_id").(uint)

	sessionID := r.URL.Query().Get("session_id")
	if sessionID != "" {
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

		tasks, err := lty_models.GetTasksBySession(userID, sessionID, pageSize, offset)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "获取会话任务失败"})
			return
		}
		if tasks == nil {
			tasks = make([]lty_models.Task, 0)
		}
		
		baseURL := GetBaseURL(r)
		for i := range tasks {
			tasks[i].ImageURL = ReplaceImageDomain(tasks[i].ImageURL, baseURL)
			tasks[i].ReferenceImage = ReplaceImageDomain(tasks[i].ReferenceImage, baseURL)
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"list": tasks,
		})
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

	tasks, total, err := lty_models.GetUserTasks(userID, pageSize, offset)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "获取任务列表失败"})
		return
	}

	if tasks == nil {
		tasks = make([]lty_models.Task, 0)
	}

	baseURL := GetBaseURL(r)
	for i := range tasks {
		tasks[i].ImageURL = ReplaceImageDomain(tasks[i].ImageURL, baseURL)
		tasks[i].ReferenceImage = ReplaceImageDomain(tasks[i].ReferenceImage, baseURL)
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"list":  tasks,
		"total": total,
	})
}

// RenameSession 重命名会话
func RenameSession(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userID := r.Context().Value("user_id").(uint)

	var req struct {
		SessionID string `json:"session_id"`
		Name      string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.SessionID == "" || req.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": 400, "message": "无效的请求参数"})
		return
	}

	if err := lty_models.RenameSession(userID, req.SessionID, req.Name); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": 500, "message": "重命名失败"})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{"code": 200, "message": "重命名成功"})
}

// DeleteTask 删除单个任务
func DeleteTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userID := r.Context().Value("user_id").(uint)

	var req struct {
		TaskID uint64 `json:"task_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.TaskID == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": 400, "message": "无效的请求参数"})
		return
	}

	err := lty_models.DeleteTask(userID, req.TaskID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": 500, "message": "删除失败"})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{"code": 200, "message": "删除成功"})
}

// DeleteSession 删除会话
func DeleteSession(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userID := r.Context().Value("user_id").(uint)

	var req struct {
		SessionID string `json:"session_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.SessionID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": 400, "message": "无效的请求参数"})
		return
	}

	err := lty_models.DeleteSession(userID, req.SessionID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": 500, "message": "删除失败"})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{"code": 200, "message": "删除成功"})
}

// GetSessionDetail 获取单个会话信息
func GetSessionDetail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID := r.Context().Value("user_id").(uint)
	sessionID := r.URL.Query().Get("session_id")
	if sessionID == "" {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code": 400,
			"msg":  "session_id is required",
		})
		return
	}

	session, err := lty_models.GetSessionByID(userID, sessionID)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code": 500,
			"msg":  err.Error(),
		})
		return
	}

	baseURL := GetBaseURL(r)
	session.ImageURL = ReplaceImageDomain(session.ImageURL, baseURL)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"data": session,
	})
}

// GetUserSessions 获取当前用户的会话列表
func GetUserSessions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID := r.Context().Value("user_id").(uint)

	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("page_size")

	page := 1
	if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
		page = p
	}

	pageSize := 13 // default to 13
	if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 {
		pageSize = ps
	}

	projectID := -1
	if pidStr := r.URL.Query().Get("project_id"); pidStr != "" {
		if pid, err := strconv.Atoi(pidStr); err == nil {
			projectID = pid
		}
	}

	offset := (page - 1) * pageSize

	sessions, total, err := lty_models.GetUserSessions(userID, projectID, pageSize, offset)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "获取会话列表失败"})
		return
	}

	if sessions == nil {
		sessions = make([]lty_models.Session, 0)
	}

	baseURL := GetBaseURL(r)
	for i := range sessions {
		sessions[i].ImageURL = ReplaceImageDomain(sessions[i].ImageURL, baseURL)
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"list":  sessions,
		"total": total,
	})
}

// ReorderSessions 拖拽排序会话
func ReorderSessions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userID := r.Context().Value("user_id").(uint)

	var req struct {
		SessionIDs []string `json:"session_ids"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || len(req.SessionIDs) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": 400, "message": "无效的请求参数"})
		return
	}

	if err := lty_models.ReorderSessions(userID, req.SessionIDs); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": 500, "message": "排序失败"})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{"code": 200, "message": "排序成功"})
}

// 原来的 extractNumImages
func extractNumImages(logContent string) int {
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
			}
		}
	}
	return numImages
}

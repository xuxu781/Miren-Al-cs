package lty_controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"lty_backend/lty_config"
	"lty_backend/lty_models"
)

// CreateModelUpstream creates a new model upstream
func CreateModelUpstream(w http.ResponseWriter, r *http.Request) {
	var req lty_models.ModelUpstream
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.IsPrimary {
		// Set other upstreams with the same logical model and model type to non-primary
		_, err := lty_config.DB.Exec("UPDATE lty_model_upstreams SET is_primary = FALSE WHERE logical_model = ? AND model_type = ?", req.LogicalModel, req.ModelType)
		if err != nil {
			http.Error(w, "Failed to update other upstreams", http.StatusInternalServerError)
			return
		}
	}

	query := `INSERT INTO lty_model_upstreams (model_type, series_id, logical_model, provider, connection_url, api_key, channel_name, resolution_tiers, aspect_ratios, image_counts, activity_tag, activity_tag_color, size_parameter, max_reference_images, resolution_configs, timeout_seconds, capabilities, billing_strategy, execution_strategy, is_primary, status, created_at, updated_at)
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	now := time.Now()
	status := "active"
	if req.Status != "" {
		status = req.Status
	}

	if req.TimeoutSeconds == 0 {
		req.TimeoutSeconds = 60
	}

	res, err := lty_config.DB.Exec(query, req.ModelType, req.SeriesID, req.LogicalModel, req.Provider, req.ConnectionURL, req.APIKey, req.ChannelName, req.ResolutionTiers, req.AspectRatios, req.ImageCounts, req.ActivityTag, req.ActivityTagColor, req.SizeParameter, req.MaxReferenceImages, req.ResolutionConfigs, req.TimeoutSeconds, req.Capabilities, req.BillingStrategy, req.ExecutionStrategy, req.IsPrimary, status, now, now)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, _ := res.LastInsertId()
	req.ID = int(id)
	req.Status = status
	req.CreatedAt = now
	req.UpdatedAt = now

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    req,
	})
}

// GetModelUpstreams returns all model upstreams
func GetModelUpstreams(w http.ResponseWriter, r *http.Request) {
	rows, err := lty_config.DB.Query("SELECT id, model_type, series_id, logical_model, provider, connection_url, api_key, channel_name, resolution_tiers, aspect_ratios, image_counts, activity_tag, activity_tag_color, size_parameter, max_reference_images, resolution_configs, timeout_seconds, capabilities, billing_strategy, execution_strategy, is_primary, status, created_at, updated_at FROM lty_model_upstreams ORDER BY id DESC")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var upstreams []lty_models.ModelUpstream
	for rows.Next() {
		var u lty_models.ModelUpstream
		var createdAt, updatedAt []uint8
		var resConfigs sql.NullString
		if err := rows.Scan(&u.ID, &u.ModelType, &u.SeriesID, &u.LogicalModel, &u.Provider, &u.ConnectionURL, &u.APIKey, &u.ChannelName, &u.ResolutionTiers, &u.AspectRatios, &u.ImageCounts, &u.ActivityTag, &u.ActivityTagColor, &u.SizeParameter, &u.MaxReferenceImages, &resConfigs, &u.TimeoutSeconds, &u.Capabilities, &u.BillingStrategy, &u.ExecutionStrategy, &u.IsPrimary, &u.Status, &createdAt, &updatedAt); err != nil {
			continue
		}
		if resConfigs.Valid {
			u.ResolutionConfigs = resConfigs.String
		}

		if t, err := time.Parse("2006-01-02 15:04:05", string(createdAt)); err == nil {
			u.CreatedAt = t
		}
		if t, err := time.Parse("2006-01-02 15:04:05", string(updatedAt)); err == nil {
			u.UpdatedAt = t
		}
		upstreams = append(upstreams, u)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    upstreams,
	})
}

// UpdateModelUpstream updates an existing model upstream
func UpdateModelUpstream(w http.ResponseWriter, r *http.Request) {
	var req lty_models.ModelUpstream
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.ID == 0 {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	if req.IsPrimary {
		// Set other upstreams with the same logical model and model type to non-primary
		_, err := lty_config.DB.Exec("UPDATE lty_model_upstreams SET is_primary = FALSE WHERE logical_model = ? AND model_type = ? AND id != ?", req.LogicalModel, req.ModelType, req.ID)
		if err != nil {
			http.Error(w, "Failed to update other upstreams", http.StatusInternalServerError)
			return
		}
	}

	query := `UPDATE lty_model_upstreams SET model_type = ?, series_id = ?, logical_model = ?, provider = ?, connection_url = ?, api_key = ?, channel_name = ?, resolution_tiers = ?, aspect_ratios = ?, image_counts = ?, activity_tag = ?, activity_tag_color = ?, size_parameter = ?, max_reference_images = ?, resolution_configs = ?, timeout_seconds = ?, capabilities = ?, billing_strategy = ?, execution_strategy = ?, is_primary = ?, status = ?, updated_at = ? WHERE id = ?`
	now := time.Now()

	_, err := lty_config.DB.Exec(query, req.ModelType, req.SeriesID, req.LogicalModel, req.Provider, req.ConnectionURL, req.APIKey, req.ChannelName, req.ResolutionTiers, req.AspectRatios, req.ImageCounts, req.ActivityTag, req.ActivityTagColor, req.SizeParameter, req.MaxReferenceImages, req.ResolutionConfigs, req.TimeoutSeconds, req.Capabilities, req.BillingStrategy, req.ExecutionStrategy, req.IsPrimary, req.Status, now, req.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Updated successfully",
	})
}

// DeleteModelUpstream deletes a model upstream
func DeleteModelUpstream(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ID int `json:"id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.ID == 0 {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	_, err := lty_config.DB.Exec("DELETE FROM lty_model_upstreams WHERE id = ?", req.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Deleted successfully",
	})
}

// GetPublicModelUpstreams returns active model upstreams for public use
func GetPublicModelUpstreams(w http.ResponseWriter, r *http.Request) {
	// 聚合去重，根据 series_id 或 logical_model 获取供用户选择的模型列表
	rows, err := lty_config.DB.Query(`
		SELECT 
			CASE WHEN series_id != '' THEN series_id ELSE logical_model END as display_id,
			MAX(logical_model), 
			model_type,
			MAX(resolution_tiers),
			MAX(aspect_ratios),
			MAX(image_counts),
			MAX(max_reference_images),
			MAX(activity_tag),
			MAX(activity_tag_color),
			MAX(resolution_configs)
		FROM lty_model_upstreams 
		WHERE status = 'active' 
		GROUP BY display_id, model_type
		ORDER BY display_id ASC
	`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var models []map[string]interface{}
	for rows.Next() {
		var displayID, logicalModel string
		var resolutionTiers, aspectRatios, imageCounts, activityTag, activityTagColor string
		var modelType, maxRefImages int
		var resConfigs sql.NullString
		if err := rows.Scan(&displayID, &logicalModel, &modelType, &resolutionTiers, &aspectRatios, &imageCounts, &maxRefImages, &activityTag, &activityTagColor, &resConfigs); err != nil {
			continue
		}

		name := displayID
		if name == "" {
			name = logicalModel
		}

		// Filter out sensitive information from resolution_configs before sending to frontend
		var filteredResConfigs string
		if resConfigs.Valid && resConfigs.String != "" {
			var configs map[string]interface{}
			if err := json.Unmarshal([]byte(resConfigs.String), &configs); err == nil {
				for _, tierConfig := range configs {
					if tierMap, ok := tierConfig.(map[string]interface{}); ok {
						delete(tierMap, "model")
					}
				}
				if filteredBytes, err := json.Marshal(configs); err == nil {
					filteredResConfigs = string(filteredBytes)
				}
			}
		}

		models = append(models, map[string]interface{}{
			"series_id":            displayID,
			"name":                 name,
			"model_type":           modelType,
			"resolution_tiers":     resolutionTiers,
			"aspect_ratios":        aspectRatios,
			"image_counts":         imageCounts,
			"max_reference_images": maxRefImages,
			"activity_tag":         activityTag,
			"activity_tag_color":   activityTagColor,
			"resolution_configs":   filteredResConfigs,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    models,
	})
}

// SetPrimaryModelUpstream sets a model upstream as primary
func SetPrimaryModelUpstream(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ID           int    `json:"id"`
		LogicalModel string `json:"logical_model"`
		ModelType    int    `json:"model_type"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.ID == 0 || req.LogicalModel == "" || req.ModelType == 0 {
		http.Error(w, "ID, LogicalModel and ModelType are required", http.StatusBadRequest)
		return
	}

	tx, err := lty_config.DB.Begin()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec("UPDATE lty_model_upstreams SET is_primary = FALSE WHERE logical_model = ? AND model_type = ?", req.LogicalModel, req.ModelType)
	if err != nil {
		tx.Rollback()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec("UPDATE lty_model_upstreams SET is_primary = TRUE WHERE id = ?", req.ID)
	if err != nil {
		tx.Rollback()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Set primary successfully",
	})
}

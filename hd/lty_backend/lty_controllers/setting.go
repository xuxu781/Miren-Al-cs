package lty_controllers

import (
	"encoding/json"
	"lty_backend/lty_models"
	"net/http"
	"strings"
)

// SettingRequest payload
type SettingRequest struct {
	KeyName     string `json:"key_name"`
	KeyValue    string `json:"key_value"`
	Description string `json:"description"`
}

func GetSettings(w http.ResponseWriter, r *http.Request) {
	settings, err := lty_models.GetSettings()
	if err != nil {
		http.Error(w, `{"error": "Failed to get settings"}`, http.StatusInternalServerError)
		return
	}

	// Rewrite domain for site_logo if needed
	for i, s := range settings {
		if s.KeyName == "site_logo" && strings.Contains(s.KeyValue, "/uploads/") {
			parts := strings.SplitN(s.KeyValue, "/uploads/", 2)
			if len(parts) == 2 {
				baseURL := GetBaseURL(r)
				settings[i].KeyValue = baseURL + "/uploads/" + parts[1]
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"list": settings,
	})
}

func GetSetting(w http.ResponseWriter, r *http.Request) {
	keyName := r.URL.Query().Get("key_name")
	if keyName == "" {
		http.Error(w, `{"error": "key_name is required"}`, http.StatusBadRequest)
		return
	}

	setting, err := lty_models.GetSetting(keyName)
	if err != nil {
		http.Error(w, `{"error": "Failed to get setting"}`, http.StatusInternalServerError)
		return
	}

	if setting == nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": nil,
		})
		return
	}

	if setting.KeyName == "site_logo" && strings.Contains(setting.KeyValue, "/uploads/") {
		parts := strings.SplitN(setting.KeyValue, "/uploads/", 2)
		if len(parts) == 2 {
			baseURL := GetBaseURL(r)
			setting.KeyValue = baseURL + "/uploads/" + parts[1]
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data": setting,
	})
}

func SaveSetting(w http.ResponseWriter, r *http.Request) {
	var req SettingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Invalid JSON format"}`, http.StatusBadRequest)
		return
	}

	if req.KeyName == "" {
		http.Error(w, `{"error": "key_name is required"}`, http.StatusBadRequest)
		return
	}

	err := lty_models.SaveSetting(req.KeyName, req.KeyValue, req.Description)
	if err != nil {
		http.Error(w, `{"error": "Failed to save setting"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "success",
	})
}

// GetPublicSettings retrieves public settings for frontend
func GetPublicSettings(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	settings, err := lty_models.GetSettings()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": 500, "message": "获取配置失败"})
		return
	}

	publicKeys := map[string]bool{
		"site_name":              true,
		"site_logo":              true,
		"site_description":       true,
		"site_copyright":         true,
		"enable_reference_image": true,
		"image_domain":           true,
		"enable_recharge":        true,
		"recharge_link":          true,
		"inspiration_categories": true,
	}

	result := make(map[string]string)
	for _, s := range settings {
		if publicKeys[s.KeyName] {
			result[s.KeyName] = s.KeyValue
		}
	}

	// Replace the domain of site_logo if it contains /uploads/
	if logo, ok := result["site_logo"]; ok && logo != "" {
		if strings.Contains(logo, "/uploads/") {
			// Extract the relative path starting from /uploads/
			parts := strings.SplitN(logo, "/uploads/", 2)
			if len(parts) == 2 {
				baseURL := GetBaseURL(r)
				result["site_logo"] = baseURL + "/uploads/" + parts[1]
			}
		}
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"data": result,
	})
}

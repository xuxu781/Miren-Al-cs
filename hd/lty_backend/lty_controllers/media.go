package lty_controllers

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"lty_backend/lty_models"
)

func GetBaseURL(r *http.Request) string {
	setting, err := lty_models.GetSetting("image_domain")
	if err == nil && setting != nil && setting.KeyValue != "" {
		domain := strings.TrimRight(setting.KeyValue, "/")
		if !strings.HasPrefix(domain, "http") {
			scheme := "http"
			if r != nil && r.TLS != nil {
				scheme = "https"
			}
			domain = scheme + "://" + domain
		}
		return domain
	}

	scheme := "http"
	if r != nil && r.TLS != nil {
		scheme = "https"
	}
	host := "localhost:8080"
	if r != nil && r.Host != "" {
		host = r.Host
	}
	return fmt.Sprintf("%s://%s", scheme, host)
}

type MediaFile struct {
	Name      string `json:"name"`
	Path      string `json:"path"`
	URL       string `json:"url"`
	Size      int64  `json:"size"`
	CreatedAt string `json:"created_at"`
	Category  string `json:"category"`
}

func GetAdminMediaList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var mediaList []MediaFile
	baseURL := GetBaseURL(r)

	// Directories to scan
	dirsToScan := []string{
		filepath.Join("uploads", "images"),
		filepath.Join("uploads", "system"),
	}

	for _, baseDir := range dirsToScan {
		filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil // skip errors
			}
			if !info.IsDir() {
				// 动态判断分类，因为 path 中可能包含 system
				category := "用户图片"
				if strings.Contains(filepath.ToSlash(path), "/system/") || strings.Contains(filepath.ToSlash(path), "uploads/system/") {
					category = "系统图片"
				}

				// calculate url
				url := fmt.Sprintf("%s/%s", baseURL, filepath.ToSlash(path))
				mediaList = append(mediaList, MediaFile{
					Name:      info.Name(),
					Path:      filepath.ToSlash(path),
					URL:       url,
					Size:      info.Size(),
					CreatedAt: info.ModTime().Format("2006-01-02 15:04:05"),
					Category:  category,
				})
			}
			return nil
		})
	}

	// sort by created_at desc (explicit sort to handle multiple directories)
	sort.Slice(mediaList, func(i, j int) bool {
		return mediaList[i].CreatedAt > mediaList[j].CreatedAt
	})

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{
			"list": mediaList,
		},
	})
}

func DeleteAdminMedia(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req struct {
		Path string `json:"path"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": 400, "message": "参数错误"})
		return
	}

	if req.Path == "" || strings.Contains(req.Path, "..") {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": 400, "message": "非法路径"})
		return
	}

	// Make sure path is under uploads/images or uploads/system
	cleanPath := filepath.Clean(req.Path)
	cleanPathSlash := filepath.ToSlash(cleanPath)
	if !strings.HasPrefix(cleanPathSlash, "uploads/images") && !strings.HasPrefix(cleanPathSlash, "uploads/system") {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": 403, "message": "无权限删除此文件"})
		return
	}

	err := os.Remove(cleanPath)
	if err != nil && !os.IsNotExist(err) {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": 500, "message": "删除失败: " + err.Error()})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code":    200,
		"message": "删除成功",
	})
}

// ensureDir ensures the directory exists
func ensureDir(dirName string) error {
	err := os.MkdirAll(dirName, 0755)
	if err == nil || os.IsExist(err) {
		return nil
	}
	return err
}

// saveBase64Image saves a base64 string to a local file and returns the local URL path
func saveBase64Image(base64Data string, r *http.Request) (string, error) {
	// e.g. data:image/png;base64,iVBORw0KGgo...
	parts := strings.SplitN(base64Data, ",", 2)
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid base64 format")
	}

	ext := ".png"
	if strings.Contains(parts[0], "jpeg") || strings.Contains(parts[0], "jpg") {
		ext = ".jpg"
	} else if strings.Contains(parts[0], "webp") {
		ext = ".webp"
	} else if strings.Contains(parts[0], "gif") {
		ext = ".gif"
	}

	data, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return "", err
	}

	dateStr := time.Now().Format("20060102")
	uploadDir := filepath.Join("uploads", "images", dateStr)
	if err := ensureDir(uploadDir); err != nil {
		return "", err
	}

	b := make([]byte, 4)
	rand.Read(b)
	fileName := fmt.Sprintf("%d_%x%s", time.Now().UnixNano(), b, ext)
	filePath := filepath.Join(uploadDir, fileName)

	if err := ioutil.WriteFile(filePath, data, 0644); err != nil {
		return "", err
	}

	// generate local url
	baseURL := GetBaseURL(r)
	return fmt.Sprintf("%s/%s", baseURL, filepath.ToSlash(filePath)), nil
}

// DownloadImageToLocal downloads an image from a URL and saves it locally
func DownloadImageToLocal(imageURL string, r *http.Request) (string, error) {
	if strings.HasPrefix(imageURL, "data:") {
		return saveBase64Image(imageURL, r)
	}

	client := &http.Client{
		Timeout: 15 * time.Second,
	}
	resp, err := client.Get(imageURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download image, status code: %d", resp.StatusCode)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	ext := ".png"
	contentType := resp.Header.Get("Content-Type")
	if strings.Contains(contentType, "jpeg") || strings.Contains(contentType, "jpg") {
		ext = ".jpg"
	} else if strings.Contains(contentType, "webp") {
		ext = ".webp"
	} else if strings.Contains(contentType, "gif") {
		ext = ".gif"
	}

	dateStr := time.Now().Format("20060102")
	uploadDir := filepath.Join("uploads", "images", dateStr)
	if err := ensureDir(uploadDir); err != nil {
		return "", err
	}

	b := make([]byte, 4)
	rand.Read(b)
	fileName := fmt.Sprintf("%d_%x%s", time.Now().UnixNano(), b, ext)
	filePath := filepath.Join(uploadDir, fileName)

	if err := ioutil.WriteFile(filePath, data, 0644); err != nil {
		return "", err
	}

	baseURL := GetBaseURL(r)
	return fmt.Sprintf("%s/%s", baseURL, filepath.ToSlash(filePath)), nil
}

// UploadSystemMedia handles system file upload from admin panel (e.g. logo)
func UploadSystemMedia(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Max 10MB
	r.ParseMultipartForm(10 << 20)

	file, header, err := r.FormFile("file")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": 400, "message": "读取文件失败"})
		return
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": 500, "message": "读取文件内容失败"})
		return
	}

	ext := filepath.Ext(header.Filename)
	if ext == "" {
		ext = ".png"
	}

	dateStr := time.Now().Format("20060102")
	uploadDir := filepath.Join("uploads", "system", dateStr)
	if err := ensureDir(uploadDir); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": 500, "message": "创建目录失败"})
		return
	}

	b := make([]byte, 4)
	rand.Read(b)
	fileName := fmt.Sprintf("sys_%d_%x%s", time.Now().UnixNano(), b, ext)
	filePath := filepath.Join(uploadDir, fileName)

	if err := ioutil.WriteFile(filePath, data, 0644); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": 500, "message": "保存文件失败"})
		return
	}

	baseURL := GetBaseURL(r)
	url := fmt.Sprintf("%s/%s", baseURL, filepath.ToSlash(filePath))

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"data": map[string]string{"url": url},
	})
}

// UploadMedia handles file upload from frontend
func UploadMedia(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Max 10MB
	r.ParseMultipartForm(10 << 20)

	file, header, err := r.FormFile("file")
	if err != nil {
		// try json base64
		var req struct {
			Image string `json:"image"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err == nil && req.Image != "" {
			url, err := saveBase64Image(req.Image, r)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]interface{}{"code": 500, "message": "保存失败"})
				return
			}
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code": 200,
				"data": map[string]string{"url": url},
			})
			return
		}

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": 400, "message": "读取文件失败"})
		return
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": 500, "message": "读取文件内容失败"})
		return
	}

	ext := filepath.Ext(header.Filename)
	if ext == "" {
		ext = ".png"
	}

	dateStr := time.Now().Format("20060102")
	uploadDir := filepath.Join("uploads", "images", dateStr)
	if err := ensureDir(uploadDir); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": 500, "message": "创建目录失败"})
		return
	}

	b := make([]byte, 4)
	rand.Read(b)
	fileName := fmt.Sprintf("%d_%x%s", time.Now().UnixNano(), b, ext)
	filePath := filepath.Join(uploadDir, fileName)

	if err := ioutil.WriteFile(filePath, data, 0644); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"code": 500, "message": "保存文件失败"})
		return
	}

	baseURL := GetBaseURL(r)
	url := fmt.Sprintf("%s/%s", baseURL, filepath.ToSlash(filePath))

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"data": map[string]string{"url": url},
	})
}

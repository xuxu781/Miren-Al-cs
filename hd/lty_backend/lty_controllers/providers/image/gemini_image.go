package image

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"lty_backend/lty_models"
)

func ProcessGeminiImageTask(
	upstream *lty_models.ModelUpstream,
	prompt, sizeForAPI, resolutionForAPI, referenceImage string,
	useURL bool,
	updateLog func(string),
	index, totalImages int,
	downloadImageToLocal func(string, *http.Request) (string, error),
	getBaseURL func(*http.Request) string,
) ([]string, error) {
	reqURL := upstream.ConnectionURL

	// Gemini 原生协议处理: 补全 URL
	if !strings.HasSuffix(reqURL, ":predict") && !strings.HasSuffix(reqURL, ":generateContent") {
		parsedURL, parseErr := url.Parse(reqURL)
		if parseErr == nil && (parsedURL.Path == "" || parsedURL.Path == "/") {
			if !strings.HasSuffix(reqURL, "/") {
				reqURL = reqURL + "/"
			}
			reqURL = reqURL + "v1beta/models/" + upstream.LogicalModel + ":generateContent"
		}
	}

	reqFormat := "b64_json" // 默认优先请求 b64_json，避免下载 URL 缓慢导致重试

	for attempt := 0; attempt <= 1; attempt++ {

		// 遵循文档的 generationConfig 结构
		imageConfig := map[string]interface{}{}
		if sizeForAPI != "" && sizeForAPI != "auto" && !strings.Contains(sizeForAPI, "x") {
			imageConfig["aspectRatio"] = sizeForAPI
		}

		resForAPI := resolutionForAPI
		if resForAPI != "" {
			resForAPI = strings.ToUpper(resForAPI)
			imageConfig["imageSize"] = resForAPI
		}

		delivery := "URI"
		if reqFormat == "b64_json" {
			delivery = "BASE64" // 或者根据实际 Gemini 官方参数，若不支持重试，可能需要调整。
		}

		generationConfig := map[string]interface{}{
			"responseModalities": []string{"TEXT", "IMAGE"},
			"responseFormat": map[string]interface{}{
				"image": map[string]interface{}{
					"delivery": delivery,
				},
			},
		}

		if len(imageConfig) > 0 {
			generationConfig["imageConfig"] = imageConfig
		}

		// 构建 contents.parts
		parts := []map[string]interface{}{
			{
				"text": prompt,
			},
		}

		// 处理参考图 (Gemini 原生格式)
		if referenceImage != "" {
			var refImages []string
			if strings.HasPrefix(referenceImage, "[") && strings.HasSuffix(referenceImage, "]") {
				json.Unmarshal([]byte(referenceImage), &refImages)
			} else {
				refImages = append(refImages, referenceImage)
			}

			for _, refImg := range refImages {
				if idx := strings.Index(refImg, "/uploads/"); idx != -1 {
					mimeType := "image/png"
					ext := strings.ToLower(filepath.Ext(refImg))
					if ext == ".jpg" || ext == ".jpeg" {
						mimeType = "image/jpeg"
					} else if ext == ".webp" {
						mimeType = "image/webp"
					}

					if useURL {
						// URL 模式：强制将本站的旧图片域名替换为当前最新的公网 URL
						baseURL := getBaseURL(nil)
						publicURL := baseURL + refImg[idx:]
						
						parts = append(parts, map[string]interface{}{
							"fileData": map[string]interface{}{
								"mimeType": mimeType,
								"fileUri":  publicURL,
							},
						})
					} else {
						// Base64 模式：读取本地文件并编码
						localPath := refImg[idx+1:]
						fileBytes, err := os.ReadFile(localPath)
						if err == nil {
							base64Str := base64.StdEncoding.EncodeToString(fileBytes)
							parts = append(parts, map[string]interface{}{
								"inlineData": map[string]interface{}{
									"mimeType": mimeType,
									"data":     base64Str,
								},
							})
						}
					}
				} else if strings.HasPrefix(refImg, "http") {
					mimeType := "image/png"
					ext := strings.ToLower(filepath.Ext(refImg))
					if ext == ".jpg" || ext == ".jpeg" {
						mimeType = "image/jpeg"
					} else if ext == ".webp" {
						mimeType = "image/webp"
					}
					parts = append(parts, map[string]interface{}{
						"fileData": map[string]interface{}{
							"mimeType": mimeType,
							"fileUri":  refImg,
						},
					})
				}
			}
		}

		requestBody := map[string]interface{}{
			"contents": []map[string]interface{}{
				{
					"role":  "user",
					"parts": parts,
				},
			},
			"generationConfig": generationConfig,
		}

		reqBytes, _ := json.Marshal(requestBody)
		reqBodyStr := string(reqBytes)

		req, reqErr := http.NewRequest("POST", reqURL, bytes.NewBuffer(reqBytes))
		if reqErr != nil {
			updateLog("Failed to create request to upstream API: " + reqErr.Error())
			return nil, reqErr
		}
		req.Header.Set("Content-Type", "application/json; charset=utf-8")

		if upstream.APIKey != "" {
			req.Header.Set("Authorization", "Bearer "+upstream.APIKey)
		}

		logReqBody := reqBodyStr
		if len(logReqBody) > 2000 {
			logReqBody = logReqBody[:2000] + "... [truncated]"
		}
		updateLog(fmt.Sprintf("Sending request to upstream API (Task %d/%d): %s\nBody: %s", index+1, totalImages, reqURL, logReqBody))

		timeout := time.Duration(upstream.TimeoutSeconds) * time.Second
		if timeout <= 0 {
			timeout = 300 * time.Second
		}
		client := &http.Client{Timeout: timeout}

		reqStartTime := time.Now()
		resp, err := client.Do(req)
		if err != nil {
			updateLog(fmt.Sprintf("Upstream API request failed (Task %d): %v", index+1, err))
			return nil, err
		}
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		reqElapsed := time.Since(reqStartTime)

		logBody := string(body)
		if len(logBody) > 2000 {
			logBody = logBody[:2000] + "... [truncated]"
		}

		// 将上游返回的真实完整数据记录到日志中，供前端任务详情查看
		updateLog(fmt.Sprintf("Gemini upstream response (Task %d): %s", index+1, logBody))

		if resp.StatusCode != http.StatusOK {
			errMsg := fmt.Sprintf("Upstream API returned error status (Task %d): %d\nBody: %s", index+1, resp.StatusCode, logBody)
			
			// 如果是因为请求了 b64_json 导致上游报错，则自动切换为 url (URI) 格式重试一次
			if attempt == 0 && reqFormat == "b64_json" {
				updateLog(fmt.Sprintf("Upstream may not support b64_json (BASE64) format, automatically switching to 'url' (URI) format and retrying... (Task %d)", index+1))
				reqFormat = "url"
				continue
			}

			updateLog(errMsg)
			return nil, fmt.Errorf("upstream returned %d", resp.StatusCode)
		}

		updateLog(fmt.Sprintf("Gemini API request took %v to complete (Task %d)", reqElapsed, index+1))

	// 解析 Gemini 响应
		var currentGenerated []string
		var geminiRes struct {
			Predictions []struct {
				BytesBase64Encoded string `json:"bytesBase64Encoded"`
				MimeType           string `json:"mimeType"`
			} `json:"predictions"`
			Candidates []struct {
				Content struct {
					Parts []struct {
						Thought    bool   `json:"thought"`
						Text       string `json:"text"`
						InlineData struct {
							MimeType string `json:"mimeType"`
							Data     string `json:"data"`
						} `json:"inlineData"`
						FileData struct {
							FileUri  string `json:"fileUri"`
							MimeType string `json:"mimeType"`
						} `json:"fileData"`
					} `json:"parts"`
				} `json:"content"`
			} `json:"candidates"`
		}

		if err := json.Unmarshal(body, &geminiRes); err != nil {
			errMsg := fmt.Sprintf("Failed to parse upstream response (Task %d): %v", index+1, err)
			updateLog(errMsg)
			return nil, err
		}

		if len(geminiRes.Predictions) > 0 {
			for _, p := range geminiRes.Predictions {
				if p.BytesBase64Encoded != "" {
					mime := p.MimeType
					if mime == "" {
						mime = "image/png"
					}
					urlStr := "data:" + mime + ";base64," + p.BytesBase64Encoded
					localUrl, err := downloadImageToLocal(urlStr, nil)
					if err == nil {
						currentGenerated = append(currentGenerated, localUrl)
					} else {
						updateLog(fmt.Sprintf("Failed to download/save image to local server (Task %d, possibly due to 15s timeout): %v", index+1, err))
						currentGenerated = append(currentGenerated, urlStr)
					}
				}
			}
		} else if len(geminiRes.Candidates) > 0 {
			// 用于记录已经下载过的 URL，防止重复下载
			downloadedURLs := make(map[string]bool)

			for _, c := range geminiRes.Candidates {
				for _, part := range c.Content.Parts {
					if part.Thought {
						continue
					}

					// 根据请求格式决定优先级
					var urlStr string
					if reqFormat == "url" && part.FileData.FileUri != "" {
						urlStr = part.FileData.FileUri
					} else if reqFormat == "b64_json" && part.InlineData.Data != "" {
						mime := part.InlineData.MimeType
						if mime == "" {
							mime = "image/png"
						}
						urlStr = "data:" + mime + ";base64," + part.InlineData.Data
					} else if part.FileData.FileUri != "" {
						urlStr = part.FileData.FileUri
					} else if part.InlineData.Data != "" {
						mime := part.InlineData.MimeType
						if mime == "" {
							mime = "image/png"
						}
						urlStr = "data:" + mime + ";base64," + part.InlineData.Data
					}

					if urlStr != "" {
						if !downloadedURLs[urlStr] {
							downloadedURLs[urlStr] = true
							localUrl, err := downloadImageToLocal(urlStr, nil)
							if err == nil {
								currentGenerated = append(currentGenerated, localUrl)
							} else {
								updateLog(fmt.Sprintf("Failed to download/save image to local server (Task %d, possibly due to 15s timeout): %v", index+1, err))
								currentGenerated = append(currentGenerated, urlStr)
							}
						}
					} else if part.Text != "" {
						updateLog(fmt.Sprintf("Received text from Gemini: %s", part.Text))
					}
				}
			}
		}

		if len(currentGenerated) == 0 {
			updateLog(fmt.Sprintf("Upstream returned empty data or no valid images (Task %d)", index+1))
			return nil, fmt.Errorf("empty response")
		}

		return currentGenerated, nil
	}

	return nil, fmt.Errorf("task failed after retries")
}

package lty_middlewares

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"lty_backend/lty_config"
	"lty_backend/lty_utils"
)

// parseAndValidateToken 提取公共的解析 Token 逻辑
func parseAndValidateToken(w http.ResponseWriter, r *http.Request) *lty_utils.Claims {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "未提供授权Token"})
		return nil
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Token格式错误"})
		return nil
	}

	claims, err := lty_utils.ParseToken(parts[1])
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "无效的Token"})
		return nil
	}
	return claims
}

// AdminAuth 管理员接口认证中间件
func AdminAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := parseAndValidateToken(w, r)
		if claims == nil {
			return
		}

		if claims.Role != "admin" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(map[string]string{"error": "无权限访问管理员接口"})
			return
		}

		// 查询数据库确认管理员依然存在
		var count int
		err := lty_config.DB.QueryRow("SELECT COUNT(*) FROM lty_admins WHERE id = ?", claims.ID).Scan(&count)
		if err != nil || count == 0 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "管理员不存在或已被删除"})
			return
		}

		ctx := context.WithValue(r.Context(), "admin_id", claims.ID)
		ctx = context.WithValue(ctx, "username", claims.Username)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

// UserAuth 用户接口认证中间件
func UserAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := parseAndValidateToken(w, r)
		if claims == nil {
			return
		}

		if claims.Role != "user" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(map[string]string{"error": "无权限访问用户接口"})
			return
		}

		// 查询数据库确认用户依然存在且未被禁用
		var status int
		err := lty_config.DB.QueryRow("SELECT status FROM lty_users WHERE id = ?", claims.ID).Scan(&status)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "用户不存在或已被删除"})
			return
		}
		if status == 0 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(map[string]string{"error": "账号已被禁用"})
			return
		}

		ctx := context.WithValue(r.Context(), "user_id", claims.ID)
		ctx = context.WithValue(ctx, "username", claims.Username)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

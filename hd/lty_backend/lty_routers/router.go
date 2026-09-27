package lty_routers

import (
	"encoding/json"
	"net/http"

	"lty_backend/lty_controllers"
	"lty_backend/lty_middlewares"
)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Length, Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func SetupRouter() http.Handler {
	mux := http.NewServeMux()

	// Health check API
	mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "ok",
			"message": "Service is healthy",
		})
	})

	// Install APIs
	mux.HandleFunc("/api/check-install", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		lty_controllers.CheckInstall(w, r)
	})
	mux.HandleFunc("/api/check-db", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		lty_controllers.CheckDBConnection(w, r)
	})
	mux.HandleFunc("/api/install", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		lty_controllers.Install(w, r)
	})

	// Public APIs
	mux.HandleFunc("/api/public/settings", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		lty_controllers.GetPublicSettings(w, r)
	})

	mux.HandleFunc("/api/admin/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		lty_controllers.Login(w, r)
	})

	mux.HandleFunc("/api/user/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		lty_controllers.UserLogin(w, r)
	})

	mux.HandleFunc("/api/user/register", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		lty_controllers.UserRegister(w, r)
	})

	mux.HandleFunc("/api/user/info", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		lty_middlewares.UserAuth(lty_controllers.GetUserInfoForFrontend)(w, r)
	})

	// Static file serving for uploads
	mux.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("uploads"))))

	// Upload API
	mux.HandleFunc("/api/upload", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		lty_middlewares.UserAuth(lty_controllers.UploadMedia)(w, r)
	})

	mux.HandleFunc("/api/admin/system/upload", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		lty_middlewares.AdminAuth(lty_controllers.UploadSystemMedia)(w, r)
	})

	// Generate API
	mux.HandleFunc("/api/generate", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		lty_middlewares.UserAuth(lty_controllers.GenerateImage)(w, r)
	})

	// Public Models API (For User Selection)
	mux.HandleFunc("/api/models/list", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		lty_controllers.GetPublicModelUpstreams(w, r)
	})

	// User Tasks API
	mux.HandleFunc("/api/user/tasks/list", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		lty_middlewares.UserAuth(lty_controllers.GetUserTasks)(w, r)
	})

	// User Sessions API
	mux.HandleFunc("/api/user/sessions/detail", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		lty_middlewares.UserAuth(lty_controllers.GetSessionDetail)(w, r)
	})
	mux.HandleFunc("/api/user/sessions/list", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		lty_middlewares.UserAuth(lty_controllers.GetUserSessions)(w, r)
	})
	mux.HandleFunc("/api/user/sessions/rename", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		lty_middlewares.UserAuth(lty_controllers.RenameSession)(w, r)
	})
	mux.HandleFunc("/api/user/sessions/delete", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		lty_middlewares.UserAuth(lty_controllers.DeleteSession)(w, r)
	})
	mux.HandleFunc("/api/user/tasks/delete", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		lty_middlewares.UserAuth(lty_controllers.DeleteTask)(w, r)
	})
	mux.HandleFunc("/api/user/sessions/move", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		lty_middlewares.UserAuth(lty_controllers.MoveSessionToGroup)(w, r)
	})
	mux.HandleFunc("/api/user/sessions/reorder", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		lty_middlewares.UserAuth(lty_controllers.ReorderSessions)(w, r)
	})

	// Project Groups API
	mux.HandleFunc("/api/user/projects/detail", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		lty_middlewares.UserAuth(lty_controllers.GetProjectGroupDetail)(w, r)
	})
	mux.HandleFunc("/api/user/projects/create", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		lty_middlewares.UserAuth(lty_controllers.CreateProjectGroup)(w, r)
	})
	mux.HandleFunc("/api/user/projects/update", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		lty_middlewares.UserAuth(lty_controllers.UpdateProjectGroup)(w, r)
	})
	mux.HandleFunc("/api/user/projects/delete", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		lty_middlewares.UserAuth(lty_controllers.DeleteProjectGroup)(w, r)
	})
	mux.HandleFunc("/api/user/projects/list", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		lty_middlewares.UserAuth(lty_controllers.GetProjectGroups)(w, r)
	})

	// Protected Admin APIs
	mux.HandleFunc("/api/admin/info", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		lty_middlewares.AdminAuth(lty_controllers.GetUserInfo)(w, r)
	})
	mux.HandleFunc("/api/admin/create", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		lty_middlewares.AdminAuth(lty_controllers.CreateAdmin)(w, r)
	})
	mux.HandleFunc("/api/admin/list", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		lty_middlewares.AdminAuth(lty_controllers.GetAdmins)(w, r)
	})
	mux.HandleFunc("/api/admin/update-password", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		lty_middlewares.AdminAuth(lty_controllers.UpdatePassword)(w, r)
	})
	mux.HandleFunc("/api/admin/delete", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		lty_middlewares.AdminAuth(lty_controllers.DeleteAdmin)(w, r)
	})

	// User Management APIs
	mux.HandleFunc("/api/admin/users/create", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		lty_middlewares.AdminAuth(lty_controllers.CreateUser)(w, r)
	})
	mux.HandleFunc("/api/admin/users/list", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		lty_middlewares.AdminAuth(lty_controllers.GetUsers)(w, r)
	})
	mux.HandleFunc("/api/admin/users/update", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		lty_middlewares.AdminAuth(lty_controllers.UpdateUser)(w, r)
	})
	mux.HandleFunc("/api/admin/users/delete", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		lty_middlewares.AdminAuth(lty_controllers.DeleteUser)(w, r)
	})
	mux.HandleFunc("/api/admin/users/points/records", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		lty_middlewares.AdminAuth(lty_controllers.GetPointRecords)(w, r)
	})
	mux.HandleFunc("/api/user/points/records", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		lty_middlewares.UserAuth(lty_controllers.GetPointRecords)(w, r)
	})

	// Tasks APIs
	mux.HandleFunc("/api/admin/tasks/list", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		lty_middlewares.AdminAuth(lty_controllers.GetTasks)(w, r)
	})
	mux.HandleFunc("/api/admin/tasks/logs", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		lty_middlewares.AdminAuth(lty_controllers.GetTaskLogs)(w, r)
	})

	// Settings APIs
	mux.HandleFunc("/api/admin/settings/list", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		lty_middlewares.AdminAuth(lty_controllers.GetSettings)(w, r)
	})
	mux.HandleFunc("/api/admin/settings/get", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		lty_middlewares.AdminAuth(lty_controllers.GetSetting)(w, r)
	})
	mux.HandleFunc("/api/admin/settings/save", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		lty_middlewares.AdminAuth(lty_controllers.SaveSetting)(w, r)
	})

	// Media APIs
	mux.HandleFunc("/api/admin/media/list", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		lty_middlewares.AdminAuth(lty_controllers.GetAdminMediaList)(w, r)
	})
	mux.HandleFunc("/api/admin/media/delete", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		lty_middlewares.AdminAuth(lty_controllers.DeleteAdminMedia)(w, r)
	})

	// Model Upstream APIs
	mux.HandleFunc("/api/admin/upstreams/create", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		lty_middlewares.AdminAuth(lty_controllers.CreateModelUpstream)(w, r)
	})
	mux.HandleFunc("/api/admin/upstreams/list", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		lty_middlewares.AdminAuth(lty_controllers.GetModelUpstreams)(w, r)
	})
	mux.HandleFunc("/api/admin/upstreams/update", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		lty_middlewares.AdminAuth(lty_controllers.UpdateModelUpstream)(w, r)
	})
	mux.HandleFunc("/api/admin/upstreams/delete", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		lty_middlewares.AdminAuth(lty_controllers.DeleteModelUpstream)(w, r)
	})
	mux.HandleFunc("/api/admin/upstreams/set-primary", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		lty_middlewares.AdminAuth(lty_controllers.SetPrimaryModelUpstream)(w, r)
	})

	// Cdkey APIs
	mux.HandleFunc("/api/admin/cdkeys/generate", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		lty_middlewares.AdminAuth(lty_controllers.GenerateCdkeys)(w, r)
	})
	mux.HandleFunc("/api/admin/cdkeys/list", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		lty_middlewares.AdminAuth(lty_controllers.GetCdkeys)(w, r)
	})
	mux.HandleFunc("/api/admin/cdkeys/usages", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		lty_middlewares.AdminAuth(lty_controllers.GetCdkeyUsages)(w, r)
	})
	mux.HandleFunc("/api/admin/cdkeys/update", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		lty_middlewares.AdminAuth(lty_controllers.UpdateCdkey)(w, r)
	})
	mux.HandleFunc("/api/admin/cdkeys/void", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		lty_middlewares.AdminAuth(lty_controllers.VoidCdkey)(w, r)
	})
	mux.HandleFunc("/api/admin/cdkeys/enable", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		lty_middlewares.AdminAuth(lty_controllers.EnableCdkey)(w, r)
	})
	mux.HandleFunc("/api/admin/cdkeys/delete", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		lty_middlewares.AdminAuth(lty_controllers.DeleteCdkey)(w, r)
	})
	mux.HandleFunc("/api/admin/cdkeys/batch-delete", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		lty_middlewares.AdminAuth(lty_controllers.BatchDeleteCdkey)(w, r)
	})
	mux.HandleFunc("/api/user/cdkeys/use", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		lty_middlewares.UserAuth(lty_controllers.UseCdkey)(w, r)
	})

	// Inspirations APIs
	mux.HandleFunc("/api/admin/inspirations/categories", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		lty_middlewares.AdminAuth(lty_controllers.GetInspirationCategories)(w, r)
	})
	mux.HandleFunc("/api/admin/inspirations/create", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		lty_middlewares.AdminAuth(lty_controllers.CreateInspiration)(w, r)
	})
	mux.HandleFunc("/api/admin/inspirations/list", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		lty_middlewares.AdminAuth(lty_controllers.GetInspirations)(w, r)
	})
	mux.HandleFunc("/api/admin/inspirations/update", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		lty_middlewares.AdminAuth(lty_controllers.UpdateInspiration)(w, r)
	})
	mux.HandleFunc("/api/admin/inspirations/delete", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		lty_middlewares.AdminAuth(lty_controllers.DeleteInspiration)(w, r)
	})
	mux.HandleFunc("/api/admin/inspirations/batch-delete", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		lty_middlewares.AdminAuth(lty_controllers.BatchDeleteInspiration)(w, r)
	})
	mux.HandleFunc("/api/public/inspirations/list", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		lty_controllers.GetPublicInspirations(w, r)
	})
	mux.HandleFunc("/api/public/inspirations/categories", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		lty_controllers.GetPublicInspirationCategories(w, r)
	})

	return corsMiddleware(mux)
}

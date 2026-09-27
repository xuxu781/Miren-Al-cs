package main

import (
	"log"
	"net/http"
	"time"

	"lty_backend/lty_config"
	"lty_backend/lty_models"
	"lty_backend/lty_routers"
)

func main() {
	// Load Config
	lty_config.LoadConfig()

	// Initialize Database if installed
	if lty_config.IsInstalled {
		if err := lty_config.InitDB(); err != nil {
			log.Fatalf("Failed to initialize DB: %v", err)
		}
		// 每次启动时自动执行数据表迁移和更新，确保结构最新
		if err := lty_config.AutoMigrate(); err != nil {
			log.Printf("Failed to auto migrate DB: %v", err)
		}

		// 启动时清理由于系统重启导致的僵尸任务
		lty_models.CleanupAllStuckTasks()

		// 启动后台定时任务，定期清理超时卡住的任务
		go func() {
			for {
				time.Sleep(30 * time.Second) // 每半分钟检查一次
				lty_models.CleanupTimeoutTasks()
			}
		}()

		// 启动后台定时任务，定期清理未被引用的历史图片
		go func() {
			for {
				lty_models.CleanupUnusedImages()
				time.Sleep(1 * time.Minute) // 每分钟检查一次，具体的过期判断在 CleanupUnusedImages 内部处理
			}
		}()
	}

	// Setup Router
	handler := lty_routers.SetupRouter()

	// Start server
	log.Println("Starting server on :8088...")
	if err := http.ListenAndServe(":8088", handler); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

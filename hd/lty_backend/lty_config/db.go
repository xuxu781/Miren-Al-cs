package lty_config

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB() error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		Config.DBUser, Config.DBPass, Config.DBHost, Config.DBPort, Config.DBName)
	
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	
	if err := db.Ping(); err != nil {
		return err
	}
	
	DB = db
	return nil
}

func AutoMigrate() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS lty_admins (
			id INT AUTO_INCREMENT PRIMARY KEY,
			username VARCHAR(100) NOT NULL UNIQUE,
			password VARCHAR(255) NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`,
		`CREATE TABLE IF NOT EXISTS lty_users (
			id INT AUTO_INCREMENT PRIMARY KEY,
			username VARCHAR(100) NOT NULL UNIQUE,
			password VARCHAR(255) NOT NULL,
			email VARCHAR(100) DEFAULT '',
			phone VARCHAR(20) DEFAULT '',
			status INT NOT NULL DEFAULT 1 COMMENT '1: active, 0: disabled',
			points DECIMAL(20,4) NOT NULL DEFAULT 0.0000,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`,
		`CREATE TABLE IF NOT EXISTS lty_point_records (
			id INT AUTO_INCREMENT PRIMARY KEY,
			user_id INT NOT NULL,
			points_change DECIMAL(20,4) NOT NULL,
			title VARCHAR(255) DEFAULT '',
			reason VARCHAR(255) DEFAULT '',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`,
		`CREATE TABLE IF NOT EXISTS lty_model_upstreams (
			id INT AUTO_INCREMENT PRIMARY KEY,
			model_type INT NOT NULL DEFAULT 4 COMMENT '1: text, 2: video, 4: image',
			series_id VARCHAR(100) NOT NULL,
			logical_model VARCHAR(100) NOT NULL,
			provider VARCHAR(50) NOT NULL,
			connection_url VARCHAR(255) DEFAULT '',
			api_key VARCHAR(255) DEFAULT '',
			channel_name VARCHAR(255) DEFAULT '',
			resolution_tiers VARCHAR(255) DEFAULT '',
			aspect_ratios VARCHAR(255) DEFAULT '',
			image_counts VARCHAR(50) DEFAULT '',
			activity_tag VARCHAR(100) DEFAULT '',
			activity_tag_color VARCHAR(20) DEFAULT '',
			size_parameter VARCHAR(50) DEFAULT '',
			max_reference_images INT DEFAULT 0,
			resolution_configs TEXT,
			timeout_seconds INT DEFAULT 300,
			capabilities TEXT,
			billing_strategy TEXT,
			execution_strategy TEXT,
			is_primary BOOLEAN DEFAULT FALSE,
			status VARCHAR(20) DEFAULT 'active',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`,
		`CREATE TABLE IF NOT EXISTS lty_tasks (
			id INT AUTO_INCREMENT PRIMARY KEY,
			user_id INT NOT NULL,
			username VARCHAR(100) NOT NULL,
			session_id VARCHAR(100) DEFAULT '',
			session_name VARCHAR(255) DEFAULT '',
			project_id INT DEFAULT 0,
			series_id VARCHAR(100) DEFAULT '',
			model VARCHAR(100) DEFAULT '',
			prompt TEXT NOT NULL,
			size VARCHAR(50) DEFAULT '',
			resolution VARCHAR(50) DEFAULT '',
			num_images INT DEFAULT 1,
			reference_image TEXT,
			image_url TEXT,
			status INT NOT NULL DEFAULT 0 COMMENT '0: pending, 1: success, 2: failed, 3: partial success, 4: deleted',
			log_content TEXT,
			sort_order INT DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`,
		`CREATE TABLE IF NOT EXISTS lty_project_groups (
			id INT AUTO_INCREMENT PRIMARY KEY,
			user_id INT NOT NULL,
			name VARCHAR(100) NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`,
		`CREATE TABLE IF NOT EXISTS lty_settings (
			id INT AUTO_INCREMENT PRIMARY KEY,
			key_name VARCHAR(100) NOT NULL UNIQUE,
			key_value TEXT,
			description VARCHAR(255) DEFAULT '',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`,
		`CREATE TABLE IF NOT EXISTS lty_cdkeys (
			id INT AUTO_INCREMENT PRIMARY KEY,
			cdkey VARCHAR(100) NOT NULL UNIQUE,
			points DECIMAL(10,2) NOT NULL DEFAULT 0,
			status INT NOT NULL DEFAULT 0 COMMENT '0: unused, 1: exhausted',
			max_uses INT NOT NULL DEFAULT 1,
			current_uses INT NOT NULL DEFAULT 0,
			used_by_user_id INT NULL,
			expires_at DATETIME NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			used_at DATETIME NULL
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`,
		`CREATE TABLE IF NOT EXISTS lty_cdkey_usages (
			id INT AUTO_INCREMENT PRIMARY KEY,
			cdkey_id INT NOT NULL,
			user_id INT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			UNIQUE KEY uk_cdkey_user (cdkey_id, user_id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`,
		`CREATE TABLE IF NOT EXISTS lty_inspirations (
			id INT AUTO_INCREMENT PRIMARY KEY,
			content TEXT NOT NULL,
			image_url VARCHAR(1000) DEFAULT '',
			main_category VARCHAR(100) DEFAULT '',
			sub_category VARCHAR(100) DEFAULT '',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`,
	}
	
	for _, q := range queries {
		if _, err := DB.Exec(q); err != nil {
			return err
		}
	}

	// 兼容已有表结构，安全地追加新字段（忽略执行错误）
	DB.Exec(`ALTER TABLE lty_cdkeys ADD COLUMN max_uses INT NOT NULL DEFAULT 1`)
	DB.Exec(`ALTER TABLE lty_cdkeys ADD COLUMN current_uses INT NOT NULL DEFAULT 0`)
	DB.Exec(`ALTER TABLE lty_cdkeys ADD COLUMN expires_at DATETIME NULL`)
	DB.Exec(`ALTER TABLE lty_model_upstreams ADD COLUMN channel_name VARCHAR(255) DEFAULT ''`)
	DB.Exec(`ALTER TABLE lty_point_records ADD COLUMN title VARCHAR(255) DEFAULT ''`)
	
	DB.Exec(`ALTER TABLE lty_inspirations ADD COLUMN image_url VARCHAR(1000) DEFAULT ''`)
	DB.Exec(`ALTER TABLE lty_inspirations ADD COLUMN main_category VARCHAR(100) DEFAULT ''`)
	DB.Exec(`ALTER TABLE lty_inspirations ADD COLUMN sub_category VARCHAR(100) DEFAULT ''`)
	DB.Exec(`ALTER TABLE lty_inspirations DROP COLUMN title`)
	DB.Exec(`ALTER TABLE lty_inspirations DROP COLUMN sort_order`)
	DB.Exec(`ALTER TABLE lty_inspirations DROP COLUMN is_active`)
	
	// 扩展积分字段类型支持更大的数值
	DB.Exec(`ALTER TABLE lty_users MODIFY COLUMN points DECIMAL(20,4) NOT NULL DEFAULT 0.0000`)
	DB.Exec(`ALTER TABLE lty_point_records MODIFY COLUMN points_change DECIMAL(20,4) NOT NULL`)
	
	// 优化任务列表加载速度，添加索引
	DB.Exec(`CREATE INDEX idx_created_at ON lty_tasks(created_at)`)
	DB.Exec(`CREATE INDEX idx_user_created ON lty_tasks(user_id, created_at)`)
	DB.Exec(`CREATE INDEX idx_status_created ON lty_tasks(status, created_at)`)
	DB.Exec(`CREATE INDEX idx_session_created ON lty_tasks(session_id, created_at)`)
	DB.Exec(`CREATE INDEX idx_user_session_created ON lty_tasks(user_id, session_id, created_at)`)
	DB.Exec(`CREATE INDEX idx_project_id ON lty_tasks(project_id)`)
	DB.Exec(`CREATE INDEX idx_series_model ON lty_model_upstreams(series_id, logical_model)`)

	return nil
}

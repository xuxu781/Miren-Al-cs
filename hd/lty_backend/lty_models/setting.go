package lty_models

import (
	"database/sql"
	"lty_backend/lty_config"
	"time"
)

type Setting struct {
	ID          int       `json:"id"`
	KeyName     string    `json:"key_name"`
	KeyValue    string    `json:"key_value"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// GetSetting retrieves a setting by its key
func GetSetting(keyName string) (*Setting, error) {
	query := "SELECT id, key_name, key_value, description, created_at, updated_at FROM lty_settings WHERE key_name = ?"
	row := lty_config.DB.QueryRow(query, keyName)

	var s Setting
	err := row.Scan(&s.ID, &s.KeyName, &s.KeyValue, &s.Description, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Not found
		}
		return nil, err
	}

	return &s, nil
}

// GetSettings retrieves multiple settings
func GetSettings() ([]Setting, error) {
	query := "SELECT id, key_name, key_value, description, created_at, updated_at FROM lty_settings"
	rows, err := lty_config.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var settings []Setting
	for rows.Next() {
		var s Setting
		err := rows.Scan(&s.ID, &s.KeyName, &s.KeyValue, &s.Description, &s.CreatedAt, &s.UpdatedAt)
		if err != nil {
			return nil, err
		}
		settings = append(settings, s)
	}
	return settings, nil
}

// SaveSetting updates or creates a setting
func SaveSetting(keyName, keyValue, description string) error {
	existing, err := GetSetting(keyName)
	if err != nil {
		return err
	}

	if existing != nil {
		// Update
		query := "UPDATE lty_settings SET key_value = ?, description = ? WHERE key_name = ?"
		_, err := lty_config.DB.Exec(query, keyValue, description, keyName)
		return err
	}

	// Insert
	query := "INSERT INTO lty_settings (key_name, key_value, description) VALUES (?, ?, ?)"
	_, err = lty_config.DB.Exec(query, keyName, keyValue, description)
	return err
}

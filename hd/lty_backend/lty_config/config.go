package lty_config

import (
	"encoding/json"
	"os"
)

type AppConfig struct {
	DBHost string `json:"db_host"`
	DBPort string `json:"db_port"`
	DBUser string `json:"db_user"`
	DBPass string `json:"db_pass"`
	DBName string `json:"db_name"`
}

var Config AppConfig
var IsInstalled bool

const configFilePath = "config.json"

func LoadConfig() {
	data, err := os.ReadFile(configFilePath)
	if err == nil {
		if err := json.Unmarshal(data, &Config); err == nil {
			IsInstalled = true
			return
		}
	}
	IsInstalled = false
}

func SaveConfig(cfg AppConfig) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	err = os.WriteFile(configFilePath, data, 0644)
	if err == nil {
		Config = cfg
		IsInstalled = true
	}
	return err
}

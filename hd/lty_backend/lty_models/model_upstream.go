package lty_models

import (
	"time"
)

type ModelUpstream struct {
	ID                int       `json:"id"`
	ModelType         int       `json:"model_type"`    // 1: text, 2: video, 4: image
	SeriesID          string    `json:"series_id"`     // 系列ID，提供给用户选择
	LogicalModel      string    `json:"logical_model"` // e.g., gpt-4, flux
	Provider          string    `json:"provider"`      // e.g., openai, custom
	ConnectionURL     string    `json:"connection_url"`
	APIKey            string    `json:"api_key"`
	ChannelName       string    `json:"channel_name"`       // 渠道名称
	ResolutionTiers   string    `json:"resolution_tiers"`   // e.g. "1K,2K,4K"
	AspectRatios      string    `json:"aspect_ratios"`      // e.g. "1:1,4:3,16:9"
	ImageCounts       string    `json:"image_counts"`       // e.g. "1,2,4"
	ActivityTag       string    `json:"activity_tag"`       // 自定义活动标签
	ActivityTagColor  string    `json:"activity_tag_color"` // 活动标签颜色
	SizeParameter     string    `json:"size_parameter"`
	MaxReferenceImages int      `json:"max_reference_images"`
	ResolutionConfigs string    `json:"resolution_configs"` // JSON string for tier-specific configs
	TimeoutSeconds    int       `json:"timeout_seconds"`
	Capabilities      string    `json:"capabilities"`       // JSON string
	BillingStrategy   string    `json:"billing_strategy"`    // JSON string
	ExecutionStrategy string    `json:"execution_strategy"`  // JSON string
	IsPrimary         bool      `json:"is_primary"`
	Status            string    `json:"status"` // active, inactive
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

package lty_models

import (
	"time"
)

type Inspiration struct {
	ID           uint      `json:"id"`
	Content      string    `json:"content"`
	ImageURL     string    `json:"image_url"`
	MainCategory string    `json:"main_category"`
	SubCategory  string    `json:"sub_category"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

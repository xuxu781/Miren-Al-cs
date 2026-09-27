package lty_models

import (
	"time"
)

type Cdkey struct {
	ID            uint       `json:"id"`
	Cdkey         string     `json:"cdkey"`
	Points        float64    `json:"points"`
	Status        int        `json:"status"` // 0: unused/active, 1: exhausted
	MaxUses       int        `json:"max_uses"` // 0 means unlimited
	CurrentUses   int        `json:"current_uses"`
	UsedByUserID  *uint      `json:"used_by_user_id"`
	ExpiresAt     *time.Time `json:"expires_at"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	UsedAt        *time.Time `json:"used_at"`
}

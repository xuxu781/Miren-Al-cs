package lty_utils

import (
	"crypto/sha256"
	"encoding/hex"
)

// HashPassword 使用标准库 sha256 实现密码哈希
func HashPassword(password string) string {
	h := sha256.New()
	h.Write([]byte(password))
	return hex.EncodeToString(h.Sum(nil))
}

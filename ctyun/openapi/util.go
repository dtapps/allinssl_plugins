package openapi

import (
	"crypto/hmac"
	"crypto/sha256"
)

// HMAC-SHA256 工具函数
func hmacSha256(key, data []byte) []byte {
	h := hmac.New(sha256.New, key)
	h.Write(data)
	return h.Sum(nil)
}

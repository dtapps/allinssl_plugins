package openapi

import (
	"crypto/hmac"
	"crypto/sha256"
)

// 生成 HMAC-SHA256 签名
func hmacSha256(key, data []byte) []byte {
	h := hmac.New(sha256.New, key)
	h.Write(data)
	return h.Sum(nil)
}

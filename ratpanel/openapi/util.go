package openapi

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/url"
	"strings"
)

func ensureAPIPath(baseURL string) (string, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return "", fmt.Errorf("invalid base URL: %w", err)
	}

	// 确保 Path 不为空
	if u.Path == "" {
		u.Path = "/"
	}

	// 如果 Path 不是以 /api 结尾，则加上
	if !strings.HasSuffix(u.Path, "/api") {
		if !strings.HasSuffix(u.Path, "/") {
			u.Path += "/"
		}
		u.Path += "api"
	}

	return u.String(), nil
}

// 生成 MD5 字符串
func md5String(s string) string {
	h := md5.New()
	io.WriteString(h, s)
	return fmt.Sprintf("%x", h.Sum(nil))
}

// 生成 SHA256 字符串
func sha256String(str string) string {
	sum := sha256.Sum256([]byte(str))
	dst := make([]byte, hex.EncodedLen(len(sum)))
	hex.Encode(dst, sum[:])
	return string(dst)
}

// 生成 HMAC-SHA256 签名
func hmacSha256(data string, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

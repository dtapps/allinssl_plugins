package openapi

import (
	"crypto/md5"
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

	// 如果 Path 不是以 /openApi 结尾，则加上
	if !strings.HasSuffix(u.Path, "/openApi") {
		if !strings.HasSuffix(u.Path, "/") {
			u.Path += "/"
		}
		u.Path += "openApi"
	}

	return u.String(), nil
}

// 生成 MD5 字符串
func md5String(s string) string {
	h := md5.New()
	io.WriteString(h, s)
	return fmt.Sprintf("%x", h.Sum(nil))
}

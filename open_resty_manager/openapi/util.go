package openapi

import (
	"fmt"
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

	// 如果 Path 不是以 /api/v1/ 结尾，则加上
	if !strings.HasSuffix(u.Path, "/api/v1/") {
		if !strings.HasSuffix(u.Path, "/") {
			u.Path += "/"
		}
		u.Path += "api/v1/"
	}

	return u.String(), nil
}

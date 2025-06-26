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

	// 如果 Path 不是以 /webapi/entry.cgi 结尾，则加上
	if !strings.HasSuffix(u.Path, "/webapi/entry.cgi") {
		if !strings.HasSuffix(u.Path, "/") {
			u.Path += "/"
		}
		u.Path += "webapi/entry.cgi"
	}

	return u.String(), nil
}

// WithDebug 设置令牌
func (c *Client) WithToken() *Client {
	c.SetHeader("X-SYNO-TOKEN", c.token)
	return c
}

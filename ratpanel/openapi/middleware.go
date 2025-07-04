package openapi

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"resty.dev/v3"
)

// PreRequestMiddleware 构造请求前的中间件
func PreRequestMiddleware(userID int, userToken string) resty.RequestMiddleware {
	return func(c *resty.Client, r *resty.Request) error {

		// 解析 URL
		parsedURL, err := url.Parse(r.URL)
		if err != nil {
			return fmt.Errorf("解析 URL 失败: %v", err)
		}

		// 获取当前时间戳（秒级）
		timestamp := time.Now().Unix()

		// 获取请求方式
		method := r.Method

		// 提取规范化路径（以 /api 开头）
		canonicalPath := parsedURL.Path
		if !strings.HasPrefix(canonicalPath, "/api") {
			idx := strings.Index(canonicalPath, "/api")
			if idx != -1 {
				canonicalPath = canonicalPath[idx:]
			} else {
				return fmt.Errorf("路径中不包含 /api 前缀")
			}
		}

		// 获取查询参数
		queryParams := r.QueryParams

		// 获取 body
		payloadStr := ""
		if r.Body != nil {
			switch b := r.Body.(type) {
			case string:
				payloadStr = b
			case []byte:
				payloadStr = string(b)
			}
		}

		canonicalRequest := fmt.Sprintf("%s\n%s\n%s\n%s",
			method,
			canonicalPath,
			queryParams.Encode(),
			sha256String(payloadStr))

		// 构造 sigture 字符串
		stringToSign := fmt.Sprintf("%s\n%d\n%s",
			"HMAC-SHA256",
			timestamp,
			sha256String(canonicalRequest))

		// 生成签名
		signature := hmacSha256(stringToSign, userToken)

		// 设置请求头
		r.SetHeader("X-Timestamp", fmt.Sprintf("%d", timestamp))
		r.SetHeader("Authorization", fmt.Sprintf("HMAC-SHA256 Credential=%d, Signature=%s", userID, signature))

		return nil
	}
}

// Ensure2xxResponseMiddleware 确保响应状态码为 2xx
func Ensure2xxResponseMiddleware(_ *resty.Client, resp *resty.Response) error {
	if !resp.IsSuccess() {
		return fmt.Errorf("请求失败: 状态码 %d, 响应: %s", resp.StatusCode(), resp.String())
	}
	return nil
}

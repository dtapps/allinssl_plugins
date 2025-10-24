package openapi

import (
	"encoding/hex"
	"fmt"
	"net/url"
	"strings"
	"time"

	"resty.dev/v3"
)

// PreRequestMiddleware 构造请求前的中间件
func PreRequestMiddleware(secret_id string, secret_key string) resty.RequestMiddleware {
	return func(c *resty.Client, r *resty.Request) error {

		// 取 endpoint/action/version/region
		endpoint := r.Header.Get("X-TC-Endpoint")
		action := r.Header.Get("X-TC-Action")
		version := r.Header.Get("X-TC-Version")
		region := r.Header.Get("X-TC-Region") // 可为空

		// 校验 endpoint/action/version 是否为空
		if endpoint == "" || action == "" || version == "" {
			return fmt.Errorf("X-TC-Endpoint, X-TC-Action and X-TC-Version required")
		}

		// 清空头
		r.Header.Del("X-TC-Endpoint")

		// 解析 service
		u, err := url.Parse(endpoint)
		if err != nil {
			return fmt.Errorf("invalid endpoint: %w", err)
		}
		host := u.Host
		parts := strings.Split(host, ".")
		if len(parts) < 1 {
			return fmt.Errorf("invalid host: %s", host)
		}
		service := parts[0]

		// ************* 步骤 1：拼接规范请求串 *************
		timestamp := time.Now().Unix()
		r.Header.Set("X-TC-Timestamp", fmt.Sprintf("%d", timestamp))

		payload := ""
		if r.Body != nil {
			if b, ok := r.Body.(string); ok {
				payload = b
			} else if b, ok := r.Body.([]byte); ok {
				payload = string(b)
			} else {
				// Resty 会在 SetBody 时支持 map/struct → json，这里最好用 SetBody(string)
				// 否则需要序列化
			}
		}
		hashedRequestPayload := sha256hex(payload)

		canonicalHeaders := fmt.Sprintf("content-type:application/json; charset=utf-8\nhost:%s\nx-tc-action:%s\n",
			host, strings.ToLower(action))
		signedHeaders := "content-type;host;x-tc-action"
		canonicalRequest := fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s",
			"POST",
			"/",
			"",
			canonicalHeaders,
			signedHeaders,
			hashedRequestPayload)

		// ************* 步骤 2：拼接待签名字符串 *************
		date := time.Unix(timestamp, 0).UTC().Format("2006-01-02")
		credentialScope := fmt.Sprintf("%s/%s/tc3_request", date, service)
		hashedCanonicalRequest := sha256hex(canonicalRequest)
		string2sign := fmt.Sprintf("%s\n%d\n%s\n%s",
			"TC3-HMAC-SHA256",
			timestamp,
			credentialScope,
			hashedCanonicalRequest)

		// ************* 步骤 3：计算签名 *************
		secretDate := hmacsha256(date, []byte("TC3"+secret_key))
		secretService := hmacsha256(service, secretDate)
		secretSigning := hmacsha256("tc3_request", secretService)
		signature := hex.EncodeToString(hmacsha256(string2sign, secretSigning))

		// ************* 步骤 4：拼接 Authorization *************
		authorization := fmt.Sprintf("%s Credential=%s/%s, SignedHeaders=%s, Signature=%s",
			"TC3-HMAC-SHA256",
			secret_id,
			credentialScope,
			signedHeaders,
			signature)

		// 设置头
		r.Header.Set("Authorization", authorization)
		r.Header.Set("Host", host)
		r.Header.Set("Content-Type", "application/json; charset=utf-8")

		if region != "" {
			r.Header.Set("X-TC-Region", region)
		}

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

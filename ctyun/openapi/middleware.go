package openapi

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"resty.dev/v3"
)

// PreRequestMiddleware 构造请求前的中间件
func PreRequestMiddleware(accessKeyId, secretAccessKey string) resty.RequestMiddleware {
	return func(c *resty.Client, r *resty.Request) error {

		// 获取当前时间（UTC）
		now := time.Now().UTC()
		eopDate := now.Format("20060102T150405Z")
		dateOnly := now.Format("20060102")

		// 生成流水号（32位随机数）
		reqID := uuid.New().String()

		// 获取查询参数
		queryParams := r.QueryParams

		// 拼接 query 字符串：按 key 排序后拼接 key=value&...
		keys := make([]string, 0, len(queryParams))
		for k := range queryParams {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		var queryPairs []string
		for _, k := range keys {
			for _, v := range queryParams[k] {
				queryPairs = append(queryPairs, fmt.Sprintf("%s=%s", url.QueryEscape(k), url.QueryEscape(v)))
			}
		}
		queryStr := strings.Join(queryPairs, "&")

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

		// 计算 body SHA256 Hex
		payloadHash := sha256.Sum256([]byte(payloadStr))
		payloadHashHex := hex.EncodeToString(payloadHash[:])

		// 需要签名的 Headers
		headersToSign := map[string]string{
			"ctyun-eop-request-id": reqID,
			"eop-date":             eopDate,
		}

		// 按 header 名称排序
		var sortedHeaders []string
		for k := range headersToSign {
			sortedHeaders = append(sortedHeaders, k)
		}
		sort.Strings(sortedHeaders)

		// 构造签名字符串中的 header 部分
		var sb strings.Builder
		for _, k := range sortedHeaders {
			sb.WriteString(fmt.Sprintf("%s:%s\n", k, headersToSign[k]))
		}
		signedHeaders := sb.String()
		headersListStr := strings.Join(sortedHeaders, ";") // 用于 Authorization 头部

		// 构造 sigture 字符串
		sigture := signedHeaders + "\n" + queryStr + "\n" + payloadHashHex

		// 构造签名密钥
		kTime := hmacSha256([]byte(secretAccessKey), []byte(eopDate))
		kAk := hmacSha256(kTime, []byte(accessKeyId))
		kDate := hmacSha256(kAk, []byte(dateOnly))

		// 生成签名
		signature := hmacSha256(kDate, []byte(sigture))
		signStr := base64.StdEncoding.EncodeToString(signature)

		// 设置请求头
		r.SetHeader("ctyun-eop-request-id", reqID)
		r.SetHeader("eop-date", eopDate)
		r.SetHeader("eop-authorization", fmt.Sprintf("%s Headers=%s Signature=%s", accessKeyId, headersListStr, signStr))

		return nil
	}
}

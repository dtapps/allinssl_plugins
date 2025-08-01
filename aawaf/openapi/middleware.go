package openapi

import (
	"fmt"
	"time"

	"resty.dev/v3"
)

// PreRequestMiddleware 构造请求前的中间件
func PreRequestMiddleware(apiKey string) resty.RequestMiddleware {
	return func(c *resty.Client, r *resty.Request) error {

		// 获取当前时间戳（秒级）
		timestamp := time.Now().Unix()

		// 1. 计算 api_sk 的 md5
		md5AK := md5String(apiKey)

		// 2. 拼接 timestamp 和 md5(api_sk)，再做一次 md5
		token := md5String(fmt.Sprintf("%d%s", timestamp, md5AK))

		// 设置请求头
		r.SetHeader("waf_request_time", fmt.Sprintf("%d", timestamp))
		r.SetHeader("waf_request_token", token)

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

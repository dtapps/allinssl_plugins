package openapi

import (
	"fmt"

	"resty.dev/v3"
)

// PreRequestMiddleware 构造请求前的中间件
func PreRequestMiddleware(apiToken string) resty.RequestMiddleware {
	return func(c *resty.Client, r *resty.Request) error {

		// 设置请求头
		r.SetHeader("X-SLCE-API-TOKEN", apiToken)

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

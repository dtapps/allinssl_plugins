package openapi

import (
	"fmt"

	"resty.dev/v3"
)

// Ensure2xxResponseMiddleware 构造请求后的中间件
func Ensure2xxResponseMiddleware(_ *resty.Client, resp *resty.Response) error {
	if !resp.IsSuccess() {
		return fmt.Errorf("请求失败: 状态码 %d, 响应: %s", resp.StatusCode(), resp.String())
	}
	return nil
}

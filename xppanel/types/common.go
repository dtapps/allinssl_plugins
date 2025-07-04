package types

const CodeSuccess = 1000

type CommonResponse[T any] struct {
	Code    int    `json:"code"`           // 状态码
	Data    T      `json:"data,omitempty"` // 返回数据（具体类型由 T 决定）
	Message string `json:"message"`        // 时间戳
}

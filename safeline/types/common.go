package types

type CommonResponse[T any] struct {
	Data T      `json:"data,omitempty"` // 返回数据（具体类型由 T 决定）
	Err  string `json:"err"`            // 错误
	Msg  string `json:"msg"`            // 信息
}

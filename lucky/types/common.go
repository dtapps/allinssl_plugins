package types

const RetSuccess = 0

type CommonResponse[T any] struct {
	Ret  int    `json:"ret"`            // 状态码
	Msg  string `json:"msg,omitempty"`  // 描述
	Data T      `json:"data,omitempty"` // 返回数据（具体类型由 T 决定）
}

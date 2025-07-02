package types

type CommonResponse[T any] struct {
	Code  int   `json:"code"`          // 状态码
	Res   T     `json:"res,omitempty"` // 返回数据（具体类型由 T 决定）
	Nonce int64 `json:"nonce"`         // 时间戳
}

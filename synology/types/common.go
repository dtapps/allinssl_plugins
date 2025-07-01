package types

type CommonResponse[T any] struct {
	Data    T    `json:"data,omitempty"` // 返回数据（具体类型由 T 决定）
	Success bool `json:"success"`        // 是否成功
}

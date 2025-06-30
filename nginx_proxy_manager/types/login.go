package types

// LoginResponse 登录响应
type LoginResponse struct {
	Token   string `json:"token"`   // 登录令牌
	Expires string `json:"expires"` // 令牌过期时间
}

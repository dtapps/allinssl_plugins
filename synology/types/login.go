package types

// 登录
// https://pve.proxmox.com/pve-docs/api-viewer/#/nodes/{node}/certificates/info
type LoginResponse struct {
	Data struct {
		// Did          string `json:"did"`            // 会话ID
		// IsPortalPort bool   `json:"is_portal_port"` // 是否通过公网端口访问
		// Sid          string `json:"sid"`            // 设备唯一标识
		Synotoken string `json:"synotoken"` // Token
	} `json:"data"` // 数据
	Success bool `json:"success"` // 是否成功
}

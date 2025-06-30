package types

// 获取证书列表 响应参数
// https://pve.proxmox.com/pve-docs/api-viewer/#/nodes/{node}/certificates/info
type CertificateListResponse struct {
	Data []struct {
		San         []string `json:"san"`         // 证书域名
		Fingerprint string   `json:"fingerprint"` // 当前证书的指纹
		NotAfter    int64    `json:"notafter"`    // 证书过期时间（Unix 时间戳）
	} `json:"data"`
}

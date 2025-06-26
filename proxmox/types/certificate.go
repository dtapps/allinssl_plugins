package types

// 获取证书列表 响应参数
// https://pve.proxmox.com/pve-docs/api-viewer/#/nodes/{node}/certificates/info
type CertificateListResponse struct {
	Data []struct {
		Fingerprint string `json:"fingerprint"` // 当前证书的指纹
	} `json:"data"`
}

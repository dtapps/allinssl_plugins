package types

// 获取证书列表 响应参数
// https://pve.proxmox.com/pve-docs/api-viewer/#/nodes/{node}/certificates/info
type CertificateListResponse struct {
	Data struct {
		Certificates []struct {
			Desc string `json:"desc"` // 证书描述
			ID   string `json:"id"`   // 证书ID
		} `json:"certificates"`
	} `json:"data"`
	Success bool `json:"success"`
}

package types

// 获取证书列表 响应参数
type ProxyHostListResponse struct {
	ID            int      `json:"id"`             // 域名ID
	DomainNames   []string `json:"domain_names"`   // 域名列表
	CertificateID int      `json:"certificate_id"` // 证书ID
}

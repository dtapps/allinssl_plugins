package types

// 获取证书列表 响应参数
type CertificateListResponse struct {
	ID          int      `json:"id"`           // 证书ID
	NiceName    string   `json:"nice_name"`    // 证书名称
	DomainNames []string `json:"domain_names"` // 域名列表
	ExpiresOn   string   `json:"expires_on"`   // 过期时间
	Meta        struct {
		Certificate    string `json:"certificate"`     // 证书
		CertificateKey string `json:"certificate_key"` // 证书私钥
	} `json:"meta,omitempty"`
}

// 创建证书 响应参数
type CertificateCreateResponse struct {
	ID int `json:"id"` // 证书ID
}

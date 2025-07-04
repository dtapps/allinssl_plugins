package types

// 获取域名列表 响应参数
type SiteQueryDomainListResponse struct {
	Items []struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		HTTPS bool   `json:"https"`
	} `json:"items"` // 网站列表
	Total int `json:"total"` // 总数
}

// 获取域名信息 响应参数
type SiteQueryDomainInfoResponse struct {
	ID                int      `json:"id"`
	Name              string   `json:"name"`
	Type              string   `json:"type"`
	Domains           []string `json:"domains"`
	HTTPS             bool     `json:"https"`
	SSLCertificate    string   `json:"ssl_certificate"`
	SSLCertificateKey string   `json:"ssl_certificate_key"`
	SSLNotBefore      string   `json:"ssl_not_before"`
	SSLNotAfter       string   `json:"ssl_not_after"`
	SSLDNSNames       any      `json:"ssl_dns_names,omitempty"`
	SSLIssuer         string   `json:"ssl_issuer"`
	SSLOCSPServer     any      `json:"ssl_ocsp_server,omitempty"`
}

package types

// 获取证书列表 响应参数
type CertificateQueryCertListResponse struct {
	Items []struct {
		ID       int      `json:"id"`
		Domains  []string `json:"domains"`
		DnsNames []string `json:"dns_names"`
		Cert     string   `json:"cert"`
		Key      string   `json:"key"`
	} `json:"items"` // 网站列表
	Total int `json:"total"` // 总数
}

// 上传证书 响应参数
type CertificateUpdateCertInfoResponse struct {
	ID int `json:"id"`
}

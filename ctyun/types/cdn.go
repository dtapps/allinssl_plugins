package types

// CDN加速 查询域名配置信息 响应参数
type CdnQueryDomainInfoResponse struct {
	Domain   string `json:"domain"`    // 域名
	CertName string `json:"cert_name"` // 证书备注名
}

// CDN加速 修改域名配置 请求参数
type CdnUpdateDomainInfoRequest struct {
	Domain      string `json:"domain"`              // 域名
	ProductCode string `json:"product_code"`        // 产品类型
	CertName    string `json:"cert_name,omitempty"` // 证书备注名
}

// CDN加速 查询证书详情 响应参数
type CdnQueryCertInfoResponse struct {
	Result struct {
		ID   int    `json:"id"`   // 证书id
		Name string `json:"name"` // 证书备注名称
	} `json:"result"` // 证书信息
}

// CDN加速 创建证书 请求参数
type CdnUpdateCertInfoRequest struct {
	Name  string `json:"name"`  // 证书备注名
	Key   string `json:"key"`   // 证书私钥
	Certs string `json:"certs"` // 证书公钥
}

// CDN加速 创建证书 响应参数
type CdnUpdateCertInfoResponse struct {
	ID int `json:"id"` // 证书id
}

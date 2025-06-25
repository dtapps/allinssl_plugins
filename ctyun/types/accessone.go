package types

// 边缘安全加速平台 域名基础及加速配置查询 响应参数
type AccessoneQueryDomainInfoResponse struct {
	Domain   string `json:"domain"`    // 域名
	CertName string `json:"cert_name"` // 证书备注名
}

// 边缘安全加速平台 域名基础及加速配置修改 请求参数
type AccessoneUpdateDomainInfoRequest struct {
	Domain      string `json:"domain"`              // 域名
	ProductCode string `json:"product_code"`        // 产品类型
	CertName    string `json:"cert_name,omitempty"` // 证书备注名
}

// 边缘安全加速平台 查询证书详情 响应参数
type AccessoneQueryCertInfoResponse struct {
	ID   int    `json:"id"`   // 证书id
	Name string `json:"name"` // 证书备注名称
}

// 边缘安全加速平台 创建证书 请求参数
type AccessoneUpdateCertInfoRequest struct {
	Name  string `json:"name"`  // 证书备注名
	Key   string `json:"key"`   // 证书私钥
	Certs string `json:"certs"` // 证书公钥
}

// 边缘安全加速平台 创建证书 响应参数
type AccessoneUpdateCertInfoResponse struct {
	ID int `json:"id"` // 证书id
}

package types

// 边缘安全加速平台 域名基础及加速配置查询 响应参数
// https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=113&api=13412&data=174&isNormal=1&vid=167
type AccessoneQueryDomainInfoResponse struct {
	Domain   string `json:"domain"`    // 域名
	CertName string `json:"cert_name"` // 证书备注名
}

// 边缘安全加速平台 查询证书详情 响应参数
// https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=113&api=13015&data=174&isNormal=1&vid=167
type AccessoneQueryCertInfoResponse struct {
	Result struct {
		ID    int    `json:"id"`    // 证书id
		Name  string `json:"name"`  // 证书备注名称
		Certs string `json:"certs"` // 证书内容
		Key   string `json:"key"`   // 私钥内容
	} `json:"result"` // 证书信息
}

// 边缘安全加速平台 创建证书 响应参数
// https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=113&api=13014&data=174&isNormal=1&vid=167
type AccessoneUpdateCertInfoResponse struct {
	ID int `json:"id"` // 证书id
}

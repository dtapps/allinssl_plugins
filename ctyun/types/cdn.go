package types

// CDN加速 查询域名配置信息 响应参数
// https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=112&api=10849&data=173&isNormal=1&vid=166
// https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=108&api=11304&data=161&isNormal=1&vid=154
type CdnQueryDomainInfoResponse struct {
	Domain   string `json:"domain"`    // 域名
	CertName string `json:"cert_name"` // 证书备注名
}

// CDN加速 查询证书详情 响应参数
// https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=112&api=10837&data=173&isNormal=1&vid=166
// https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=108&api=10899&data=161&isNormal=1&vid=154
type CdnQueryCertInfoResponse struct {
	Result struct {
		ID   int    `json:"id"`   // 证书id
		Name string `json:"name"` // 证书备注名称
		Cert string `json:"cert"` // 证书内容
		Key  string `json:"key"`  // 私钥内容
	} `json:"result"` // 证书信息
}

// CDN加速 创建证书 响应参数
// https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=112&api=10835&data=173&isNormal=1&vid=166
// https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=108&api=10893&data=161&isNormal=1&vid=154
type CdnUpdateCertInfoResponse struct {
	ID int `json:"id"` // 证书id
}

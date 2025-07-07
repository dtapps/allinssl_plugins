package types

import "encoding/json"

// 获取证书列表 响应参数
type CertificateListResponse struct {
	ID      int    `json:"id"`      // 证书ID
	Name    string `json:"name"`    // 证书名称
	Domains string `json:"domains"` // 域名列表
	Crt     string `json:"crt"`     // 证书
	Key     string `json:"key"`     // 证书私钥
}

// ParseDomains 将 Domains 字符串解析为字符串切片
func (c CertificateListResponse) ParseDomains() []string {
	var domains []string
	_ = json.Unmarshal([]byte(c.Domains), &domains)
	return domains
}

package types

import "encoding/json"

// 获取证书列表 响应参数
type CertificateListV6Response struct {
	ID   int    `json:"id"`   // 证书ID
	Name string `json:"name"` // 证书名称
	Sni  string `json:"sni"`  // 域名列表
	Cert string `json:"cert"` // 证书
	Key  string `json:"key"`  // 证书私钥
}

// ParseSni 将 Sni 字符串解析为字符串切片
func (c CertificateListV6Response) ParseSni() []string {
	var domains []string
	_ = json.Unmarshal([]byte(c.Sni), &domains)
	return domains
}

// 获取证书列表 响应参数
type CertificateListV7Response struct {
	ID   int    `json:"id"`   // 证书ID
	Name string `json:"name"` // 证书名称
	Sni  string `json:"sni"`  // 域名列表
	Crt  string `json:"crt"`  // 证书
	Key  string `json:"key"`  // 证书私钥
}

// ParseSni 将 Sni 字符串解析为字符串切片
func (c CertificateListV7Response) ParseSni() []string {
	var domains []string
	_ = json.Unmarshal([]byte(c.Sni), &domains)
	return domains
}

package types

import "encoding/json"

// 获取证书列表 响应参数
type SiteListResponse struct {
	// CertOptions []struct {
	// 	Label string `json:"label"` // 证书选项
	// 	Value int    `json:"value"` // 证书选项值
	// } `json:"cert_options"` // 证书选项
	Sites []SiteListSitesResponse `json:"sites"` // 站点列表
}

type SiteListSitesResponse struct {
	CertID    int    `json:"cert_id"`   // 证书ID
	Domains   string `json:"domains"`   // 域名列表
	ForceSsl  bool   `json:"force_ssl"` // 是否强制HTTPS
	Hsts      bool   `json:"hsts"`      // 是否开启HSTS
	Http2     bool   `json:"http2"`     // 是否开启HTTP/2
	ID        int    `json:"id"`        // 域名ID
	Ipv6      bool   `json:"ipv6"`      // 是否开启IPv6
	Listeners string `json:"listeners"` // 监听端口
	Locations string `json:"locations"` // 路径列表
	Name      string `json:"name"`      // 站点名称
}

// ParseDomains 将 Domains 字符串解析为字符串切片
func (c SiteListSitesResponse) ParseDomains() []string {
	var domains []string
	_ = json.Unmarshal([]byte(c.Domains), &domains)
	return domains
}

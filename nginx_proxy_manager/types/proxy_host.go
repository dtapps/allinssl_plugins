package types

// 获取证书列表 响应参数
type ProxyHostListResponse struct {
	ID          int      `json:"id"`           // 域名ID
	DomainNames []string `json:"domain_names"` // 域名列表
	Certificate struct {
		ID       int    `json:"id"`        // 证书ID
		Nickname string `json:"nice_name"` // 证书名称
	} `json:"certificate,omitempty"`
}

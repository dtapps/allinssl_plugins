package types

type SiteQueryDomainListResponse struct {
	List []struct {
		SiteID string `json:"site_id"` // 网站ID
		Server struct {
			ServerName []string `json:"server_name"` // 网站域名
			Ssl        struct {
				FullChain  string `json:"full_chain"`  // SSL证书完整链
				PrivateKey string `json:"private_key"` // SSL证书私钥
			} `json:"ssl"` // SSL证书信息
		} `json:"server"` // 服务器
	} `json:"list"` // 域名列表
	Total int `json:"total"` // 总数
}

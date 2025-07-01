package types

// 获取证书列表 响应参数
type CertificateListResponse struct {
	Certificates []struct {
		ID      string `json:"id"`   // 证书ID
		Desc    string `json:"desc"` // 证书描述
		Subject struct {
			SubAltName []string `json:"sub_alt_name"` // 证书域名
		} `json:"subject"`
		ValidTill string `json:"valid_till"` // 证书过期时间
	} `json:"certificates"`
}

type CertificateUpdateResponse struct {
	ID string `json:"id"`
}

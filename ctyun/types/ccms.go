package types

// 证书管理服务 查询用户证书列表 请求参数
type CcmsQueryCertListRequest struct {
	PageSize int    `json:"pageSize"` // 当前页码
	PageNum  int    `json:"pageNum"`  // 每页记录数
	Origin   string `json:"origin"`   // 证书来源
}

// 证书管理服务 查询用户证书列表 响应参数
type CcmsCertificateListResponse struct {
	TotalSize int                               `json:"totalSize"` // 证书总数量
	List      []CcmsCertificateListResponseList `json:"list"`      // 证书列表
}

// 证书管理服务 查询用户证书列表 响应参数 证书列表
type CcmsCertificateListResponseList struct {
	ID   string `json:"id"`   // 证书id
	Name string `json:"name"` // 证书名称
}

// 证书管理服务 上传证书 请求参数
type CcmsUpdateCertInfoRequest struct {
	Name               string `json:"name"`               // 证书名称
	Certificate        string `json:"certificate"`        // 证书字符串
	PrivateKey         string `json:"privateKey"`         // 私钥字符串
	EncryptionStandard string `json:"encryptionStandard"` // 加密标准
}

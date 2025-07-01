package types

// 证书管理服务 查询用户证书列表 响应参数
// https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=152&api=17233&data=204&isNormal=1&vid=283
type CcmsCertificateListResponse struct {
	TotalSize int                               `json:"totalSize"` // 证书总数量
	List      []CcmsCertificateListResponseList `json:"list"`      // 证书列表
}

// 证书管理服务 查询用户证书列表 响应参数 证书列表
type CcmsCertificateListResponseList struct {
	ID          string `json:"id"`          // 证书id
	Name        string `json:"name"`        // 证书名称
	Fingerprint string `json:"fingerprint"` // 证书指纹
	ExpireTime  string `json:"expireTime"`  // 证书过期时间
}

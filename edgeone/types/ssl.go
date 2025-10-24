package types

// 获取证书列表 响应参数
// https://cloud.tencent.com/document/product/400/41671
type SslDescribeCertificatesResponse struct {
	Response struct {
		Error ErrorResponse `json:"Error,omitempty"` // 公共错误，没有时表示成功

		TotalCount   int64 `json:"TotalCount"` // 总数量
		Certificates []struct {
			Domain         string   `json:"Domain"`         // 主域名
			Alias          string   `json:"Alias"`          // 备注名称
			Status         int64    `json:"Status"`         // 证书状态：0 = 审核中，1 = 已通过，2 = 审核失败，3 = 已过期，4 = 自动添加DNS记录，5 = 企业证书，待提交资料，6 = 订单取消中，7 = 已取消，8 = 已提交资料， 待上传确认函，9 = 证书吊销中，10 = 已吊销，11 = 重颁发中，12 = 待上传吊销确认函，13 = 免费证书待提交资料。14 = 证书已退款。 15 = 证书迁移中
			CERTBeginTime  string   `json:"CertBeginTime"`  // 证书生效时间
			CERTEndTime    string   `json:"CertEndTime"`    // 证书过期时间
			CertificateID  string   `json:"CertificateId"`  // 证书 ID
			SubjectAltName []string `json:"SubjectAltName"` // 证书包含的多个域名（包含主域名）
		} `json:"Certificates"` // 列表
		RequestID string `json:"RequestId"`
	} `json:"Response"`
}

// 上传证书 响应参数
// https://cloud.tencent.com/document/product/400/41665
type SslUploadCertificateResponse struct {
	Response struct {
		Error ErrorResponse `json:"Error,omitempty"` // 公共错误，没有时表示成功

		CertificateID string `json:"CertificateId"` // 证书 ID
		RequestID     string `json:"RequestId"`
	} `json:"Response"`
}

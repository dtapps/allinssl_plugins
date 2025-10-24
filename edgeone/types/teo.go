package types

// TeoDescribeAccelerationDomainsResponse 查询加速域名列表 响应参数
// https://cloud.tencent.com/document/product/1552/86336
type TeoDescribeAccelerationDomainsResponse struct {
	Response struct {
		Error ErrorResponse `json:"Error,omitempty"` // 公共错误，没有时表示成功

		TotalCount          int64 `json:"TotalCount"` // 总数量
		AccelerationDomains []struct {
			DomainName   string `json:"DomainName"`   // 加速域名名称
			DomainStatus string `json:"DomainStatus"` // 加速域名状态 online：已生效；process：部署中；offline：已停用；forbidden：已封禁；init：未生效，待激活站点
			Certificate  struct {
				Mode string `json:"Mode"` // 配置服务端证书的模式
				List []struct {
					CertID     string `json:"CertId"`     // 证书 ID
					Alias      string `json:"Alias"`      // 备注名称
					ExpireTime string `json:"ExpireTime"` // 证书过期时间
					Status     string `json:"Status"`     // 证书状态：deployed：已部署；processing：部署中；applying：申请中；failed：申请失败；issued：绑定失败
				} `json:"List"`
			} `json:"Certificate"` // 加速域名所对应的证书信息
		} `json:"AccelerationDomains"`
		RequestId string `json:"RequestId"`
	} `json:"Response"`
}

// TeoModifyHostsCertificateResponse 配置域名证书 响应参数
// https://cloud.tencent.com/document/product/1552/80764
type TeoModifyHostsCertificateResponse struct {
	Response struct {
		Error ErrorResponse `json:"Error,omitempty"` // 公共错误，没有时表示成功

		RequestId string `json:"RequestId"`
	} `json:"Response"`
}

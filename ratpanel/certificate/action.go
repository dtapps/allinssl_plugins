package certificate

import (
	"fmt"

	"github.com/dtapps/allinssl_plugins/core"
	"github.com/dtapps/allinssl_plugins/ratpanel/openapi"
	"github.com/dtapps/allinssl_plugins/ratpanel/types"
)

// 上传证书
// certID: 证书ID
// isExist: 是否存在
func Action(openapiClient *openapi.Client, certBundle *core.CertBundle) (certID int, isExist bool, err error) {

	// 1. 获取证书列表
	var queryCertListResp types.CommonResponse[types.CertificateQueryCertListResponse]
	_, err = openapiClient.R().
		SetQueryParams(map[string]string{
			"page":  "1",   // 页码
			"limit": "100", // 每页数量
		}).
		SetResult(&queryCertListResp).
		SetContentType("application/json").
		SetForceResponseContentType("application/json").
		Get("/cert/cert")
	if err != nil {
		return 0, false, fmt.Errorf("获取证书列表错误: %w", err)
	}

	// 2. 遍历证书列表，查找是否存在相同的证书
	for _, item := range queryCertListResp.Data.Items {
		if item.Cert != "" && item.Key != "" {
			// 获取接口证书信息
			apiCertBundle, err := core.ParseCertBundle([]byte(item.Cert), []byte(item.Key))
			if err != nil {
				return 0, false, fmt.Errorf("解析接口证书信息错误: %w", err)
			}
			// 如果接口证书没有过期就对比是否与传入的证书信息一致
			if !apiCertBundle.IsExpired() {
				if apiCertBundle.GetFingerprintSHA256() == certBundle.GetFingerprintSHA256() {
					// 域名已经绑定了证书
					return item.ID, true, nil
				}
			}
		}
	}

	// 3. 上传证书
	var updateCertInfoResp types.CommonResponse[types.CertificateUpdateCertInfoResponse]
	_, err = openapiClient.R().
		SetBodyMap(map[string]any{
			"cert": certBundle.Certificate, // 证书完整链
			"key":  certBundle.PrivateKey,  // 证书私钥
		}).
		SetResult(&updateCertInfoResp).
		SetContentType("application/json").
		SetForceResponseContentType("application/json").
		Post("/cert/cert/upload")
	if err != nil {
		return 0, false, fmt.Errorf("上传证书错误: %w", err)
	}

	return updateCertInfoResp.Data.ID, false, nil
}

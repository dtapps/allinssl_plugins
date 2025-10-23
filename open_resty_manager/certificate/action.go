package certificate

import (
	"fmt"

	"github.com/dtapps/allinssl_plugins/core"
	"github.com/dtapps/allinssl_plugins/open_resty_manager/openapi"
	"github.com/dtapps/allinssl_plugins/open_resty_manager/types"
)

// 上传证书
// isExist: 是否存在
func Action(openapiClient *openapi.Client, certBundle *core.CertBundle) (isExist bool, err error) {

	// 1. 获取证书列表
	var certListResp []types.CertificateListResponse
	_, err = openapiClient.R().
		SetResult(&certListResp).
		SetContentType("application/json").
		Get("/admin/certs")
	if err != nil {
		return false, fmt.Errorf("获取证书列表错误: %w", err)
	}
	for _, certInfo := range certListResp {
		if certInfo.Crt != "" && certInfo.Key != "" && certInfo.Domains != "" {
			if certBundle.IsDNSNamesMatch(certInfo.ParseDomains()) {
				// 获取接口证书信息
				apiCertBundle, err := core.ParseCertBundle([]byte(certInfo.Crt), []byte(certInfo.Key))
				if err != nil {
					return false, fmt.Errorf("解析接口证书信息错误: %w", err)
				}
				// 如果接口证书没有过期就对比是否与传入的证书信息一致
				if !apiCertBundle.IsExpired() {
					if apiCertBundle.GetFingerprintSHA256() == certBundle.GetFingerprintSHA256() {
						// 证书已存在且未过期
						return true, nil
					}
				}
			}
		}
	}

	// 2. 上传证书
	_, err = openapiClient.R().
		SetBody(map[string]any{
			"type": 1,
			"name": certBundle.GetNoteShort(),
			"crt":  certBundle.Certificate,
			"key":  certBundle.PrivateKey,
		}).
		SetContentType("application/json").
		Post("/admin/certs")
	if err != nil {
		return false, fmt.Errorf("创建证书错误: %w", err)
	}

	return false, nil
}

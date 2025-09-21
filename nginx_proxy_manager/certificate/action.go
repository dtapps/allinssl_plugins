package certificate

import (
	"fmt"
	"strings"
	"time"

	"github.com/dtapps/allinssl_plugins/core"
	"github.com/dtapps/allinssl_plugins/nginx_proxy_manager/openapi"
	"github.com/dtapps/allinssl_plugins/nginx_proxy_manager/types"
)

// 上传证书
// certID: 证书ID
// isExist: 是否存在
func Action(openapiClient *openapi.Client, certBundle *core.CertBundle) (certID int, isExist bool, err error) {

	// 1. 获取证书列表
	var certListResp []types.CertificateListResponse
	_, err = openapiClient.R().
		SetResult(&certListResp).
		Get("/nginx/certificates")
	if err != nil {
		return 0, false, fmt.Errorf("获取证书列表错误: %w", err)
	}
	for _, certInfo := range certListResp {
		certID = certInfo.ID
		if strings.EqualFold(certInfo.NiceName, certBundle.GetNote()) {
			if certInfo.Meta.Certificate != "" && certInfo.Meta.CertificateKey != "" {
				if certBundle.IsDNSNamesMatch(certInfo.DomainNames) {
					var expiresOn time.Time
					expiresOn, err = time.Parse(time.DateTime, certInfo.ExpiresOn)
					if err != nil {
						return 0, false, fmt.Errorf("解析过期时间失败: %w", err)
					}
					if expiresOn.After(time.Now()) {
						// 证书已存在且未过期
						return certInfo.ID, true, nil
					}
				} else {
					// 证书已存在
					return certInfo.ID, true, nil
				}
			} else {
				certID = 0
			}
		} else {
			certID = 0
		}
	}

	// 2. 创建证书
	if certID == 0 {
		var certCreateResp types.CertificateCreateResponse
		_, err = openapiClient.R().
			SetBody(map[string]string{
				"provider":  "other",
				"nice_name": certBundle.GetNote(),
			}).
			SetResult(&certCreateResp).
			Post("/nginx/certificates")
		if err != nil {
			return 0, false, fmt.Errorf("创建证书错误: %w", err)
		}
		certID = certCreateResp.ID
	}

	// 3. 上传证书
	_, err = openapiClient.R().
		SetFileReader("certificate", "certificate.pem", strings.NewReader(certBundle.Certificate)).
		SetFileReader("certificate_key", "private_key.pem", strings.NewReader(certBundle.PrivateKey)).
		SetFileReader("intermediate_certificate", "ca_bundle.pem", strings.NewReader(certBundle.CertificateChain)).
		SetPathParams(map[string]string{
			"certID": fmt.Sprintf("%d", certID),
		}).
		Post("/nginx/certificates/{certID}/upload")
	if err != nil {
		err = fmt.Errorf("上传证书错误: %w", err)
		return
	}

	return certID, false, nil
}

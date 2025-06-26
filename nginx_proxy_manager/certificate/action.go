package certificate

import (
	"fmt"
	"strings"

	"github.com/dtapps/allinssl_plugins/nginx_proxy_manager/core"
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
		SetContentType("application/json").
		SetResult(&certListResp).
		Get("/nginx/certificates")
	if err != nil {
		err = fmt.Errorf("获取证书列表错误: %w", err)
		return
	}
	for _, cert := range certListResp {
		if cert.NiceName == certBundle.GetNote() {
			if cert.Meta.Certificate != "" && cert.Meta.CertificateKey != "" {
				// 证书已存在
				return cert.ID, true, nil
			} else {
				certID = cert.ID
			}
		}
	}

	// 2. 创建证书
	if certID == 0 {
		var certCreateResp types.CertificateCreateResponse
		_, err = openapiClient.R().
			SetContentType("application/json").
			SetBody(map[string]string{
				"provider":  "other",
				"nice_name": certBundle.GetNote(),
			}).
			SetResult(&certCreateResp).
			Post("/nginx/certificates")
		if err != nil {
			err = fmt.Errorf("创建证书错误: %w", err)
			return
		}
		certID = certCreateResp.ID
	}

	// 3. 上传证书
	_, err = openapiClient.R().
		SetFileReader("certificate", "certificate.pem", strings.NewReader(certBundle.Certificate)).
		SetFileReader("certificate_key", "private_key.pem", strings.NewReader(certBundle.PrivateKey)).
		Post(fmt.Sprintf("/nginx/certificates/%v/upload", certID))
	if err != nil {
		err = fmt.Errorf("上传证书错误: %w", err)
		return
	}

	return certID, false, nil
}

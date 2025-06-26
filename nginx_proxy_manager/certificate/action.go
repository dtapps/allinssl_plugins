package certificate

import (
	"fmt"
	"strings"

	"github.com/dtapps/allinssl_plugins/nginx_proxy_manager/core"
	"github.com/dtapps/allinssl_plugins/nginx_proxy_manager/openapi"
)

// 上传证书
// certID: 证书ID
func Action(openapiClient *openapi.Client, certBundle *core.CertBundle) (certID int, err error) {

	// 1. 获取证书列表
	var certListResp []struct {
		ID       int    `json:"id"`        // 证书ID
		NiceName string `json:"nice_name"` // 证书名称
		Meta     struct {
			Certificate    string `json:"certificate"`
			CertificateKey string `json:"certificate_key"`
		} `json:"meta,omitempty"`
	}
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
				return cert.ID, nil
			} else {
				certID = cert.ID
			}
		}
	}

	// 2. 创建证书
	if certID == 0 {
		var certCreateResp struct {
			ID int `json:"id"` // 唯一标识符
		}
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

	return certID, nil
}

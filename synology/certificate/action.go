package certificate

import (
	"fmt"
	"strings"

	"github.com/dtapps/allinssl_plugins/synology/core"
	"github.com/dtapps/allinssl_plugins/synology/openapi"
	"github.com/dtapps/allinssl_plugins/synology/types"
)

// 上传证书
// isExist: 是否存在
func Action(openapiClient *openapi.Client, certBundle *core.CertBundle) (isExist bool, err error) {

	// 1. 获取证书列表
	var certListResp types.CertificateListResponse
	_, err = openapiClient.R().
		SetQueryParam("api", "SYNO.Core.Certificate.CRT").
		SetQueryParam("version", "1").
		SetQueryParam("method", "list").
		SetResult(&certListResp).
		Get("")
	if err != nil {
		err = fmt.Errorf("获取证书列表错误: %w", err)
		return
	}
	for _, certInfo := range certListResp.Data.Certificates {
		if strings.EqualFold(certInfo.Desc, certBundle.GetNote()) {
			// 证书已存在
			return true, nil
		}
	}

	// 2. 上传证书
	_, err = openapiClient.R().
		SetQueryParam("api", "SYNO.Core.Certificate").
		SetQueryParam("version", "1").
		SetQueryParam("method", "import").
		SetFormData(map[string]string{
			"desc": certBundle.GetNote(),
		}).
		SetFileReader("cert", "certificate.pem", strings.NewReader(certBundle.Certificate)).
		SetFileReader("key", "private_key.pem", strings.NewReader(certBundle.PrivateKey)).
		Post("")
	if err != nil {
		err = fmt.Errorf("上传证书错误: %w", err)
		return
	}

	return false, nil
}

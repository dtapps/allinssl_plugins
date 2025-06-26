package certificate

import (
	"fmt"
	"strings"

	"github.com/dtapps/allinssl_plugins/proxmox/core"
	"github.com/dtapps/allinssl_plugins/proxmox/openapi"
)

// 上传证书
// pveNode: 证书备注名
func Action(openapiClient *openapi.Client, pveNode string, certBundle *core.CertBundle) (err error) {

	// 1. 获取证书列表
	var certListResp struct {
		Data []struct {
			Fingerprint string `json:"fingerprint"` // 当前证书的指纹
		} `json:"data"`
	}
	_, err = openapiClient.R().
		SetContentType("application/json").
		SetResult(&certListResp).
		Get(fmt.Sprintf("/api2/json/nodes/%s/certificates/info", pveNode))
	if err != nil {
		err = fmt.Errorf("获取证书列表错误: %w", err)
		return
	}
	fmt.Println("数据", certListResp)
	for _, certInfo := range certListResp.Data {
		if strings.EqualFold(strings.ReplaceAll(certInfo.Fingerprint, ":", ""), certBundle.FingerprintSHA1) {
			// 证书已存在
			return nil
		}
		if strings.EqualFold(strings.ReplaceAll(certInfo.Fingerprint, ":", ""), certBundle.FingerprintSHA256) {
			// 证书已存在
			return nil
		}
	}

	// 2. 上传证书
	_, err = openapiClient.R().
		SetContentType("application/json").
		SetBody(map[string]string{
			"certificates": certBundle.Certificate,
			"key":          certBundle.PrivateKey,
		}).
		Post(fmt.Sprintf("/api2/json/nodes/%s/certificates/custom", pveNode))
	if err != nil {
		err = fmt.Errorf("上传证书错误: %w", err)
		return
	}

	return nil
}

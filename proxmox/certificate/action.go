package certificate

import (
	"fmt"
	"strings"
	"time"

	"github.com/dtapps/allinssl_plugins/core"
	"github.com/dtapps/allinssl_plugins/proxmox/openapi"
	"github.com/dtapps/allinssl_plugins/proxmox/types"
)

// 上传证书
// pveNode: 证书备注名
// isExist: 是否存在
func Action(openapiClient *openapi.Client, pveNode string, certBundle *core.CertBundle) (isExist bool, err error) {

	// 1. 获取证书列表
	var certListResp types.CertificateListResponse
	_, err = openapiClient.R().
		SetContentType("application/json").
		SetResult(&certListResp).
		Get(fmt.Sprintf("/api2/json/nodes/%s/certificates/info", pveNode))
	if err != nil {
		return false, fmt.Errorf("获取证书列表错误: %w", err)
	}
	for _, certInfo := range certListResp.Data {
		apiFingerprint := strings.ReplaceAll(certInfo.Fingerprint, ":", "")
		if strings.EqualFold(apiFingerprint, certBundle.GetFingerprintSHA256()) {
			if len(certInfo.San) > 0 && len(certBundle.DNSNames) > 0 && core.EqualStringSlices(certInfo.San, certBundle.DNSNames) {
				notAfter := time.Unix(certInfo.NotAfter, 0)
				if notAfter.After(time.Now()) {
					// 证书已存在且未过期
					return true, nil
				}
			} else {
				// 证书已存在
				return true, nil
			}
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
		return false, fmt.Errorf("上传证书错误: %w", err)
	}

	return false, nil
}

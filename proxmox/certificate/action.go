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
// pveNode: 节点
// isExist: 是否存在
func Action(openapiClient *openapi.Client, pveNode string, certBundle *core.CertBundle) (isExist bool, err error) {

	// 1. 获取证书列表
	// https://pve.proxmox.com/pve-docs/api-viewer/#/nodes/{node}/certificates/info
	var certListResp types.CertificateListResponse
	_, err = openapiClient.R().
		SetPathParams(map[string]string{
			"node": pveNode,
		}).
		SetResult(&certListResp).
		SetContentType("application/json").
		Get("/api2/json/nodes/{node}/certificates/info")
	if err != nil {
		return false, fmt.Errorf("获取证书列表错误: %w", err)
	}
	for _, certInfo := range certListResp.Data {
		apiFingerprint := strings.ReplaceAll(certInfo.Fingerprint, ":", "")
		if strings.EqualFold(apiFingerprint, certBundle.GetFingerprintSHA256()) {
			if certBundle.IsDNSNamesMatch(certInfo.San) {
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
	// https://pve.proxmox.com/pve-docs/api-viewer/#/nodes/{node}/certificates/custom
	_, err = openapiClient.R().
		SetBody(map[string]string{
			"certificates": certBundle.Certificate,
			"key":          certBundle.PrivateKey,
		}).
		SetPathParams(map[string]string{
			"node": pveNode,
		}).
		SetContentType("application/json").
		Post("/api2/json/nodes/{node}/certificates/custom")
	if err != nil {
		return false, fmt.Errorf("上传证书错误: %w", err)
	}

	return false, nil
}

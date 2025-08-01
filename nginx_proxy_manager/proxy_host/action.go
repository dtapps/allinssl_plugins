package proxy_host

import (
	"fmt"
	"strings"

	"github.com/dtapps/allinssl_plugins/core"
	"github.com/dtapps/allinssl_plugins/nginx_proxy_manager/openapi"
	"github.com/dtapps/allinssl_plugins/nginx_proxy_manager/types"
)

// 域名绑定证书
// domain: 域名
// certID: 证书ID
// hostID: 域名ID
func Action(openapiClient *openapi.Client, domain string, certID int, certBundle *core.CertBundle) (hostID int, err error) {

	// 1. 获取域名列表
	var proxyHostsListResp []types.ProxyHostListResponse
	_, err = openapiClient.R().
		SetResult(&proxyHostsListResp).
		SetContentType("application/json").
		Get("/nginx/proxy-hosts")
	if err != nil {
		return 0, fmt.Errorf("获取域名列表错误: %w", err)
	}
	for _, item := range proxyHostsListResp {
		if strings.Contains(strings.Join(item.DomainNames, ","), domain) {
			if item.CertificateID == certID {
				// 证书已绑定域名
				return item.ID, nil
			} else {
				hostID = item.ID
			}
		}
	}
	if hostID == 0 {
		return 0, fmt.Errorf("域名 %s 不存在", domain)
	}

	// 2. 绑定证书
	_, err = openapiClient.R().
		SetBody(map[string]any{
			"certificate_id": certID,
		}).
		SetPathParams(map[string]string{
			"hostID": fmt.Sprintf("%d", hostID),
		}).
		SetContentType("application/json").
		Put("/nginx/proxy-hosts/{hostID}")
	if err != nil {
		return hostID, fmt.Errorf("域名 %s 绑定证书错误: %w", domain, err)
	}

	return hostID, nil
}

package proxy_host

import (
	"fmt"
	"strings"

	"github.com/dtapps/allinssl_plugins/nginx_proxy_manager/core"
	"github.com/dtapps/allinssl_plugins/nginx_proxy_manager/openapi"
)

// 域名绑定证书
// domain: 域名
// note: 证书备注名
// hostID: 域名ID
func Action(openapiClient *openapi.Client, domain string, certID int, certBundle *core.CertBundle) (hostID int, err error) {

	// 1. 获取域名列表
	var proxyHostsListResp []struct {
		ID          int      `json:"id"`           // 域名ID
		DomainNames []string `json:"domain_names"` // 域名列表
		Certificate struct {
			ID       int    `json:"id"`        // 证书ID
			Nickname string `json:"nice_name"` // 证书名称
		} `json:"certificate,omitempty"`
	}
	_, err = openapiClient.R().
		SetContentType("application/json").
		SetResult(&proxyHostsListResp).
		Get("/nginx/proxy-hosts")
	if err != nil {
		err = fmt.Errorf("获取域名列表错误: %w", err)
		return
	}
	for _, item := range proxyHostsListResp {
		if strings.Contains(strings.Join(item.DomainNames, ","), domain) {
			if item.Certificate.Nickname == certBundle.GetNote() {
				// 证书已绑定域名
				return item.ID, nil
			} else {
				hostID = item.ID
			}
		} else {
			err = fmt.Errorf("域名不存在")
			return
		}
	}

	// 2. 绑定证书
	_, err = openapiClient.R().
		SetContentType("application/json").
		SetBody(map[string]any{
			"certificate_id": certID,
		}).
		Put(fmt.Sprintf("/nginx/proxy-hosts/%v", hostID))
	if err != nil {
		err = fmt.Errorf("绑定证书错误: %w", err)
		return
	}

	return hostID, nil
}

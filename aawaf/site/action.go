package site

import (
	"fmt"

	"github.com/dtapps/allinssl_plugins/aawaf/openapi"
	"github.com/dtapps/allinssl_plugins/aawaf/types"
	"github.com/dtapps/allinssl_plugins/core"
)

// 域名绑定证书
// domain: 域名
// isBind: 是否已绑定
func Action(openapiClient *openapi.Client, domain string, certBundle *core.CertBundle) (isBound bool, err error) {

	// 1. 获取域名列表
	var queryDomainListResp types.CommonResponse[types.SiteQueryDomainListResponse]
	_, err = openapiClient.R().
		SetBody(map[string]any{
			"p":         1,      // 页码
			"p_size":    10,     // 每页数量
			"site_name": domain, // 域名
		}).
		SetResult(&queryDomainListResp).
		SetContentType("application/json").
		SetForceResponseContentType("application/json").
		Post("/wafmastersite/get_site_list")
	if err != nil {
		return false, fmt.Errorf("获取域名列表错误: %w", err)
	}

	// 2. 检查域名是否配置了现存证书
	var siteID string
	for _, item := range queryDomainListResp.Res.List {
		// 检查域名是否包含传入的域名
		for _, serverName := range item.Server.ServerName {
			if serverName == domain {
				siteID = item.SiteID // 找到匹配的域名，记录站点ID
				break
			}
		}
		if siteID != "" && item.Server.Ssl.FullChain != "" && item.Server.Ssl.PrivateKey != "" {
			// 获取接口证书信息
			apiCertBundle, err := core.ParseCertBundle([]byte(item.Server.Ssl.FullChain), []byte(item.Server.Ssl.PrivateKey))
			if err != nil {
				return false, fmt.Errorf("解析接口证书信息错误: %w", err)
			}
			// 如果接口证书没有过期就对比是否与传入的证书信息一致
			if !apiCertBundle.IsExpired() {
				if apiCertBundle.GetFingerprintSHA256() == certBundle.GetFingerprintSHA256() {
					// 域名已经绑定了证书
					return true, nil
				}
			}
		}
	}

	// 3. 更新域名信息
	var updateDomainInfoResp types.CommonResponse[string]
	_, err = openapiClient.R().
		SetBody(map[string]any{
			"site_ids":    []string{siteID},       // 域名
			"full_chain":  certBundle.Certificate, // 证书完整链
			"private_key": certBundle.PrivateKey,  // 证书私钥
			"ssl_name":    certBundle.GetNote(),   // 证书备注名
		}).
		SetResult(&updateDomainInfoResp).
		SetContentType("application/json").
		SetForceResponseContentType("application/json").
		Post("/wafmastersite/deploy_ssl")
	if err != nil {
		return false, fmt.Errorf("更新域名信息错误: %w", err)
	}

	return false, nil
}

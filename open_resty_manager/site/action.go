package site

import (
	"fmt"
	"strings"

	"github.com/dtapps/allinssl_plugins/core"
	"github.com/dtapps/allinssl_plugins/open_resty_manager/openapi"
	"github.com/dtapps/allinssl_plugins/open_resty_manager/types"
)

// 域名绑定证书
// domain: 域名
func Action(openapiClient *openapi.Client, domain string, certBundle *core.CertBundle) (err error) {

	// 1. 获取域名列表
	var siteListResp types.SiteListResponse
	_, err = openapiClient.R().
		SetResult(&siteListResp).
		SetContentType("application/json").
		Get("/admin/sites")
	if err != nil {
		return fmt.Errorf("获取域名列表错误: %w", err)
	}
	var siteInfo types.SiteListSitesResponse
	for _, item := range siteListResp.Sites {
		if strings.Contains(strings.Join(item.ParseDomains(), ","), domain) {
			siteInfo = item
		}
	}
	if siteInfo.ID == 0 {
		return fmt.Errorf("域名 %s 不存在", domain)
	}

	// 2. 获取证书列表
	var certListResp []types.CertificateListResponse
	_, err = openapiClient.R().
		SetResult(&certListResp).
		SetContentType("application/json").
		Get("/admin/certs")
	if err != nil {
		return fmt.Errorf("获取证书列表错误: %w", err)
	}
	for _, certInfo := range certListResp {
		if certInfo.Crt != "" && certInfo.Key != "" && certInfo.Domains != "" {
			if certBundle.IsDNSNamesMatch(certInfo.ParseDomains()) {
				// 获取接口证书信息
				apiCertBundle, err := core.ParseCertBundle([]byte(certInfo.Crt), []byte(certInfo.Key))
				if err != nil {
					return fmt.Errorf("解析接口证书信息错误: %w", err)
				}
				// 如果接口证书没有过期就对比是否与传入的证书信息一致
				if !apiCertBundle.IsExpired() {
					if apiCertBundle.GetFingerprintSHA256() == certBundle.GetFingerprintSHA256() {
						if siteInfo.CertID == certInfo.ID {
							// 证书已绑定
							return nil
						} else {
							siteInfo.CertID = certInfo.ID
						}
					}
				}
			}
		}
	}

	// 3. 更新域名
	_, err = openapiClient.R().
		SetBodyStruct(siteInfo).
		SetContentType("application/json").
		Put("/admin/sites")
	if err != nil {
		return fmt.Errorf("域名 %s 绑定证书错误: %w", domain, err)
	}

	return nil
}

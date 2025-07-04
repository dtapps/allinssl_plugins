package site

import (
	"fmt"

	"github.com/dtapps/allinssl_plugins/core"
	"github.com/dtapps/allinssl_plugins/ratpanel/openapi"
	"github.com/dtapps/allinssl_plugins/ratpanel/types"
)

// 域名绑定证书
// domain: 域名
// certID: 证书ID
// isBind: 是否已绑定
func Action(openapiClient *openapi.Client, domain string, certID int, certBundle *core.CertBundle) (isBound bool, err error) {

	// 1. 获取域名列表
	var queryDomainListResp types.CommonResponse[types.SiteQueryDomainListResponse]
	_, err = openapiClient.R().
		SetQueryParams(map[string]string{
			"page":  "1",   // 页码
			"limit": "100", // 每页数量
		}).
		SetResult(&queryDomainListResp).
		SetContentType("application/json").
		SetForceResponseContentType("application/json").
		Get("/website")
	if err != nil {
		return false, fmt.Errorf("获取域名列表错误: %w", err)
	}

	// 2. 检查域名是否存在
	var siteID int
	for _, item := range queryDomainListResp.Data.Items {
		// 检查域名是否包含传入的域名
		if item.Name == domain {
			siteID = item.ID // 找到匹配的域名，记录站点ID
			break
		}
	}

	// 3. 如果没有找到匹配的域名，返回错误
	if siteID == 0 {
		return false, fmt.Errorf("没有找到匹配的域名: %s", domain)
	}

	// 4. 通过站点ID获取域名详细信息
	var queryDomainInfoResp types.CommonResponse[types.SiteQueryDomainInfoResponse]
	_, err = openapiClient.R().
		SetPathParams(map[string]string{
			"id": fmt.Sprintf("%d", siteID), // 站点ID
		}).
		SetResult(&queryDomainListResp).
		SetContentType("application/json").
		SetForceResponseContentType("application/json").
		Get("/website/{id}")
	if err != nil {
		return false, fmt.Errorf("获取域名列表错误: %w", err)
	}

	// 5. 检查域名是否配置了现存证书
	if queryDomainInfoResp.Data.SSLCertificate != "" && queryDomainInfoResp.Data.SSLCertificateKey != "" {
		// 获取接口证书信息
		apiCertBundle, err := core.ParseCertBundle([]byte(queryDomainInfoResp.Data.SSLCertificate), []byte(queryDomainInfoResp.Data.SSLCertificateKey))
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

	// 6. 部署证书到站点
	var updateDomainInfoResp types.CommonResponse[string]
	_, err = openapiClient.R().
		SetBodyMap(map[string]any{
			"id":         certID, // 证书ID
			"website_id": siteID, // 站点ID
		}).
		SetPathParams(map[string]string{
			"id": fmt.Sprintf("%d", certID), // 证书ID
		}).
		SetResult(&updateDomainInfoResp).
		SetContentType("application/json").
		SetForceResponseContentType("application/json").
		Post("/cert/cert/{id}/deploy")
	if err != nil {
		return false, fmt.Errorf("更新域名信息错误: %w", err)
	}

	return false, nil
}

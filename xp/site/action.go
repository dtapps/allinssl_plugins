package site

import (
	"fmt"

	"github.com/dtapps/allinssl_plugins/core"
	"github.com/dtapps/allinssl_plugins/xp/openapi"
	"github.com/dtapps/allinssl_plugins/xp/types"
)

// 域名绑定证书
// domain: 域名
// isBind: 是否已绑定
func Action(openapiClient *openapi.Client, domain string, certBundle *core.CertBundle) (isBound bool, err error) {

	// 1. 获取域名列表
	var queryDomainListResp types.CommonResponse[[]types.SiteQueryDomainListResponse]
	_, err = openapiClient.R().
		SetBody(map[string]any{
			"p":         1,      // 页码
			"p_size":    10,     // 每页数量
			"site_name": domain, // 域名
		}).
		SetResult(&queryDomainListResp).
		SetContentType("application/json").
		SetForceResponseContentType("application/json").
		Get("/siteList")
	if err != nil {
		return false, fmt.Errorf("获取域名列表错误: %w", err)
	}
	if queryDomainListResp.Code != types.CodeSuccess {
		return false, fmt.Errorf("获取域名列表失败: %s", queryDomainListResp.Message)
	}
	if len(queryDomainListResp.Data) <= 0 {
		return false, fmt.Errorf("没有找到匹配的域名: %s", domain)
	}

	// 2. 检查域名是否配置了现存证书
	var siteID int
	for _, item := range queryDomainListResp.Data {
		// 检查域名是否包含传入的域名
		if item.Name == domain {
			siteID = item.ID // 找到匹配的域名，记录站点ID
			break
		}
	}

	// 3. 更新域名信息
	var updateDomainInfoResp types.CommonResponse[any]
	_, err = openapiClient.R().
		SetBody(map[string]any{
			"id":  siteID,                 // 域名
			"key": certBundle.Certificate, // 证书完整链
			"pem": certBundle.PrivateKey,  // 证书私钥
		}).
		SetResult(&updateDomainInfoResp).
		SetContentType("application/json").
		SetForceResponseContentType("application/json").
		Post("/setSSL")
	if err != nil {
		return false, fmt.Errorf("更新域名信息错误: %w", err)
	}
	if updateDomainInfoResp.Code != types.CodeSuccess {
		return false, fmt.Errorf("更新域名信息失败: %s", updateDomainInfoResp.Message)
	}

	return false, nil
}

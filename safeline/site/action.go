package site

import (
	"fmt"

	"github.com/dtapps/allinssl_plugins/core"
	"github.com/dtapps/allinssl_plugins/safeline/openapi"
	"github.com/dtapps/allinssl_plugins/safeline/types"
)

// 域名绑定证书
// domain: 域名
// isBind: 是否已绑定
func Action(openapiClient *openapi.Client, domain string, certBundle *core.CertBundle) (isBound bool, err error) {

	// 1. 获取域名列表
	var queryDomainListResp types.CommonResponse[types.SiteQueryDomainListResponse]
	_, err = openapiClient.R().
		SetQueryParam("page", "1").        // 页码
		SetQueryParam("page_size", "900"). // 每页数量
		SetQueryParam("site", domain).     // 域名
		SetResult(&queryDomainListResp).
		SetContentType("application/json").
		SetForceResponseContentType("application/json").
		Get("/open/site")
	if err != nil {
		return false, fmt.Errorf("获取域名列表错误: %w", err)
	}

	// 2. 检查域名是否配置了现存证书
	updateDomainInfoReq := types.SiteQueryDomainListDataResponse{}
	for _, item := range queryDomainListResp.Data.Data {
		// 检查域名是否包含传入的域名
		for _, serverName := range item.ServerNames {
			if serverName == domain {
				updateDomainInfoReq = item // 记录域名信息
				break
			}
		}
	}

	// 3. 如果没有找到匹配的域名，说明域名未绑定
	if updateDomainInfoReq.ID == 0 {
		return false, fmt.Errorf("域名 %s 不存在", domain)
	}

	// 4. 如果站点ID和证书ID都存在，说明域名已经绑定了证书，检查证书是否过期
	if updateDomainInfoReq.CertID > 0 {
		// 获取证书信息
		var queryCertInfoResp types.CommonResponse[types.SiteQueryCertInfoResponse]
		_, err = openapiClient.R().
			SetResult(&queryCertInfoResp).
			SetContentType("application/json").
			SetForceResponseContentType("application/json").
			Get(fmt.Sprintf("/open/cert/%d", updateDomainInfoReq.CertID))
		if err != nil {
			return false, fmt.Errorf("获取证书信息错误: %w", err)
		}
		// 获取接口证书信息
		apiCertBundle, err := core.ParseCertBundle([]byte(queryCertInfoResp.Data.Manual.Crt), []byte(queryCertInfoResp.Data.Manual.Key))
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

	// 5. 获取证书列表
	var queryCertListResp types.CommonResponse[types.SiteQueryCertListResponse]
	_, err = openapiClient.R().
		SetResult(&queryCertListResp).
		SetContentType("application/json").
		SetForceResponseContentType("application/json").
		Get("/open/cert")
	if err != nil {
		return false, fmt.Errorf("获取证书列表错误: %w", err)
	}

	// 6. 检查证书是否已经存在
	isCartExist := false
	for _, item := range queryCertListResp.Data.Nodes {
		// 检查证书域名是否包含传入的域名
		if core.CanDomainsUseCert([]string{domain}, item.Domains) {
			// 获取证书信息
			var queryCertInfoResp types.CommonResponse[types.SiteQueryCertInfoResponse]
			_, err = openapiClient.R().
				SetResult(&queryCertInfoResp).
				SetContentType("application/json").
				SetForceResponseContentType("application/json").
				Get(fmt.Sprintf("/open/cert/%d", item.ID))
			if err != nil {
				return false, fmt.Errorf("获取证书信息错误: %w", err)
			}
			// 获取接口证书信息
			apiCertBundle, err := core.ParseCertBundle([]byte(queryCertInfoResp.Data.Manual.Crt), []byte(queryCertInfoResp.Data.Manual.Key))
			if err != nil {
				return false, fmt.Errorf("解析接口证书信息错误: %w", err)
			}
			// 如果接口证书没有过期就对比是否与传入的证书信息一致
			if !apiCertBundle.IsExpired() {
				if apiCertBundle.GetFingerprintSHA256() == certBundle.GetFingerprintSHA256() {
					updateDomainInfoReq.CertID = item.ID // 记录证书ID
					isCartExist = true                   // 证书已存在
				}
			}
		}
	}

	// 7. 上传证书
	if !isCartExist {
		var updateCertInfoResp types.CommonResponse[int]
		_, err = openapiClient.R().
			SetBodyMap(map[string]any{
				"type": 2,
				"manual": map[string]any{
					"crt": certBundle.Certificate, // 证书完整链
					"key": certBundle.PrivateKey,  // 证书私钥
				},
			}).
			SetResult(&updateCertInfoResp).
			SetContentType("application/json").
			SetForceResponseContentType("application/json").
			Post("/open/cert")
		if err != nil {
			return false, fmt.Errorf("上传证书错误: %w", err)
		}
		updateDomainInfoReq.CertID = updateCertInfoResp.Data // 获取新证书ID
	}

	// 8. 更新域名信息
	var updateDomainInfoResp types.CommonResponse[any]
	_, err = openapiClient.R().
		SetBodyStruct(updateDomainInfoReq).
		SetResult(&updateDomainInfoResp).
		SetContentType("application/json").
		SetForceResponseContentType("application/json").
		Put("/open/site")
	if err != nil {
		return false, fmt.Errorf("更新域名信息错误: %w", err)
	}

	return false, nil
}

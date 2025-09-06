package cdn

import (
	"fmt"

	"github.com/dtapps/allinssl_plugins/core"
	"github.com/dtapps/allinssl_plugins/ctyun/openapi"
	"github.com/dtapps/allinssl_plugins/ctyun/types"
)

// 域名绑定证书
// domain: 域名
// isBind: 是否已绑定
func Action(openapiClient *openapi.Client, domain string, certBundle *core.CertBundle) (isBound bool, err error) {

	// 1. 获取域名信息
	var queryDomainInfoResp types.CommonResponse[types.CdnQueryDomainInfoResponse]
	_, err = openapiClient.R().
		SetQueryParam("product_code", productCode). // 产品类型
		SetQueryParam("domain", domain).            // 域名
		SetResult(&queryDomainInfoResp).
		SetContentType("application/json").
		Get("/v1/domain/query-domain-detail")
	if err != nil {
		return false, fmt.Errorf("获取 %s 域名信息错误: %w", domain, err)
	}
	if queryDomainInfoResp.StatusCode != types.StatusCodeSuccess {
		return false, fmt.Errorf("获取 %s 域名信息失败: %s", domain, queryDomainInfoResp.Message)
	}

	// 2. 检查域名是否配置了现存证书
	if queryDomainInfoResp.ReturnObj.CertName == certBundle.GetNoteShort() {
		return true, nil
	}

	// 3. 查询证书信息
	var queryCertInfoResp types.CommonResponse[types.CdnQueryCertInfoResponse]
	_, err = openapiClient.R().
		SetQueryParam("name", certBundle.GetNoteShort()). // 证书备注名
		SetResult(&queryCertInfoResp).
		SetContentType("application/json").
		Get("/v1/cert/query-cert-detail")
	if err != nil {
		return false, fmt.Errorf("查询证书信息错误: %w", err)
	}
	// if queryCertInfoResp.StatusCode != types.StatusCodeSuccess {
	// 	return false, fmt.Errorf("查询证书信息失败: %s", queryCertInfoResp.Message)
	// }

	// 4. 证书不存在就上传证书
	if queryCertInfoResp.ReturnObj.Result.Name != certBundle.GetNoteShort() {
		certificate, privateKey := core.BuildCertsForAPI(certBundle)
		var updateCertInfoResp types.CommonResponse[types.CdnUpdateCertInfoResponse]
		_, err = openapiClient.R().
			SetBodyMap(map[string]any{
				"name":  certBundle.GetNoteShort(), // 证书备注名
				"key":   privateKey,                // 证书私钥
				"certs": certificate,               // 证书公钥
			}).
			SetResult(&updateCertInfoResp).
			SetContentType("application/json").
			Post("/v1/cert/creat-cert")
		if err != nil {
			return false, fmt.Errorf("上传证书错误: %w", err)
		}
		if updateCertInfoResp.StatusCode != types.StatusCodeSuccess {
			return false, fmt.Errorf("上传证书失败: %s", updateCertInfoResp.Message)
		}
	}

	// 5. 更新域名信息
	var updateDomainInfoResp types.CommonResponse[any]
	_, err = openapiClient.R().
		SetBodyMap(map[string]any{
			"domain":       domain,                    // 域名
			"product_code": productCode,               // 产品类型
			"cert_name":    certBundle.GetNoteShort(), // 证书备注名
		}).
		SetResult(&updateDomainInfoResp).
		SetContentType("application/json").
		Post("/v1/domain/update-domain")
	if err != nil {
		return false, fmt.Errorf("更新域名信息错误: %w", err)
	}
	if updateDomainInfoResp.StatusCode != types.StatusCodeSuccess {
		return false, fmt.Errorf("更新域名信息失败: %s", updateDomainInfoResp.Message)
	}

	return false, nil
}

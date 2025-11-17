package cdn

import (
	"context"
	"fmt"

	"github.com/dtapps/allinssl_plugins/core"
	"github.com/dtapps/allinssl_plugins/ctyun/openapi"
	"github.com/dtapps/allinssl_plugins/ctyun/types"
)

// 域名绑定证书
// domain: 域名
// isBind: 是否已绑定
// 查询域名配置信息 https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=108&api=11304&data=161&isNormal=1&vid=154
// 查询证书详情 https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=108&api=10899&data=161&isNormal=1&vid=154
// 创建证书 https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=108&api=10893&data=161&isNormal=1&vid=154
// 修改域名配置 https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=108&api=11308&data=161&isNormal=1&vid=154
func Action(ctx context.Context, openapiClient *openapi.Client, domain string, certBundle *core.CertBundle) (isBound bool, err error) {

	// 1. 获取域名信息
	var queryDomainInfo types.CommonResponse[types.CdnQueryDomainInfoResponse]
	_, err = openapiClient.R().
		SetQueryParam("product_code", productCode). // 产品类型
		SetQueryParam("domain", domain).            // 域名
		SetResult(&queryDomainInfo).
		SetContext(ctx).
		Get("/v1/domain/query-domain-detail")
	if err != nil {
		return false, fmt.Errorf("获取 %s 域名信息错误: %w", domain, err)
	}
	if queryDomainInfo.StatusCode != types.StatusCodeSuccess {
		return false, fmt.Errorf("获取 %s 域名信息失败:%s ; %s", domain, queryDomainInfo.Message, queryDomainInfo.ErrorMessage)
	}

	// 2. 检查域名是否配置了现存证书
	if certBundle.IsSameCertificateNote(certBundle.GetNoteShort(), queryDomainInfo.ReturnObj.CertName) {
		return true, nil
	}

	// 3. 证书不存在就上传证书
	if !certBundle.IsSameCertificateNote(certBundle.GetNoteShort(), queryDomainInfo.ReturnObj.CertName) {

		// 是否上传
		isUpload := false

		// 查询证书是否存在
		var queryCertInfo types.CommonResponse[types.CdnQueryCertInfoResponse]
		_, err = openapiClient.R().
			SetQueryParam("name", certBundle.GetNoteShort()). // 证书备注名
			SetResult(&queryCertInfo).
			SetContext(ctx).
			Get("/v1/cert/query-cert-detail")
		if err != nil {
			return false, fmt.Errorf("查询证书信息错误: %w", err)
		}
		if queryCertInfo.ReturnObj.Result.Name != "" {
			isUpload = false
		} else {
			isUpload = true
		}

		// 不存在就上传
		if isUpload {
			privateKey, certificate := core.BuildCertsForAPI(certBundle)
			var uploadCertInfo types.CommonResponse[types.CdnUpdateCertInfoResponse]
			_, err = openapiClient.R().
				SetBodyMap(map[string]any{
					"name":  certBundle.GetNoteShort(), // 证书备注名
					"key":   privateKey,                // 证书私钥
					"certs": certificate,               // 证书公钥
				}).
				SetResult(&uploadCertInfo).
				SetContext(ctx).
				Post("/v1/cert/creat-cert")
			if err != nil {
				return false, fmt.Errorf("上传证书错误: %w", err)
			}
			if uploadCertInfo.StatusCode != types.StatusCodeSuccess {
				return false, fmt.Errorf("上传证书失败: %s ; %s", uploadCertInfo.Message, uploadCertInfo.ErrorMessage)
			}
		}
	}

	// 4. 更新域名信息
	var updateDomainInfo types.CommonResponse[any]
	_, err = openapiClient.R().
		SetBodyMap(map[string]any{
			"domain":       domain,                    // 域名
			"product_code": productCode,               // 产品类型
			"cert_name":    certBundle.GetNoteShort(), // 证书备注名
		}).
		SetResult(&updateDomainInfo).
		SetContext(ctx).
		Post("/v1/domain/update-domain")
	if err != nil {
		return false, fmt.Errorf("更新域名信息错误: %w", err)
	}
	if updateDomainInfo.StatusCode != types.StatusCodeSuccess {
		return false, fmt.Errorf("更新域名信息失败:%s ; %s", updateDomainInfo.Message, updateDomainInfo.ErrorMessage)
	}

	return false, nil
}

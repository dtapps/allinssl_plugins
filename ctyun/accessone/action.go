package accessone

import (
	"fmt"

	"github.com/dtapps/allinssl_plugins/core"
	"github.com/dtapps/allinssl_plugins/ctyun/openapi"
	"github.com/dtapps/allinssl_plugins/ctyun/types"
)

// 域名绑定证书
// domain: 域名
// isBind: 是否已绑定
// 域名基础及加速配置查询 https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=113&api=13412&data=174&isNormal=1&vid=167
// 查询证书详情 https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=113&api=13015&data=174&isNormal=1&vid=167
// 创建证书 https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=113&api=13014&data=174&isNormal=1&vid=167
// 域名基础及加速配置修改 https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=113&api=13413&data=174&isNormal=1&vid=167
func Action(openapiClient *openapi.Client, domain string, certBundle *core.CertBundle) (isBound bool, err error) {

	// 1. 获取域名信息
	var queryDomainInfo types.CommonResponse[types.AccessoneQueryDomainInfoResponse]
	_, err = openapiClient.R().
		SetBodyMap(map[string]any{
			"product_code": productCode, // 产品类型
			"domain":       domain,      // 域名
		}).
		SetResult(&queryDomainInfo).
		Post("/ctapi/v1/accessone/domain/config")
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

	// 3. 查询证书信息
	var queryCertInfo types.CommonResponse[types.AccessoneQueryCertInfoResponse]
	_, err = openapiClient.R().
		SetQueryParam("name", certBundle.GetNoteShort()). // 证书备注名
		SetResult(&queryCertInfo).
		Get("/ctapi/v1/accessone/cert/query")
	if err != nil {
		return false, fmt.Errorf("查询证书信息错误: %w", err)
	}

	// 4. 证书不存在就上传证书
	if !certBundle.IsSameCertificateNote(certBundle.GetNoteShort(), queryCertInfo.ReturnObj.Name) {
		// 加载 API 证书
		apiCertBundle, err := certBundle.LoadApiCert([]byte(queryCertInfo.ReturnObj.Cert), []byte(queryCertInfo.ReturnObj.Key))
		if err != nil {
			return false, fmt.Errorf("加载 API 证书错误: %w", err)
		}
		// 是否上传
		isUpload := false
		if !apiCertBundle.IsCertsEqual() {
			isUpload = true
		}
		if apiCertBundle.IsAPIExpired() {
			isUpload = true
		}
		if isUpload {
			privateKey, certificate := core.BuildCertsForAPI(apiCertBundle.Local)
			var uploadCertInfo types.CommonResponse[types.AccessoneUpdateCertInfoResponse]
			_, err = openapiClient.R().
				SetBodyMap(map[string]any{
					"name":  certBundle.GetNoteShort(), // 证书备注名
					"key":   privateKey,                // 证书私钥
					"certs": certificate,               // 证书公钥
				}).
				SetResult(&uploadCertInfo).
				Post("/ctapi/v1/accessone/cert/create")
			if err != nil {
				return false, fmt.Errorf("上传证书错误: %w", err)
			}
			if uploadCertInfo.StatusCode != types.StatusCodeSuccess {
				return false, fmt.Errorf("上传证书失败: %s ; %s", uploadCertInfo.Message, uploadCertInfo.ErrorMessage)
			}
		}
	}

	// 5. 更新域名信息
	var updateDomainInfo types.CommonResponse[any]
	_, err = openapiClient.R().
		SetBodyMap(map[string]any{
			"domain":       domain,                    // 域名
			"product_code": productCode,               // 产品类型
			"cert_name":    certBundle.GetNoteShort(), // 证书备注名
		}).
		SetResult(&updateDomainInfo).
		Post("/ctapi/v1/accessone/domain/modify_config")
	if err != nil {
		return false, fmt.Errorf("更新域名信息错误: %w", err)
	}
	if updateDomainInfo.StatusCode != types.StatusCodeSuccess {
		return false, fmt.Errorf("更新域名信息失败:%s ; %s", updateDomainInfo.Message, updateDomainInfo.ErrorMessage)
	}

	return false, nil
}

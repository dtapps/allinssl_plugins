package certificate

import (
	"fmt"
	"strings"

	"github.com/dtapps/allinssl_plugins/core"
	"github.com/dtapps/allinssl_plugins/lucky/openapi"
	"github.com/dtapps/allinssl_plugins/lucky/types"
)

// 上传证书
// isExist: 是否存在
func Action(openapiClient *openapi.Client, certBundle *core.CertBundle) (isExist bool, err error) {

	// 1. 获取证书列表
	var certListResp types.CertificateListResponse
	_, err = openapiClient.R().
		SetResult(&certListResp).
		SetContentType("application/json").
		Get("/ssl")
	if err != nil {
		return false, fmt.Errorf("获取证书列表错误: %w", err)
	}
	if certListResp.Ret != types.RetSuccess {
		return false, fmt.Errorf("获取证书列表错误: %s", certListResp.Msg)
	}
	for _, certInfo := range certListResp.List {
		if strings.EqualFold(certInfo.Remark, certBundle.GetNoteShort()) {
			// 证书已存在
			return true, nil
		}
	}

	// 2. 上传证书公钥
	var certUpdateCertResp types.CertificateUpdateResponse
	_, err = openapiClient.R().
		SetFileReader("file", "certificate.pem", strings.NewReader(certBundle.Certificate)).
		SetResult(&certUpdateCertResp).
		Post("/getfilebase64")
	if err != nil {
		err = fmt.Errorf("上传证书公钥错误: %w", err)
		return
	}
	if certUpdateCertResp.Ret != types.RetSuccess {
		return false, fmt.Errorf("上传证书公钥错误: %s", certUpdateCertResp.Msg)
	}

	// 3. 上传证书私钥
	var certUpdateKeyResp types.CertificateUpdateResponse
	_, err = openapiClient.R().
		SetFileReader("file", "private_key.pem", strings.NewReader(certBundle.PrivateKey)).
		SetResult(&certUpdateKeyResp).
		Post("/getfilebase64")
	if err != nil {
		err = fmt.Errorf("上传证书私钥错误: %w", err)
		return
	}
	if certUpdateKeyResp.Ret != types.RetSuccess {
		return false, fmt.Errorf("上传证书私钥错误: %s", certUpdateKeyResp.Msg)
	}

	// 3. 创建证书
	var certCreateResp types.CertificateCreateResponse
	_, err = openapiClient.R().
		SetBody(map[string]any{
			"AddFrom":    "file",
			"CertBase64": certUpdateCertResp.Base64,
			"KeyBase64":  certUpdateKeyResp.Base64,
			"Remark":     certBundle.GetNoteShort(),
		}).
		SetResult(&certCreateResp).
		SetContentType("application/json").
		Post("/ssl")
	if err != nil {
		return false, fmt.Errorf("创建证书错误: %w", err)
	}
	if certCreateResp.Ret != types.RetSuccess {
		return false, fmt.Errorf("创建证书错误: %s", certCreateResp.Msg)
	}

	return false, nil
}

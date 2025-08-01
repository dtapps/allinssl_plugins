package certificate

import (
	"encoding/json"
	"fmt"

	"github.com/dtapps/allinssl_plugins/core"
	"github.com/dtapps/allinssl_plugins/uuwaf/openapi"
	"github.com/dtapps/allinssl_plugins/uuwaf/types"
)

// 上传证书
// isExist: 是否存在
func V6Action(openapiClient *openapi.Client, certBundle *core.CertBundle) (isExist bool, err error) {

	// 1. 获取证书列表
	var certListResp []types.CertificateListV6Response
	_, err = openapiClient.R().
		SetResult(&certListResp).
		SetContentType("application/json").
		Get("/certs")
	if err != nil {
		return false, fmt.Errorf("获取证书列表错误: %w", err)
	}
	for _, certInfo := range certListResp {
		if certInfo.Cert != "" && certInfo.Key != "" && certInfo.Sni != "" {
			if certBundle.IsDNSNamesMatch(certInfo.ParseSni()) {
				// 获取接口证书信息
				apiCertBundle, err := core.ParseCertBundle([]byte(certInfo.Cert), []byte(certInfo.Key))
				if err != nil {
					return false, fmt.Errorf("解析接口证书信息错误: %w", err)
				}
				// 如果接口证书没有过期就对比是否与传入的证书信息一致
				if !apiCertBundle.IsExpired() {
					if apiCertBundle.GetFingerprintSHA256() == certBundle.GetFingerprintSHA256() {
						// 证书已存在且未过期
						return true, nil
					}
				}
			}
		}
	}

	// 2. 检查证书
	_, err = openapiClient.R().
		SetBody(map[string]any{
			"mode": 0,
			"cert": certBundle.Certificate,
			"key":  certBundle.PrivateKey,
		}).
		SetContentType("application/json").
		Post("/certs/check")
	if err != nil {
		return false, fmt.Errorf("检查证书错误: %w", err)
	}

	// 3. 上传证书
	sniStr, err := json.Marshal(certBundle.DNSNames)
	if err != nil {
		panic(err)
	}
	_, err = openapiClient.R().
		SetBody(map[string]any{
			"cert":        certBundle.Certificate,
			"expire_time": certBundle.NotAfter.Format("2006-01-02 15:04:05"),
			"key":         certBundle.PrivateKey,
			"sni":         string(sniStr),
		}).
		SetContentType("application/json").
		Post("/certs/config")
	if err != nil {
		return false, fmt.Errorf("创建证书错误: %w", err)
	}

	return false, nil
}

// 上传证书
// isExist: 是否存在
func V7Action(openapiClient *openapi.Client, certBundle *core.CertBundle) (isExist bool, err error) {

	// 1. 获取证书列表
	var certListResp []types.CertificateListV7Response
	_, err = openapiClient.R().
		SetResult(&certListResp).
		SetContentType("application/json").
		Get("/certs")
	if err != nil {
		return false, fmt.Errorf("获取证书列表错误: %w", err)
	}
	for _, certInfo := range certListResp {
		if certInfo.Crt != "" && certInfo.Key != "" && certInfo.Sni != "" {
			if certBundle.IsDNSNamesMatch(certInfo.ParseSni()) {
				// 获取接口证书信息
				apiCertBundle, err := core.ParseCertBundle([]byte(certInfo.Crt), []byte(certInfo.Key))
				if err != nil {
					return false, fmt.Errorf("解析接口证书信息错误: %w", err)
				}
				// 如果接口证书没有过期就对比是否与传入的证书信息一致
				if !apiCertBundle.IsExpired() {
					if apiCertBundle.GetFingerprintSHA256() == certBundle.GetFingerprintSHA256() {
						// 证书已存在且未过期
						return true, nil
					}
				}
			}
		}
	}

	// 2. 上传证书
	_, err = openapiClient.R().
		SetBody(map[string]any{
			"type": 1,
			"name": certBundle.GetNoteShort(),
			"crt":  certBundle.Certificate,
			"key":  certBundle.PrivateKey,
		}).
		SetContentType("application/json").
		Post("/certs")
	if err != nil {
		return false, fmt.Errorf("创建证书错误: %w", err)
	}

	return false, nil
}

package main

import (
	"fmt"

	"github.com/dtapps/allinssl_plugins/ctyun/accessone"
	"github.com/dtapps/allinssl_plugins/ctyun/ccms"
	"github.com/dtapps/allinssl_plugins/ctyun/cdn"
	"github.com/dtapps/allinssl_plugins/ctyun/core"
	"github.com/dtapps/allinssl_plugins/ctyun/icdn"
)

// 部署到天翼云CDN加速
func deployCdnAction(cfg map[string]any) (*Response, error) {

	if cfg == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}
	certStr, ok := cfg["cert"].(string)
	if !ok || certStr == "" {
		return nil, fmt.Errorf("cert is required and must be a string")
	}
	keyStr, ok := cfg["key"].(string)
	if !ok || keyStr == "" {
		return nil, fmt.Errorf("key is required and must be a string")
	}
	accessKey, ok := cfg["access_key"].(string)
	if !ok || accessKey == "" {
		return nil, fmt.Errorf("access_key is required and must be a string")
	}
	secretKey, ok := cfg["secret_key"].(string)
	if !ok || secretKey == "" {
		return nil, fmt.Errorf("secret_key is required and must be a string")
	}
	domain, ok := cfg["domain"].(string)
	if !ok || domain == "" {
		return nil, fmt.Errorf("domain is required and must be a string")
	}

	// 1. 解析证书字符串
	certBundle, err := core.ParseCertBundle([]byte(certStr), []byte(keyStr))
	if err != nil {
		return nil, fmt.Errorf("failed to parse cert bundle: %w", err)
	}

	// 2. 计算证书字符串的SHA256值
	sha256, err := GetSHA256(certBundle.Certificate)
	if err != nil {
		return nil, fmt.Errorf("failed to get SHA256 of cert: %w", err)
	}
	note := fmt.Sprintf("allinssl-%s", sha256)

	// 创建CDN加速客户端
	cdnClient, err := cdn.NewClient(accessKey, secretKey)
	if err != nil {
		return nil, fmt.Errorf("创建CDN加速客户端失败: %w", err)
	}

	// 1. 查询域名是否存在和现存的证书信息
	queryDomainInfo, err := cdnClient.GetQueryDomainInfo(domain)
	if err != nil {
		return nil, fmt.Errorf("查询域名是否存在和现存的证书信息失败: %w", err)
	}
	if queryDomainInfo.ReturnObj.Domain == "" {
		return nil, fmt.Errorf("域名不存在")
	}

	// 2. 检查域名是否配置了现存证书
	if queryDomainInfo.ReturnObj.CertName == note {
		return &Response{
			Status:  "success",
			Message: "证书已绑定域名",
			Result: map[string]any{
				"domain": domain,
				"cert":   certBundle,
			},
		}, nil
	}

	// 3. 查询证书是否已存，不存在就上传证书
	queryCertInfo, _ := cdnClient.GetQueryCertInfo(note)
	if queryCertInfo.ReturnObj.Result.Name != note {
		_, err = cdnClient.PostUpdateCertInfo(note, certBundle.PrivateKey, certBundle.Certificate)
		if err != nil {
			return nil, fmt.Errorf("上传证书失败: %w", err)
		}
	}

	// 4. 更新证书到域名
	_, err = cdnClient.PostUpdateDomainInfo(domain, note)
	if err != nil {
		return nil, fmt.Errorf("更新证书到域名失败: %w", err)
	}

	return &Response{
		Status:  "success",
		Message: "更新域名证书成功",
		Result: map[string]any{
			"domain": domain,
			"cert":   certBundle,
		},
	}, nil
}

// 部署到天翼云全站加速
func deployIcdnAction(cfg map[string]any) (*Response, error) {

	if cfg == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}
	certStr, ok := cfg["cert"].(string)
	if !ok || certStr == "" {
		return nil, fmt.Errorf("cert is required and must be a string")
	}
	keyStr, ok := cfg["key"].(string)
	if !ok || keyStr == "" {
		return nil, fmt.Errorf("key is required and must be a string")
	}
	accessKey, ok := cfg["access_key"].(string)
	if !ok || accessKey == "" {
		return nil, fmt.Errorf("access_key is required and must be a string")
	}
	secretKey, ok := cfg["secret_key"].(string)
	if !ok || secretKey == "" {
		return nil, fmt.Errorf("secret_key is required and must be a string")
	}
	domain, ok := cfg["domain"].(string)
	if !ok || domain == "" {
		return nil, fmt.Errorf("domain is required and must be a string")
	}

	// 1. 解析证书字符串
	certBundle, err := core.ParseCertBundle([]byte(certStr), []byte(keyStr))
	if err != nil {
		return nil, fmt.Errorf("failed to parse cert bundle: %w", err)
	}

	// 2. 计算证书字符串的SHA256值
	sha256, err := GetSHA256(certBundle.Certificate)
	if err != nil {
		return nil, fmt.Errorf("failed to get SHA256 of cert: %w", err)
	}
	note := fmt.Sprintf("allinssl-%s", sha256)

	// 创建全站加速客户端
	icdnClient, err := icdn.NewClient(accessKey, secretKey)
	if err != nil {
		return nil, fmt.Errorf("创建全站加速客户端失败: %w", err)
	}

	// 1. 查询域名是否存在和现存的证书信息
	queryDomainInfo, err := icdnClient.GetQueryDomainInfo(domain)
	if err != nil {
		return nil, fmt.Errorf("查询域名是否存在和现存的证书信息失败: %w", err)
	}
	if queryDomainInfo.ReturnObj.Domain == "" {
		return nil, fmt.Errorf("域名不存在")
	}

	// 2. 检查域名是否配置了现存证书
	if queryDomainInfo.ReturnObj.CertName == note {
		return &Response{
			Status:  "success",
			Message: "证书已绑定域名",
			Result: map[string]any{
				"domain": domain,
				"cert":   certBundle,
			},
		}, nil
	}

	// 3. 查询证书是否已存，不存在就上传证书
	queryCertInfo, _ := icdnClient.GetQueryCertInfo(note)
	if queryCertInfo.ReturnObj.Result.Name != note {
		_, err = icdnClient.PostUpdateCertInfo(note, certBundle.PrivateKey, certBundle.Certificate)
		if err != nil {
			return nil, fmt.Errorf("上传证书失败: %w", err)
		}
	}

	// 4. 更新证书到域名
	_, err = icdnClient.PostUpdateDomainInfo(domain, note)
	if err != nil {
		return nil, fmt.Errorf("更新证书到域名失败: %w", err)
	}

	return &Response{
		Status:  "success",
		Message: "更新域名证书成功",
		Result: map[string]any{
			"domain": domain,
			"cert":   certBundle,
		},
	}, nil
}

// 部署到天翼云边缘安全加速平台
func deployAccessoneAction(cfg map[string]any) (*Response, error) {

	if cfg == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}
	certStr, ok := cfg["cert"].(string)
	if !ok || certStr == "" {
		return nil, fmt.Errorf("cert is required and must be a string")
	}
	keyStr, ok := cfg["key"].(string)
	if !ok || keyStr == "" {
		return nil, fmt.Errorf("key is required and must be a string")
	}
	accessKey, ok := cfg["access_key"].(string)
	if !ok || accessKey == "" {
		return nil, fmt.Errorf("access_key is required and must be a string")
	}
	secretKey, ok := cfg["secret_key"].(string)
	if !ok || secretKey == "" {
		return nil, fmt.Errorf("secret_key is required and must be a string")
	}
	domain, ok := cfg["domain"].(string)
	if !ok || domain == "" {
		return nil, fmt.Errorf("domain is required and must be a string")
	}

	// 1. 解析证书字符串
	certBundle, err := core.ParseCertBundle([]byte(certStr), []byte(keyStr))
	if err != nil {
		return nil, fmt.Errorf("failed to parse cert bundle: %w", err)
	}

	// 2. 计算证书字符串的SHA256值
	sha256, err := GetSHA256(certBundle.Certificate)
	if err != nil {
		return nil, fmt.Errorf("failed to get SHA256 of cert: %w", err)
	}
	note := fmt.Sprintf("allinssl-%s", sha256)

	// 创建边缘安全加速平台客户端
	accessoneClient, err := accessone.NewClient(accessKey, secretKey)
	if err != nil {
		return nil, fmt.Errorf("创建边缘安全加速平台客户端失败: %w", err)
	}

	// 1. 查询域名是否存在和现存的证书信息
	queryDomainInfo, err := accessoneClient.GetQueryDomainInfo(domain)
	if err != nil {
		return nil, fmt.Errorf("查询域名是否存在和现存的证书信息失败: %w", err)
	}
	if queryDomainInfo.ReturnObj.Domain == "" {
		return nil, fmt.Errorf("域名不存在")
	}

	// 2. 检查域名是否配置了现存证书
	if queryDomainInfo.ReturnObj.CertName == note {
		return &Response{
			Status:  "success",
			Message: "证书已绑定域名",
			Result: map[string]any{
				"domain": domain,
				"cert":   certBundle,
			},
		}, nil
	}

	// 3. 查询证书是否已存，不存在就上传证书
	queryCertInfo, _ := accessoneClient.GetQueryCertInfo(note)
	if queryCertInfo.ReturnObj.Name != note {
		_, err = accessoneClient.PostUpdateCertInfo(note, certBundle.PrivateKey, certBundle.Certificate)
		if err != nil {
			return nil, fmt.Errorf("上传证书失败: %w", err)
		}
	}

	// 4. 更新证书到域名
	_, err = accessoneClient.PostUpdateDomainInfo(domain, note)
	if err != nil {
		return nil, fmt.Errorf("更新证书到域名失败: %w", err)
	}

	return &Response{
		Status:  "success",
		Message: "更新域名证书成功",
		Result: map[string]any{
			"domain": domain,
			"cert":   certBundle,
		},
	}, nil
}

// 上传证书到天翼云证书管理
func deployCcmsAction(cfg map[string]any) (*Response, error) {

	if cfg == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}
	certStr, ok := cfg["cert"].(string)
	if !ok || certStr == "" {
		return nil, fmt.Errorf("cert is required and must be a string")
	}

	keyStr, ok := cfg["key"].(string)
	if !ok || keyStr == "" {
		return nil, fmt.Errorf("key is required and must be a string")
	}

	accessKey, ok := cfg["access_key"].(string)
	if !ok || accessKey == "" {
		return nil, fmt.Errorf("access_key is required and must be a string")
	}

	secretKey, ok := cfg["secret_key"].(string)
	if !ok || secretKey == "" {
		return nil, fmt.Errorf("secret_key is required and must be a string")
	}

	// 1. 解析证书字符串
	certBundle, err := core.ParseCertBundle([]byte(certStr), []byte(keyStr))
	if err != nil {
		return nil, fmt.Errorf("failed to parse cert bundle: %w", err)
	}

	// 2. 计算证书字符串的SHA256值
	sha256, err := GetSHA256(certBundle.Certificate)
	if err != nil {
		return nil, fmt.Errorf("failed to get SHA256 of cert: %w", err)
	}

	// 3. 取SHA256的前6位作为唯一标识
	sha256Short := sha256[:6]
	name := fmt.Sprintf("allinssl-%s", sha256Short)

	// 创建证书管理服务客户端
	ccmsClient, err := ccms.NewClient(accessKey, secretKey)
	if err != nil {
		return nil, fmt.Errorf("创建证书管理服务客户端错误: %w", err)
	}

	// 1. 查询证书是否存在
	queryCertList, err := ccmsClient.GetQueryCertList()
	if err != nil {
		return nil, fmt.Errorf("查询证书是否存在错误: %w", err)
	}
	for _, certificate := range queryCertList.ReturnObj.List {
		if certificate.Name == name {
			return &Response{
				Status:  "success",
				Message: "证书已存在",
				Result: map[string]any{
					"cert": certBundle,
				},
			}, nil
		}
	}

	// 2. 上传证书
	_, err = ccmsClient.PostUpdateCertInfo(name, certBundle.Certificate, certBundle.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("上传证书错误: %w", err)
	}

	return &Response{
		Status:  "success",
		Message: "上传证书成功",
		Result: map[string]any{
			"cert": certBundle,
		},
	}, nil
}

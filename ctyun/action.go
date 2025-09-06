package main

import (
	"fmt"
	"time"

	"github.com/dtapps/allinssl_plugins/core"
	"github.com/dtapps/allinssl_plugins/ctyun/accessone"
	"github.com/dtapps/allinssl_plugins/ctyun/ccms"
	"github.com/dtapps/allinssl_plugins/ctyun/cdn"
	"github.com/dtapps/allinssl_plugins/ctyun/icdn"
	"github.com/dtapps/allinssl_plugins/ctyun/openapi"
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

	ctAccessKey, ok := cfg["access_key"].(string)
	if !ok || ctAccessKey == "" {
		return nil, fmt.Errorf("access_key is required and must be a string")
	}
	ctSecretKey, ok := cfg["secret_key"].(string)
	if !ok || ctSecretKey == "" {
		return nil, fmt.Errorf("secret_key is required and must be a string")
	}
	ctDomain, ok := cfg["domain"].(string)
	if !ok || ctDomain == "" {
		return nil, fmt.Errorf("domain is required and must be a string")
	}

	// 解析证书字符串
	certBundle, err := core.ParseCertBundle([]byte(certStr), []byte(keyStr))
	if err != nil {
		return nil, fmt.Errorf("failed to parse cert bundle: %w", err)
	}

	// 验证证书链
	var verifyChainText string
	if err := certBundle.VerifyChain(); err != nil {
		verifyChainText = "❌ 证书链不完整或不被信任:" + err.Error()
	} else {
		verifyChainText = "✅ 证书链完整有效"
	}

	// 1. 检查证书是否过期
	if certBundle.IsExpired() {
		return nil, fmt.Errorf("证书已过期 %s", certBundle.NotAfter.Format(time.DateTime))
	}

	// 2. 解析传入域名
	userDomains, isMultiple := core.ParseDomainsFixedSeparator(ctDomain, ",")
	if isMultiple {
		if !certBundle.CanDomainsUseCert(userDomains) {
			return nil, fmt.Errorf("域名和证书不匹配，证书支持域名：%v，传入域名：%v", certBundle.DNSNames, userDomains)
		}
	}

	// 创建请求客户端
	openapiClient, err := openapi.NewClient(cdn.Endpoint, ctAccessKey, ctSecretKey)
	if err != nil {
		return nil, fmt.Errorf("创建请求客户端错误: %w", err)
	}
	// openapiClient.WithDebug()

	// 1. 域名绑定证书
	for _, domain := range userDomains {
		_, err = cdn.Action(openapiClient, domain, certBundle)
		if err != nil {
			return nil, err
		}
	}

	return &Response{
		Status:  "success",
		Message: "更新域名证书成功",
		Result: map[string]any{
			"domain":      ctDomain,
			"cert":        certBundle,
			"verifyChain": verifyChainText,
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

	ctAccessKey, ok := cfg["access_key"].(string)
	if !ok || ctAccessKey == "" {
		return nil, fmt.Errorf("access_key is required and must be a string")
	}
	ctSecretKey, ok := cfg["secret_key"].(string)
	if !ok || ctSecretKey == "" {
		return nil, fmt.Errorf("secret_key is required and must be a string")
	}
	ctDomain, ok := cfg["domain"].(string)
	if !ok || ctDomain == "" {
		return nil, fmt.Errorf("domain is required and must be a string")
	}

	// 解析证书字符串
	certBundle, err := core.ParseCertBundle([]byte(certStr), []byte(keyStr))
	if err != nil {
		return nil, fmt.Errorf("failed to parse cert bundle: %w", err)
	}

	// 验证证书链
	var verifyChainText string
	if err := certBundle.VerifyChain(); err != nil {
		verifyChainText = "❌ 证书链不完整或不被信任:" + err.Error()
	} else {
		verifyChainText = "✅ 证书链完整有效"
	}

	// 1. 检查证书是否过期
	if certBundle.IsExpired() {
		return nil, fmt.Errorf("证书已过期 %s", certBundle.NotAfter.Format(time.DateTime))
	}

	// 2. 解析传入域名
	userDomains, isMultiple := core.ParseDomainsFixedSeparator(ctDomain, ",")
	if isMultiple {
		if !certBundle.CanDomainsUseCert(userDomains) {
			return nil, fmt.Errorf("域名和证书不匹配，证书支持域名：%v，传入域名：%v", certBundle.DNSNames, userDomains)
		}
	}

	// 创建请求客户端
	openapiClient, err := openapi.NewClient(icdn.Endpoint, ctAccessKey, ctSecretKey)
	if err != nil {
		return nil, fmt.Errorf("创建请求客户端错误: %w", err)
	}
	// openapiClient.WithDebug()

	// 1. 域名绑定证书
	for _, domain := range userDomains {
		_, err := icdn.Action(openapiClient, domain, certBundle)
		if err != nil {
			return nil, err
		}
	}

	return &Response{
		Status:  "success",
		Message: "更新域名证书成功",
		Result: map[string]any{
			"domain":      ctDomain,
			"cert":        certBundle,
			"verifyChain": verifyChainText,
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

	ctAccessKey, ok := cfg["access_key"].(string)
	if !ok || ctAccessKey == "" {
		return nil, fmt.Errorf("access_key is required and must be a string")
	}
	ctSecretKey, ok := cfg["secret_key"].(string)
	if !ok || ctSecretKey == "" {
		return nil, fmt.Errorf("secret_key is required and must be a string")
	}
	ctDomain, ok := cfg["domain"].(string)
	if !ok || ctDomain == "" {
		return nil, fmt.Errorf("domain is required and must be a string")
	}

	// 解析证书字符串
	certBundle, err := core.ParseCertBundle([]byte(certStr), []byte(keyStr))
	if err != nil {
		return nil, fmt.Errorf("failed to parse cert bundle: %w", err)
	}

	// 验证证书链
	var verifyChainText string
	if err := certBundle.VerifyChain(); err != nil {
		verifyChainText = "❌ 证书链不完整或不被信任:" + err.Error()
	} else {
		verifyChainText = "✅ 证书链完整有效"
	}

	// 1. 检查证书是否过期
	if certBundle.IsExpired() {
		return nil, fmt.Errorf("证书已过期 %s", certBundle.NotAfter.Format(time.DateTime))
	}

	// 2. 解析传入域名
	userDomains, isMultiple := core.ParseDomainsFixedSeparator(ctDomain, ",")
	if isMultiple {
		if !certBundle.CanDomainsUseCert(userDomains) {
			return nil, fmt.Errorf("域名和证书不匹配，证书支持域名：%v，传入域名：%v", certBundle.DNSNames, userDomains)
		}
	}

	// 创建请求客户端
	openapiClient, err := openapi.NewClient(accessone.Endpoint, ctAccessKey, ctSecretKey)
	if err != nil {
		return nil, fmt.Errorf("创建请求客户端错误: %w", err)
	}
	// openapiClient.WithDebug()

	// 1. 域名绑定证书
	for _, domain := range userDomains {
		_, err := accessone.Action(openapiClient, domain, certBundle)
		if err != nil {
			return nil, err
		}
	}

	return &Response{
		Status:  "success",
		Message: "更新域名证书成功",
		Result: map[string]any{
			"domain":      ctDomain,
			"cert":        certBundle,
			"verifyChain": verifyChainText,
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

	ctAccessKey, ok := cfg["access_key"].(string)
	if !ok || ctAccessKey == "" {
		return nil, fmt.Errorf("access_key is required and must be a string")
	}
	ctSecretKey, ok := cfg["secret_key"].(string)
	if !ok || ctSecretKey == "" {
		return nil, fmt.Errorf("secret_key is required and must be a string")
	}

	// 解析证书字符串
	certBundle, err := core.ParseCertBundle([]byte(certStr), []byte(keyStr))
	if err != nil {
		return nil, fmt.Errorf("failed to parse cert bundle: %w", err)
	}

	// 验证证书链
	var verifyChainText string
	if err := certBundle.VerifyChain(); err != nil {
		verifyChainText = "❌ 证书链不完整或不被信任:" + err.Error()
	} else {
		verifyChainText = "✅ 证书链完整有效"
	}

	// 1. 检查证书是否过期
	if certBundle.IsExpired() {
		return nil, fmt.Errorf("证书已过期 %s", certBundle.NotAfter.Format(time.DateTime))
	}

	// 创建请求客户端
	openapiClient, err := openapi.NewClient(ccms.Endpoint, ctAccessKey, ctSecretKey)
	if err != nil {
		return nil, fmt.Errorf("创建请求客户端错误: %w", err)
	}
	// openapiClient.WithDebug()

	// 1. 上传证书
	isExist, err := ccms.Action(openapiClient, certBundle)
	if err != nil {
		return nil, err
	}
	if isExist {
		return &Response{
			Status:  "success",
			Message: "证书已存在",
			Result: map[string]any{
				"cert":        certBundle,
				"verifyChain": verifyChainText,
			},
		}, nil
	}

	return &Response{
		Status:  "success",
		Message: "上传证书成功",
		Result: map[string]any{
			"cert":        certBundle,
			"verifyChain": verifyChainText,
		},
	}, nil
}

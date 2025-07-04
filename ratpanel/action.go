package main

import (
	"fmt"
	"time"

	"github.com/dtapps/allinssl_plugins/core"
	"github.com/dtapps/allinssl_plugins/ratpanel/certificate"
	"github.com/dtapps/allinssl_plugins/ratpanel/openapi"
	"github.com/dtapps/allinssl_plugins/ratpanel/site"
)

// 部署到网站
func deploySiteAction(cfg map[string]any) (*Response, error) {

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

	ratURL, ok := cfg["url"].(string)
	if !ok || ratURL == "" {
		return nil, fmt.Errorf("url is required and must be a string")
	}
	ratUserID, ok := cfg["user_id"].(int)
	if !ok || ratUserID == 0 {
		return nil, fmt.Errorf("user_id is required and must be a string")
	}
	ratUserToken, ok := cfg["user_token"].(string)
	if !ok || ratUserToken == "" {
		return nil, fmt.Errorf("user_token is required and must be a string")
	}
	ratDomain, ok := cfg["domain"].(string)
	if !ok || ratDomain == "" {
		return nil, fmt.Errorf("domain is required and must be a string")
	}

	// 解析证书字符串
	certBundle, err := core.ParseCertBundle([]byte(certStr), []byte(keyStr))
	if err != nil {
		return nil, fmt.Errorf("failed to parse cert bundle: %w", err)
	}

	// 1. 检查证书是否过期
	if certBundle.IsExpired() {
		return nil, fmt.Errorf("证书已过期 %s", certBundle.NotAfter.Format(time.DateTime))
	}

	// 2. 解析传入域名
	userDomains, isMultiple := core.ParseDomainsFixedSeparator(ratDomain, ",")
	if isMultiple {
		if !certBundle.CanDomainsUseCert(userDomains) {
			return nil, fmt.Errorf("域名和证书不匹配，证书支持域名：%v，传入域名：%v", certBundle.DNSNames, userDomains)
		}
	}

	// 创建请求客户端
	openapiClient, err := openapi.NewClient(ratURL, ratUserID, ratUserToken)
	if err != nil {
		return nil, fmt.Errorf("创建请求客户端错误: %w", err)
	}
	// openapiClient.WithDebug()
	openapiClient.WithSkipVerify()

	// 1. 上传证书
	cardID, _, err := certificate.Action(openapiClient, certBundle)
	if err != nil {
		return nil, err
	}

	// 2. 域名绑定证书
	for _, domain := range userDomains {
		_, err = site.Action(openapiClient, domain, cardID, certBundle)
		if err != nil {
			return nil, err
		}
	}

	return &Response{
		Status:  "success",
		Message: "更新域名证书成功",
		Result: map[string]any{
			"domain": ratDomain,
			"cert":   certBundle,
		},
	}, nil
}

// 上传到证书管理
func deployCertificatesAction(cfg map[string]any) (*Response, error) {

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

	ratURL, ok := cfg["url"].(string)
	if !ok || ratURL == "" {
		return nil, fmt.Errorf("url is required and must be a string")
	}
	ratUserID, ok := cfg["user_id"].(int)
	if !ok || ratUserID == 0 {
		return nil, fmt.Errorf("user_id is required and must be a string")
	}
	ratUserToken, ok := cfg["user_token"].(string)
	if !ok || ratUserToken == "" {
		return nil, fmt.Errorf("user_token is required and must be a string")
	}

	// 解析证书字符串
	certBundle, err := core.ParseCertBundle([]byte(certStr), []byte(keyStr))
	if err != nil {
		return nil, fmt.Errorf("failed to parse cert bundle: %w", err)
	}

	// 1. 检查证书是否过期
	if certBundle.IsExpired() {
		return nil, fmt.Errorf("证书已过期 %s", certBundle.NotAfter.Format(time.DateTime))
	}

	// 创建请求客户端
	openapiClient, err := openapi.NewClient(ratURL, ratUserID, ratUserToken)
	if err != nil {
		return nil, fmt.Errorf("创建请求客户端错误: %w", err)
	}
	// openapiClient.WithDebug()
	openapiClient.WithSkipVerify()

	// 1. 上传证书
	_, isExist, err := certificate.Action(openapiClient, certBundle)
	if err != nil {
		return nil, err
	}
	if isExist {
		return &Response{
			Status:  "success",
			Message: "证书已存在",
			Result: map[string]any{
				"cert": certBundle,
			},
		}, nil
	}

	return &Response{
		Status:  "success",
		Message: "上传证书成功",
		Result: map[string]any{
			"cert": certBundle,
		},
	}, nil
}

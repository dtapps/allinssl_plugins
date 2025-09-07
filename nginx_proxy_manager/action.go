package main

import (
	"fmt"
	"time"

	"github.com/dtapps/allinssl_plugins/core"
	"github.com/dtapps/allinssl_plugins/nginx_proxy_manager/certificate"
	"github.com/dtapps/allinssl_plugins/nginx_proxy_manager/openapi"
	"github.com/dtapps/allinssl_plugins/nginx_proxy_manager/proxy_host"
)

// 部署证书到代理网站
func deployProxyHostsAction(cfg map[string]any) (*core.Response, error) {

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

	npmURL, ok := cfg["url"].(string)
	if !ok || npmURL == "" {
		return nil, fmt.Errorf("url is required and must be a string")
	}
	npmEmail, ok := cfg["email"].(string)
	if !ok || npmEmail == "" {
		return nil, fmt.Errorf("email is required and must be a string")
	}
	npmPassword, ok := cfg["password"].(string)
	if !ok || npmPassword == "" {
		return nil, fmt.Errorf("password is required and must be a string")
	}
	npmDomain, ok := cfg["domain"].(string)
	if !ok || npmDomain == "" {
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
	userDomains, isMultiple := core.ParseDomainsFixedSeparator(npmDomain, ",")
	if isMultiple {
		if !certBundle.CanDomainsUseCert(userDomains) {
			return nil, fmt.Errorf("域名和证书不匹配，证书支持域名：%v，传入域名：%v", certBundle.DNSNames, userDomains)
		}
	}

	// 创建请求客户端
	openapiClient, err := openapi.NewClient(npmURL, npmEmail, npmPassword)
	if err != nil {
		return nil, fmt.Errorf("创建请求客户端错误: %w", err)
	}
	// openapiClient.WithDebug()

	// 1. 先登录获取令牌
	openapiClient, err = openapiClient.WithLogin()
	if err != nil {
		return nil, fmt.Errorf("登录错误: %w", err)
	}

	// 2.设置令牌
	openapiClient.WithToken()

	// 3. 上传证书
	certID, _, err := certificate.Action(openapiClient, certBundle)
	if err != nil {
		return nil, err
	}

	// 4. 域名绑定证书
	for _, domain := range userDomains {
		_, err = proxy_host.Action(openapiClient, domain, certID, certBundle)
		if err != nil {
			return nil, err
		}
	}

	return &core.Response{
		Status:  "success",
		Message: "更新域名证书成功",
		Result: map[string]any{
			"domain": npmDomain,
			"cert":   certBundle.ResultInfo(),
		},
	}, nil
}

// 上传证书到证书管理
func deployCertificatesAction(cfg map[string]any) (*core.Response, error) {

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

	npmURL, ok := cfg["url"].(string)
	if !ok || npmURL == "" {
		return nil, fmt.Errorf("url is required and must be a string")
	}
	npmEmail, ok := cfg["email"].(string)
	if !ok || npmEmail == "" {
		return nil, fmt.Errorf("email is required and must be a string")
	}
	npmPassword, ok := cfg["password"].(string)
	if !ok || npmPassword == "" {
		return nil, fmt.Errorf("password is required and must be a string")
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
	openapiClient, err := openapi.NewClient(npmURL, npmEmail, npmPassword)
	if err != nil {
		return nil, fmt.Errorf("创建请求客户端错误: %w", err)
	}
	// openapiClient.WithDebug()

	// 1. 先登录获取令牌
	openapiClient, err = openapiClient.WithLogin()
	if err != nil {
		return nil, fmt.Errorf("登录错误: %w", err)
	}

	// 2.设置令牌
	openapiClient.WithToken()

	// 3. 上传证书
	_, isExist, err := certificate.Action(openapiClient, certBundle)
	if err != nil {
		return nil, err
	}
	if isExist {
		return &core.Response{
			Status:  "success",
			Message: "证书已存在",
			Result: map[string]any{
				"cert": certBundle.ResultInfo(),
			},
		}, nil
	}

	return &core.Response{
		Status:  "success",
		Message: "上传证书成功",
		Result: map[string]any{
			"cert": certBundle.ResultInfo(),
		},
	}, nil
}

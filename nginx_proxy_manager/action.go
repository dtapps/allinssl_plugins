package main

import (
	"fmt"

	"github.com/dtapps/allinssl_plugins/nginx_proxy_manager/certificate"
	"github.com/dtapps/allinssl_plugins/nginx_proxy_manager/core"
	"github.com/dtapps/allinssl_plugins/nginx_proxy_manager/openapi"
	"github.com/dtapps/allinssl_plugins/nginx_proxy_manager/proxy_host"
)

// 部署证书到代理网站
func deployProxyHostsAction(cfg map[string]any) (*Response, error) {

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

	baseURL, ok := cfg["url"].(string)
	if !ok || baseURL == "" {
		return nil, fmt.Errorf("url is required and must be a string")
	}
	email, ok := cfg["email"].(string)
	if !ok || email == "" {
		return nil, fmt.Errorf("email is required and must be a string")
	}
	password, ok := cfg["password"].(string)
	if !ok || password == "" {
		return nil, fmt.Errorf("password is required and must be a string")
	}
	domain, ok := cfg["domain"].(string)
	if !ok || domain == "" {
		return nil, fmt.Errorf("domain is required and must be a string")
	}

	// 解析证书字符串
	certBundle, err := core.ParseCertBundle([]byte(certStr), []byte(keyStr))
	if err != nil {
		return nil, fmt.Errorf("failed to parse cert bundle: %w", err)
	}

	// 创建请求客户端
	openapiClient, err := openapi.NewClient(baseURL, email, password)
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
	certID, err := certificate.Action(openapiClient, certBundle)
	if err != nil {
		return nil, err
	}

	// 4. 域名绑定证书
	_, err = proxy_host.Action(openapiClient, domain, certID, certBundle)
	if err != nil {
		return nil, err
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

// 上传证书到证书管理
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

	baseURL, ok := cfg["url"].(string)
	if !ok || baseURL == "" {
		return nil, fmt.Errorf("url is required and must be a string")
	}
	email, ok := cfg["email"].(string)
	if !ok || email == "" {
		return nil, fmt.Errorf("email is required and must be a string")
	}
	password, ok := cfg["password"].(string)
	if !ok || password == "" {
		return nil, fmt.Errorf("password is required and must be a string")
	}

	// 解析证书字符串
	certBundle, err := core.ParseCertBundle([]byte(certStr), []byte(keyStr))
	if err != nil {
		return nil, fmt.Errorf("failed to parse cert bundle: %w", err)
	}

	// 创建请求客户端
	openapiClient, err := openapi.NewClient(baseURL, email, password)
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
	_, err = certificate.Action(openapiClient, certBundle)
	if err != nil {
		return nil, err
	}

	return &Response{
		Status:  "success",
		Message: "上传证书成功",
		Result: map[string]any{
			"cert": certBundle,
		},
	}, nil
}

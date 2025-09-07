package main

import (
	"fmt"
	"time"

	"github.com/dtapps/allinssl_plugins/core"
	"github.com/dtapps/allinssl_plugins/safeline/openapi"
	"github.com/dtapps/allinssl_plugins/safeline/site"
)

// 部署到雷池WAF网站
func deploySiteAction(cfg map[string]any) (*core.Response, error) {

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

	wafURL, ok := cfg["url"].(string)
	if !ok || wafURL == "" {
		return nil, fmt.Errorf("url is required and must be a string")
	}
	wafApiKey, ok := cfg["api_token"].(string)
	if !ok || wafApiKey == "" {
		return nil, fmt.Errorf("api_token is required and must be a string")
	}
	wafDomain, ok := cfg["domain"].(string)
	if !ok || wafDomain == "" {
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
	userDomains, isMultiple := core.ParseDomainsFixedSeparator(wafDomain, ",")
	if isMultiple {
		if !certBundle.CanDomainsUseCert(userDomains) {
			return nil, fmt.Errorf("域名和证书不匹配，证书支持域名：%v，传入域名：%v", certBundle.DNSNames, userDomains)
		}
	}

	// 创建请求客户端
	openapiClient, err := openapi.NewClient(wafURL, wafApiKey)
	if err != nil {
		return nil, fmt.Errorf("创建请求客户端错误: %w", err)
	}
	// openapiClient.WithDebug()
	openapiClient.WithSkipVerify()

	// 1. 域名绑定证书
	for _, domain := range userDomains {
		_, err = site.Action(openapiClient, domain, certBundle)
		if err != nil {
			return nil, err
		}
	}

	return &core.Response{
		Status:  "success",
		Message: "更新域名证书成功",
		Result: map[string]any{
			"domain": wafDomain,
			"cert":   certBundle.ResultInfo(),
		},
	}, nil
}

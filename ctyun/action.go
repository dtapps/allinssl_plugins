package main

import (
	"context"
	"fmt"
	"time"

	"github.com/dtapps/allinssl_plugins/core"
	"github.com/dtapps/allinssl_plugins/ctyun/accessone"
	"github.com/dtapps/allinssl_plugins/ctyun/ccms"
	"github.com/dtapps/allinssl_plugins/ctyun/cdn"
	"github.com/dtapps/allinssl_plugins/ctyun/icdn"
	"github.com/dtapps/allinssl_plugins/ctyun/openapi"
)

type CommonConfig struct {
	Debug     bool   `json:"debug"`
	Cert      string `json:"cert"`
	Key       string `json:"key"`
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
}

// 解析公共配置
func parseCommonConfig(cfg map[string]any) (commonConfig *CommonConfig, err error) {
	// 检查是否开启调试模式
	isDebug, ok := cfg["debug"].(bool)
	if !ok {
		isDebug = false
	}

	// 证书字符串
	certStr, ok := cfg["cert"].(string)
	if !ok || certStr == "" {
		return nil, fmt.Errorf("cert is required and must be a string")
	}

	// 证书私钥字符串
	keyStr, ok := cfg["key"].(string)
	if !ok || keyStr == "" {
		return nil, fmt.Errorf("key is required and must be a string")
	}

	// access_key
	access_key, ok := cfg["access_key"].(string)
	if !ok || access_key == "" {
		return nil, fmt.Errorf("access_key is required and must be a string")
	}

	// secret_key
	secret_key, ok := cfg["secret_key"].(string)
	if !ok || secret_key == "" {
		return nil, fmt.Errorf("secret_key is required and must be a string")
	}

	// 返回公共配置
	return &CommonConfig{
		Debug:     isDebug,
		Cert:      certStr,
		Key:       keyStr,
		AccessKey: access_key,
		SecretKey: secret_key,
	}, nil
}

// 部署到天翼云CDN加速
func deployCdnAction(cfg map[string]any) (*core.Response, error) {

	ctx := context.Background()

	if cfg == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}
	commonConfig, err := parseCommonConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("parse common config error: %w", err)
	}

	// 域名
	domain, ok := cfg["domain"].(string)
	if !ok || domain == "" {
		return nil, fmt.Errorf("domain is required and must be a string")
	}

	// 解析证书字符串
	certBundle, err := core.ParseCertBundle([]byte(commonConfig.Cert), []byte(commonConfig.Key))
	if err != nil {
		return nil, fmt.Errorf("failed to parse cert bundle: %w", err)
	}

	// 检查证书是否过期
	if certBundle.IsExpired() {
		return nil, fmt.Errorf("证书已过期 %s", certBundle.NotAfter.Format(time.DateTime))
	}

	// 解析传入域名
	userDomains, isMultiple := core.ParseDomainsFixedSeparator(domain, ",")
	if isMultiple {
		if !certBundle.CanDomainsUseCert(userDomains) {
			return nil, fmt.Errorf("域名和证书不匹配，证书支持域名：%v，传入域名：%v", certBundle.DNSNames, userDomains)
		}
	}

	// 创建请求客户端
	openapiClient, err := openapi.NewClient(cdn.Endpoint, commonConfig.AccessKey, commonConfig.SecretKey)
	if err != nil {
		return nil, fmt.Errorf("创建请求客户端错误: %w", err)
	}
	if commonConfig.Debug {
		openapiClient.WithDebug()
	}

	// 域名绑定证书
	for _, domain := range userDomains {
		if domain == "" {
			continue
		}
		_, err = cdn.Action(ctx, openapiClient, domain, certBundle)
		if err != nil {
			return nil, err
		}
	}

	return &core.Response{
		Status:  "success",
		Message: "更新域名证书成功",
		Result: map[string]any{
			"domain":           domain,
			"cert_note":        certBundle.GetNoteShort(),
			"cert_expiry_time": certBundle.NotAfter.Format(time.DateTime),
		},
	}, nil
}

// 部署到天翼云全站加速
func deployIcdnAction(cfg map[string]any) (*core.Response, error) {

	ctx := context.Background()

	if cfg == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}
	commonConfig, err := parseCommonConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("parse common config error: %w", err)
	}

	// 域名
	domain, ok := cfg["domain"].(string)
	if !ok || domain == "" {
		return nil, fmt.Errorf("domain is required and must be a string")
	}

	// 解析证书字符串
	certBundle, err := core.ParseCertBundle([]byte(commonConfig.Cert), []byte(commonConfig.Key))
	if err != nil {
		return nil, fmt.Errorf("failed to parse cert bundle: %w", err)
	}

	// 检查证书是否过期
	if certBundle.IsExpired() {
		return nil, fmt.Errorf("证书已过期 %s", certBundle.NotAfter.Format(time.DateTime))
	}

	// 解析传入域名
	userDomains, isMultiple := core.ParseDomainsFixedSeparator(domain, ",")
	if isMultiple {
		if !certBundle.CanDomainsUseCert(userDomains) {
			return nil, fmt.Errorf("域名和证书不匹配，证书支持域名：%v，传入域名：%v", certBundle.DNSNames, userDomains)
		}
	}

	// 创建请求客户端
	openapiClient, err := openapi.NewClient(icdn.Endpoint, commonConfig.AccessKey, commonConfig.SecretKey)
	if err != nil {
		return nil, fmt.Errorf("创建请求客户端错误: %w", err)
	}
	if commonConfig.Debug {
		openapiClient.WithDebug()
	}

	// 域名绑定证书
	for _, domain := range userDomains {
		if domain == "" {
			continue
		}
		_, err := icdn.Action(ctx, openapiClient, domain, certBundle)
		if err != nil {
			return nil, err
		}
	}

	return &core.Response{
		Status:  "success",
		Message: "更新域名证书成功",
		Result: map[string]any{
			"domain":           domain,
			"cert_note":        certBundle.GetNoteShort(),
			"cert_expiry_time": certBundle.NotAfter.Format(time.DateTime),
		},
	}, nil
}

// 部署到天翼云边缘安全加速平台
func deployAccessoneAction(cfg map[string]any) (*core.Response, error) {

	ctx := context.Background()

	if cfg == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}
	commonConfig, err := parseCommonConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("parse common config error: %w", err)
	}

	// 域名
	domain, ok := cfg["domain"].(string)
	if !ok || domain == "" {
		return nil, fmt.Errorf("domain is required and must be a string")
	}

	// 解析证书字符串
	certBundle, err := core.ParseCertBundle([]byte(commonConfig.Cert), []byte(commonConfig.Key))
	if err != nil {
		return nil, fmt.Errorf("failed to parse cert bundle: %w", err)
	}

	// 检查证书是否过期
	if certBundle.IsExpired() {
		return nil, fmt.Errorf("证书已过期 %s", certBundle.NotAfter.Format(time.DateTime))
	}

	// 解析传入域名
	userDomains, isMultiple := core.ParseDomainsFixedSeparator(domain, ",")
	if isMultiple {
		if !certBundle.CanDomainsUseCert(userDomains) {
			return nil, fmt.Errorf("域名和证书不匹配，证书支持域名：%v，传入域名：%v", certBundle.DNSNames, userDomains)
		}
	}

	// 创建请求客户端
	openapiClient, err := openapi.NewClient(accessone.Endpoint, commonConfig.AccessKey, commonConfig.SecretKey)
	if err != nil {
		return nil, fmt.Errorf("创建请求客户端错误: %w", err)
	}
	if commonConfig.Debug {
		openapiClient.WithDebug()
	}

	// 1. 域名绑定证书
	for _, domain := range userDomains {
		if domain == "" {
			continue
		}
		_, err := accessone.Action(ctx, openapiClient, domain, certBundle)
		if err != nil {
			return nil, err
		}
	}

	return &core.Response{
		Status:  "success",
		Message: "更新域名证书成功",
		Result: map[string]any{
			"domain":           domain,
			"cert_note":        certBundle.GetNoteShort(),
			"cert_expiry_time": certBundle.NotAfter.Format(time.DateTime),
		},
	}, nil
}

// 上传证书到天翼云证书管理
func deployCcmsAction(cfg map[string]any) (*core.Response, error) {

	ctx := context.Background()

	if cfg == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}
	commonConfig, err := parseCommonConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("parse common config error: %w", err)
	}

	// 解析证书字符串
	certBundle, err := core.ParseCertBundle([]byte(commonConfig.Cert), []byte(commonConfig.Key))
	if err != nil {
		return nil, fmt.Errorf("failed to parse cert bundle: %w", err)
	}

	// 检查证书是否过期
	if certBundle.IsExpired() {
		return nil, fmt.Errorf("证书已过期 %s", certBundle.NotAfter.Format(time.DateTime))
	}

	// 创建请求客户端
	openapiClient, err := openapi.NewClient(ccms.Endpoint, commonConfig.AccessKey, commonConfig.SecretKey)
	if err != nil {
		return nil, fmt.Errorf("创建请求客户端错误: %w", err)
	}
	if commonConfig.Debug {
		openapiClient.WithDebug()
	}

	// 上传证书
	isExist, err := ccms.Action(ctx, openapiClient, certBundle)
	if err != nil {
		return nil, err
	}
	if isExist {
		return &core.Response{
			Status:  "success",
			Message: "证书已存在",
			Result: map[string]any{
				"cert_note":        certBundle.GetNoteShort(),
				"cert_expiry_time": certBundle.NotAfter.Format(time.DateTime),
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

package main

import (
	"context"
	"fmt"
	"time"

	"github.com/dtapps/allinssl_plugins/core"
	"github.com/dtapps/allinssl_plugins/edgeone/openapi"
	"github.com/dtapps/allinssl_plugins/edgeone/ssl"
	"github.com/dtapps/allinssl_plugins/edgeone/teo"
)

type CommonConfig struct {
	Debug     bool   `json:"debug"`
	Cert      string `json:"cert"`
	Key       string `json:"key"`
	SecretID  string `json:"secret_id"`
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

	// 腾讯云SecretId
	secret_id, ok := cfg["secret_id"].(string)
	if !ok || secret_id == "" {
		return nil, fmt.Errorf("secret_id is required and must be a string")
	}

	// 腾讯云SecretKey
	secret_key, ok := cfg["secret_key"].(string)
	if !ok || secret_key == "" {
		return nil, fmt.Errorf("secret_key is required and must be a string")
	}

	// 返回公共配置
	return &CommonConfig{
		Debug:     isDebug,
		Cert:      certStr,
		Key:       keyStr,
		SecretID:  secret_id,
		SecretKey: secret_key,
	}, nil
}

// 部署到边缘安全加速平台EO
func deployTeoAction(cfg map[string]any) (*core.Response, error) {

	ctx := context.Background()

	if cfg == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}
	commonConfig, err := parseCommonConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("parse common config error: %w", err)
	}

	// 加速域名所属站点 ID
	zoneID, ok := cfg["zone_id"].(string)
	if !ok || zoneID == "" {
		return nil, fmt.Errorf("zone_id is required and must be a string")
	}

	// 加速域名
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
	openapiClient, err := openapi.NewClient(commonConfig.SecretID, commonConfig.SecretKey)
	if err != nil {
		return nil, fmt.Errorf("创建请求客户端 错误: %w", err)
	}
	if commonConfig.Debug {
		openapiClient.WithDebug()
	}

	// 域名绑定证书
	for _, domain := range userDomains {
		_, err := teo.Action(ctx, openapiClient, certBundle, &teo.Params{
			Debug:  commonConfig.Debug,
			ZoneID: zoneID,
			Domain: domain,
		})
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

// 上传到SSL证书
func deploySslAction(cfg map[string]any) (*core.Response, error) {

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
	openapiClient, err := openapi.NewClient(commonConfig.SecretID, commonConfig.SecretKey)
	if err != nil {
		return nil, fmt.Errorf("创建请求客户端 错误: %w", err)
	}
	if commonConfig.Debug {
		openapiClient.WithDebug()
	}

	// 上传证书
	sslResp, err := ssl.Action(ctx, openapiClient, certBundle, &ssl.Params{
		Debug: commonConfig.Debug,
	})
	if err != nil {
		return nil, err
	}
	if sslResp.IsExist {
		return &core.Response{
			Status:  "success",
			Message: "证书已存在",
			Result: map[string]any{
				"cert_note":        certBundle.GetNoteShort(),
				"cert_expiry_time": certBundle.NotAfter.Format(time.DateTime),
				"cert_id":          sslResp.CertID,
			},
		}, nil
	}

	return &core.Response{
		Status:  "success",
		Message: "上传证书成功",
		Result: map[string]any{
			"cert_note":        certBundle.GetNoteShort(),
			"cert_expiry_time": certBundle.NotAfter.Format(time.DateTime),
			"cert_id":          sslResp.CertID,
		},
	}, nil
}

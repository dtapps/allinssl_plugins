package main

import (
	"fmt"
	"time"

	"github.com/dtapps/allinssl_plugins/core"
	"github.com/dtapps/allinssl_plugins/lucky/certificate"
	"github.com/dtapps/allinssl_plugins/lucky/openapi"
)

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

	luckyURL, ok := cfg["url"].(string)
	if !ok || luckyURL == "" {
		return nil, fmt.Errorf("url is required and must be a string")
	}
	luckyOpenToken, ok := cfg["open_token"].(string)
	if !ok || luckyOpenToken == "" {
		return nil, fmt.Errorf("open_token is required and must be a string")
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
	openapiClient, err := openapi.NewClient(luckyURL, luckyOpenToken)
	if err != nil {
		return nil, fmt.Errorf("创建请求客户端错误: %w", err)
	}
	// openapiClient.WithDebug()
	openapiClient.WithSkipVerify()

	// 3. 上传证书
	isExist, err := certificate.Action(openapiClient, certBundle)
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

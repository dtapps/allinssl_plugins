package main

import (
	"fmt"

	"github.com/dtapps/allinssl_plugins/synology/certificate"
	"github.com/dtapps/allinssl_plugins/synology/core"
	"github.com/dtapps/allinssl_plugins/synology/openapi"
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

	dsmURL, ok := cfg["url"].(string)
	if !ok || dsmURL == "" {
		return nil, fmt.Errorf("url is required and must be a string")
	}
	dsmUsername, ok := cfg["username"].(string)
	if !ok || dsmUsername == "" {
		return nil, fmt.Errorf("username is required and must be a string")
	}
	dsmPassword, ok := cfg["password"].(string)
	if !ok || dsmPassword == "" {
		return nil, fmt.Errorf("password is required and must be a string")
	}

	// 解析证书字符串
	certBundle, err := core.ParseCertBundle([]byte(certStr), []byte(keyStr))
	if err != nil {
		return nil, fmt.Errorf("failed to parse cert bundle: %w", err)
	}

	// 创建请求客户端
	openapiClient, err := openapi.NewClient(dsmURL, dsmUsername, dsmPassword)
	if err != nil {
		return nil, fmt.Errorf("创建请求客户端错误: %w", err)
	}
	// openapiClient.WithDebug()
	openapiClient.WithSkipVerify()

	// 1. 先登录获取令牌
	openapiClient, err = openapiClient.WithLogin()
	if err != nil {
		return nil, fmt.Errorf("登录错误: %w", err)
	}

	// 2.设置令牌
	openapiClient.WithToken()

	// 上传证书
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

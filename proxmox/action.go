package main

import (
	"fmt"

	"github.com/dtapps/allinssl_plugins/proxmox/certificate"
	"github.com/dtapps/allinssl_plugins/proxmox/core"
	"github.com/dtapps/allinssl_plugins/proxmox/openapi"
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

	pveURL, ok := cfg["url"].(string)
	if !ok || pveURL == "" {
		return nil, fmt.Errorf("url is required and must be a string")
	}
	pveNode, ok := cfg["node"].(string)
	if !ok || pveNode == "" {
		return nil, fmt.Errorf("node is required and must be a string")
	}
	pveUser, ok := cfg["user"].(string)
	if !ok || pveUser == "" {
		return nil, fmt.Errorf("user is required and must be a string")
	}
	pveTokenID, ok := cfg["token_id"].(string)
	if !ok || pveTokenID == "" {
		return nil, fmt.Errorf("token_id is required and must be a string")
	}
	pveTokenSecret, ok := cfg["token_secret"].(string)
	if !ok || pveTokenSecret == "" {
		return nil, fmt.Errorf("token_secret is required and must be a string")
	}

	// 解析证书字符串
	certBundle, err := core.ParseCertBundle([]byte(certStr), []byte(keyStr))
	if err != nil {
		return nil, fmt.Errorf("failed to parse cert bundle: %w", err)
	}

	// 创建请求客户端
	openapiClient, err := openapi.NewClient(pveURL, pveUser, pveTokenID, pveTokenSecret)
	if err != nil {
		return nil, fmt.Errorf("创建请求客户端错误: %w", err)
	}
	// openapiClient.WithDebug()
	openapiClient.WithSkipVerify()

	// 上传证书
	isExist, err := certificate.Action(openapiClient, pveNode, certBundle)
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

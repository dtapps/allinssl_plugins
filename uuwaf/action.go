package main

import (
	"fmt"
	"time"

	"github.com/dtapps/allinssl_plugins/core"
	"github.com/dtapps/allinssl_plugins/uuwaf/certificate"
	"github.com/dtapps/allinssl_plugins/uuwaf/openapi"
)

// 上传证书到证书管理
func deployCertificatesAction(cfg map[string]any, apiVersion string) (*core.Response, error) {

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

	omURL, ok := cfg["url"].(string)
	if !ok || omURL == "" {
		return nil, fmt.Errorf("url is required and must be a string")
	}
	omUsername, ok := cfg["username"].(string)
	if !ok || omUsername == "" {
		return nil, fmt.Errorf("username is required and must be a string")
	}
	omPassword, ok := cfg["password"].(string)
	if !ok || omPassword == "" {
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
	openapiClient, err := openapi.NewClient(omURL, omUsername, omPassword)
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

	if apiVersion == "v6" {

		// 2.设置令牌
		openapiClient.WithV6Token()

		// 3. 上传证书
		isExist, err := certificate.V6Action(openapiClient, certBundle)
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

	} else if apiVersion == "v7" {

		// 2.设置令牌
		openapiClient.WithV7Token()

		// 3. 上传证书
		isExist, err := certificate.V7Action(openapiClient, certBundle)
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

	} else {
		return nil, fmt.Errorf("不支持的 API 版本: %s", apiVersion)
	}

	return &core.Response{
		Status:  "success",
		Message: "上传证书成功",
		Result: map[string]any{
			"cert": certBundle.ResultInfo(),
		},
	}, nil
}

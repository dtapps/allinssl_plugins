package main

import (
	"github.com/dtapps/allinssl_plugins/core"
)

var pluginMeta = map[string]any{
	"name":        "aawaf",
	"description": "部署到堡塔云WAF",
	"version":     "1.0.1",
	"author":      "dtapps",
	"config": map[string]any{
		"url":     "堡塔云WAF 主机IP或域名，包含协议和端口，https://example.com 或 https://0.0.0.0:8379",
		"api_key": "堡塔云WAF 接口密钥",
	},
	"actions": []core.ActionInfo{
		{
			Name:        "site",
			Description: "部署到网站",
			Params: map[string]any{
				"domain": "域名，多个域名使用逗号分隔（需要是泛域名证书） 例如：example.com,www.example.com",
			},
		},
	},
}

func main() {
	req, err := core.ReadRequest()
	if err != nil {
		core.OutputError("请求处理失败", err)
		return
	}

	// 处理标准动作
	if core.HandleStandardActions(req, pluginMeta) {
		return
	}

	// 处理插件特有动作
	switch req.Action {
	case "site":
		rep, err := deploySiteAction(req.Params)
		if err != nil {
			core.OutputError("部署到网站 失败", err)
			return
		}
		core.OutputJSON(rep)
	default:
		core.OutputJSON(&core.Response{
			Status:  "error",
			Message: "未知 action: " + req.Action,
		})
	}
}

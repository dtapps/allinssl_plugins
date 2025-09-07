package main

import (
	"github.com/dtapps/allinssl_plugins/core"
)

var pluginMeta = map[string]any{
	"name":        "synology",
	"description": "部署到Synology",
	"version":     "1.0.3",
	"author":      "dtapps",
	"config": map[string]any{
		"url":      "Synology 主机IP或域名，包含协议和端口，例如：https://example.com 或 https://0.0.0.0:5001",
		"username": "Synology 用户名",
		"password": "Synology 密码",
	},
	"actions": []core.ActionInfo{
		{
			Name:        "certificates",
			Description: "上传到证书",
			Params: map[string]any{
				"as_default": "是否将证书设置为默认证书",
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
	case "certificates":
		rep, err := deployCertificatesAction(req.Params)
		if err != nil {
			core.OutputError("上传到证书 失败", err)
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

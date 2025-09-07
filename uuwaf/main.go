package main

import (
	"github.com/dtapps/allinssl_plugins/core"
)

var pluginMeta = map[string]any{
	"name":        "uuwaf",
	"description": "部署到南墙WEB应用防火墙",
	"version":     "1.0.0",
	"author":      "dtapps",
	"config": map[string]any{
		"url":      "南墙WEB应用防火墙 主机IP或域名，包含协议和端口，https://example.com 或 https://0.0.0.0:4443",
		"username": "南墙WEB应用防火墙 登录用户名",
		"password": "南墙WEB应用防火墙 登录密码",
	},
	"actions": []core.ActionInfo{
		{
			Name:        "certificates_v6",
			Description: "上传到网站证书（V6版本）",
			Params:      map[string]any{},
		},
		{
			Name:        "certificates_v7",
			Description: "上传到网站证书（V7版本）",
			Params:      map[string]any{},
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
	case "certificates_v6":
		rep, err := deployCertificatesAction(req.Params, "v6")
		if err != nil {
			core.OutputError("上传到网站证书（V6版本） 失败", err)
			return
		}
		core.OutputJSON(rep)
	case "certificates_v7":
		rep, err := deployCertificatesAction(req.Params, "v7")
		if err != nil {
			core.OutputError("上传到网站证书（V7版本） 失败", err)
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

package main

import (
	"github.com/dtapps/allinssl_plugins/core"
)

var pluginMeta = map[string]any{
	"name":        "lucky",
	"description": "部署到Lucky",
	"version":     "1.0.0",
	"author":      "dtapps",
	"config": map[string]any{
		"url":        "主机IP或域名，包含协议和端口加入口，http://example.com/xxx 或 http://0.0.0.0:16601/xxx",
		"open_token": "设置中心的OpenToken",
	},
	"actions": []core.ActionInfo{
		{
			Name:        "certificates",
			Description: "上传到SSL/TLS证书",
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
	case "certificates":
		rep, err := deployCertificatesAction(req.Params)
		if err != nil {
			core.OutputError("上传到SSL/TLS证书 失败", err)
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

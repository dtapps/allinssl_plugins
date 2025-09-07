package main

import (
	"github.com/dtapps/allinssl_plugins/core"
)

var pluginMeta = map[string]any{
	"name":        "open_resty_manager",
	"description": "部署到OpenResty",
	"version":     "1.0.0",
	"author":      "dtapps",
	"config": map[string]any{
		"url":      "OpenResty 主机IP或域名，包含协议和端口，https://example.com 或 https://0.0.0.0:34567",
		"username": "OpenResty 登录用户名",
		"password": "OpenResty 登录密码",
	},
	"actions": []core.ActionInfo{
		{
			Name:        "site",
			Description: "部署到网站管理",
			Params: map[string]any{
				"domain": "域名，多个域名使用逗号分隔（需要是泛域名证书） 例如：example.com,www.example.com",
			},
		},
		{
			Name:        "certificates",
			Description: "上传到网站证书",
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

	case "site":
		rep, err := deploySiteAction(req.Params)
		if err != nil {
			core.OutputError("部署到网站管理 失败", err)
			return
		}
		core.OutputJSON(rep)
	case "certificates":
		rep, err := deployCertificatesAction(req.Params)
		if err != nil {
			core.OutputError("上传到网站证书 失败", err)
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

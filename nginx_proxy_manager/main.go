package main

import (
	"github.com/dtapps/allinssl_plugins/core"
)

var pluginMeta = map[string]any{
	"name":        "nginx_proxy_manager",
	"description": "部署到Nginx Proxy Manager",
	"version":     "1.0.4",
	"author":      "dtapps",
	"config": map[string]any{
		"url":      "Nginx Proxy Manager 主机IP或域名，包含协议和端口，http://example.com 或 http://0.0.0.0:81",
		"email":    "Nginx Proxy Manager 登录邮箱",
		"password": "Nginx Proxy Manager 登录密码",
	},
	"actions": []core.ActionInfo{
		{
			Name:        "proxy_hosts",
			Description: "部署到Proxy Hosts",
			Params: map[string]any{
				"domain": "域名，多个域名使用逗号分隔（需要是泛域名证书） 例如：example.com,www.example.com",
			},
		},
		{
			Name:        "certificates",
			Description: "上传到SSL Certificates",
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
	case "proxy_hosts":
		rep, err := deployProxyHostsAction(req.Params)
		if err != nil {
			core.OutputError("部署到Proxy Hosts 失败", err)
			return
		}
		core.OutputJSON(rep)
	case "certificates":
		rep, err := deployCertificatesAction(req.Params)
		if err != nil {
			core.OutputError("上传到SSL Certificates 失败", err)
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

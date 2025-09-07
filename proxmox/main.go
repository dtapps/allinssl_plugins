package main

import (
	"github.com/dtapps/allinssl_plugins/core"
)

var pluginMeta = map[string]any{
	"name":        "proxmox",
	"description": "部署到Proxmox VE",
	"version":     "1.0.2",
	"author":      "dtapps",
	"config": map[string]any{
		"url":          "Proxmox VE 主机IP或域名，包含协议和端口，例如：https://example.com 或 https://0.0.0.0:8006",
		"node":         "Proxmox VE 节点名称，例如：pve",
		"user":         "Proxmox VE 用户名和领域，例如：root@pam",
		"token_id":     "Proxmox VE 令牌 ID",
		"token_secret": "Proxmox VE 令牌 密钥",
	},
	"actions": []core.ActionInfo{
		{
			Name:        "certificates",
			Description: "上传到证书",
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

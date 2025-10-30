package main

import (
	"github.com/dtapps/allinssl_plugins/core"
)

var pluginMeta = map[string]any{
	"name":        "edgeone",
	"description": "部署到腾讯云/EdgeOne",
	"version":     "1.0.1",
	"author":      "dtapps",
	"config": map[string]any{
		"secret_id":  "腾讯云 SecretId",
		"secret_key": "腾讯云 SecretKey",
	},
	"actions": []core.ActionInfo{
		{
			Name:        "teo",
			Description: "部署到边缘安全加速平台EO",
			Params: map[string]any{
				"zone_id": "站点 ID",
				"domain":  "域名，多个域名使用逗号分隔（需要是泛域名证书） 例如：example.com,www.example.com",
			},
		},
		{
			Name:        "ssl",
			Description: "上传到SSL证书",
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
	case "teo":
		rep, err := deployTeoAction(req.Params)
		if err != nil {
			core.OutputError("部署到部署到边缘安全加速平台EO 失败", err)
			return
		}
		core.OutputJSON(rep)
	case "ssl":
		rep, err := deploySslAction(req.Params)
		if err != nil {
			core.OutputError("上传到SSL证书 失败", err)
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

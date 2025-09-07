package main

import (
	"github.com/dtapps/allinssl_plugins/core"
)

var pluginMeta = map[string]any{
	"name":        "ctyun",
	"description": "部署到天翼云",
	"version":     "1.0.7",
	"author":      "dtapps",
	"config": map[string]any{
		"access_key": "天翼云 AccessKey",
		"secret_key": "天翼云 SecretKey",
	},
	"actions": []core.ActionInfo{
		{
			Name:        "cdn",
			Description: "部署到CDN加速",
			Params: map[string]any{
				"domain": "域名，多个域名使用逗号分隔（需要是泛域名证书） 例如：example.com,www.example.com",
			},
		},
		{
			Name:        "icdn",
			Description: "部署到全站加速",
			Params: map[string]any{
				"domain": "域名，多个域名使用逗号分隔（需要是泛域名证书） 例如：example.com,www.example.com",
			},
		},
		{
			Name:        "accessone",
			Description: "部署到边缘安全加速平台",
			Params: map[string]any{
				"domain": "域名，多个域名使用逗号分隔（需要是泛域名证书） 例如：example.com,www.example.com",
			},
		},
		{
			Name:        "ccms",
			Description: "上传到证书管理",
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
	case "cdn":
		rep, err := deployCdnAction(req.Params)
		if err != nil {
			core.OutputError("部署到CDN加速 失败", err)
			return
		}
		core.OutputJSON(rep)
	case "icdn":
		rep, err := deployIcdnAction(req.Params)
		if err != nil {
			core.OutputError("部署到全站加速 失败", err)
			return
		}
		core.OutputJSON(rep)
	case "accessone":
		rep, err := deployAccessoneAction(req.Params)
		if err != nil {
			core.OutputError("部署到边缘安全加速平台 失败", err)
			return
		}
		core.OutputJSON(rep)
	case "ccms":
		rep, err := deployCcmsAction(req.Params)
		if err != nil {
			core.OutputError("上传到证书管理 失败", err)
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

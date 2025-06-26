package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type ActionInfo struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Params      map[string]any `json:"params,omitempty"` // 可选参数
}

type Request struct {
	Action string         `json:"action"`
	Params map[string]any `json:"params"`
}

type Response struct {
	Status  string         `json:"status"`
	Message string         `json:"message"`
	Result  map[string]any `json:"result"`
}

var pluginMeta = map[string]any{
	"name":        "nginx_proxy_manager",
	"description": "部署到Nginx Proxy Manager",
	"version":     "1.0.1",
	"author":      "dtapps",
	"config": map[string]any{
		"url":      "Nginx Proxy Manager URL，for example：http://xxxx:81",
		"email":    "Nginx Proxy Manager Login Email address",
		"password": "Nginx Proxy Manager Login Password",
	},
	"actions": []ActionInfo{
		{
			Name:        "proxy_hosts",
			Description: "部署到 Nginx Proxy Manager Proxy Hosts",
			Params: map[string]any{
				"domain": "域名",
			},
		},
		{
			Name:        "certificates",
			Description: "上传到 Nginx Proxy Manager SSL Certificates",
		},
	},
}

func outputJSON(resp *Response) {
	_ = json.NewEncoder(os.Stdout).Encode(resp)
}

func outputError(msg string, err error) {
	outputJSON(&Response{
		Status:  "error",
		Message: fmt.Sprintf("%s: %v", msg, err),
	})
}

func main() {

	var req Request
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		outputError("读取输入失败", err)
		return
	}

	if err := json.Unmarshal(input, &req); err != nil {
		outputError("解析请求失败", err)
		return
	}

	switch req.Action {
	case "get_metadata":
		outputJSON(&Response{
			Status:  "success",
			Message: "插件信息",
			Result:  pluginMeta,
		})
	case "list_actions":
		outputJSON(&Response{
			Status:  "success",
			Message: "支持的动作",
			Result:  map[string]any{"actions": pluginMeta["actions"]},
		})
	case "proxy_hosts":
		rep, err := deployProxyHostsAction(req.Params)
		if err != nil {
			outputError("部署证书到代理网站失败", err)
			return
		}
		outputJSON(rep)
	case "certificates":
		rep, err := deployCertificatesAction(req.Params)
		if err != nil {
			outputError("上传证书到证书管理失败", err)
			return
		}
		outputJSON(rep)
	default:
		outputJSON(&Response{
			Status:  "error",
			Message: "未知 action: " + req.Action,
		})
	}
}

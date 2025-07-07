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
	"name":        "open_resty_manager",
	"description": "部署到OpenResty",
	"version":     "1.0.0",
	"author":      "dtapps",
	"config": map[string]any{
		"url":      "OpenResty 主机IP或域名，包含协议和端口，https://example.com 或 https://0.0.0.0:34567",
		"username": "OpenResty 登录用户名",
		"password": "OpenResty 登录密码",
	},
	"actions": []ActionInfo{
		{
			Name:        "site",
			Description: "部署到 网站管理",
			Params: map[string]any{
				"domain": "域名，多个域名使用逗号分隔（需要是泛域名证书） 例如：example.com,www.example.com",
			},
		},
		{
			Name:        "certificates",
			Description: "上传到 网站证书",
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
	case "site":
		rep, err := deploySiteAction(req.Params)
		if err != nil {
			outputError("部署证书到网站失败", err)
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

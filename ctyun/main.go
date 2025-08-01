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
	"name":        "ctyun",
	"description": "部署到天翼云",
	"version":     "1.0.3",
	"author":      "dtapps",
	"config": map[string]any{
		"access_key": "天翼云 AccessKey",
		"secret_key": "天翼云 SecretKey",
	},
	"actions": []ActionInfo{
		{
			Name:        "cdn",
			Description: "部署到天翼云CDN加速",
			Params: map[string]any{
				"domain": "域名，多个域名使用逗号分隔（需要是泛域名证书） 例如：example.com,www.example.com",
			},
		},
		{
			Name:        "icdn",
			Description: "部署到天翼云全站加速",
			Params: map[string]any{
				"domain": "域名，多个域名使用逗号分隔（需要是泛域名证书） 例如：example.com,www.example.com",
			},
		},
		{
			Name:        "accessone",
			Description: "部署到天翼云边缘安全加速平台",
			Params: map[string]any{
				"domain": "域名，多个域名使用逗号分隔（需要是泛域名证书） 例如：example.com,www.example.com",
			},
		},
		{
			Name:        "ccms",
			Description: "上传到天翼云证书管理",
			Params:      map[string]any{},
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
	case "cdn":
		rep, err := deployCdnAction(req.Params)
		if err != nil {
			outputError("CDN加速 部署失败", err)
			return
		}
		outputJSON(rep)
	case "icdn":
		rep, err := deployIcdnAction(req.Params)
		if err != nil {
			outputError("全站加速 部署失败", err)
			return
		}
		outputJSON(rep)
	case "accessone":
		rep, err := deployAccessoneAction(req.Params)
		if err != nil {
			outputError("边缘安全加速平台 部署失败", err)
			return
		}
		outputJSON(rep)
	case "ccms":
		rep, err := deployCcmsAction(req.Params)
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

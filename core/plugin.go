package core

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// ActionInfo 插件动作信息
type ActionInfo struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Params      map[string]any `json:"params,omitempty"` // 可选参数
}

// Request 插件请求结构
type Request struct {
	Action string         `json:"action"`
	Params map[string]any `json:"params"`
}

// Response 插件响应结构
type Response struct {
	Status  string         `json:"status"`
	Message string         `json:"message"`
	Result  map[string]any `json:"result"`
}

// OutputJSON 输出JSON响应
func OutputJSON(resp *Response) {
	_ = json.NewEncoder(os.Stdout).Encode(resp)
}

// OutputError 输出错误响应
func OutputError(msg string, err error) {
	OutputJSON(&Response{
		Status:  "error",
		Message: fmt.Sprintf("%s: %v", msg, err),
	})
}

// OutputSuccess 输出成功响应
func OutputSuccess(msg string, result map[string]any) {
	OutputJSON(&Response{
		Status:  "success",
		Message: msg,
		Result:  result,
	})
}

// ReadRequest 从标准输入读取并解析请求
func ReadRequest() (*Request, error) {
	var req Request
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		return nil, fmt.Errorf("读取输入失败: %w", err)
	}

	if err := json.Unmarshal(input, &req); err != nil {
		return nil, fmt.Errorf("解析请求失败: %w", err)
	}

	return &req, nil
}

// HandleStandardActions 处理标准动作（get_metadata, list_actions）
func HandleStandardActions(req *Request, pluginMeta map[string]any) bool {
	switch req.Action {
	case "get_metadata":
		OutputSuccess("插件信息", pluginMeta)
		return true
	case "list_actions":
		OutputSuccess("支持的动作", map[string]any{"actions": pluginMeta["actions"]})
		return true
	}
	return false
}

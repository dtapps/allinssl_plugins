package types

const StatusCodeSuccess = 100000

type BaseResponse struct {
	StatusCode   int    `json:"statusCode"`             // 状态码
	Message      string `json:"message"`                // 结果简述
	Error        string `json:"error,omitempty"`        // 错误码
	ErrorMessage string `json:"errorMessage,omitempty"` // 错误详情
}

type CommonResponse[T any] struct {
	BaseResponse   // 嵌入公共基础结构体
	ReturnObj    T `json:"returnObj,omitempty"` // 返回对象（具体类型由 T 决定）
}

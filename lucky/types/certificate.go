package types

// 获取证书列表 响应参数
type CertificateListResponse struct {
	Ret  int    `json:"ret"`           // 状态码
	Msg  string `json:"msg,omitempty"` // 描述
	List []struct {
		Key       string `json:"Key"`
		Remark    string `json:"Remark"`
		Enable    bool   `json:"Enable"`
		AddTime   string `json:"AddTime"`
		CertsInfo struct {
			Domains       []string `json:"Domains"`
			NotBeforeTime string   `json:"NotBeforeTime"`
			NotAfterTime  string   `json:"NotAfterTime"`
		} `json:"CertsInfo"`
		AddFrom              string         `json:"AddFrom"`
		UpdateTime           string         `json:"UpdateTime"`
		ExtParams            map[string]any `json:"ExtParams"`
		MappingToPath        bool           `json:"MappingToPath"`
		MappingPath          string         `json:"MappingPath"`
		MappingChangeScript  string         `json:"MappingChangeScript"`
		ACMEErrMsg           string         `json:"ACMEErrMsg"`
		ACMEing              bool           `json:"ACMEing"`
		ACMECancelingProcess bool           `json:"ACMECancelingProcess"`
		SyncInfo             any            `json:"SyncInfo"`
	} `json:"list,omitempty"`
}

// 上传证书 响应参数
type CertificateUpdateResponse struct {
	Ret    int    `json:"ret"`           // 状态码
	Msg    string `json:"msg,omitempty"` // 描述
	Base64 string `json:"base64"`        // 证书base64
	File   string `json:"file"`          // 证书文件
}

// 创建证书 响应参数
type CertificateCreateResponse struct {
	Ret int    `json:"ret"`           // 状态码
	Msg string `json:"msg,omitempty"` // 描述
}

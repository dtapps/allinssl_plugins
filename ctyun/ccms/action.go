package ccms

import (
	"encoding/json"
	"fmt"

	"github.com/dtapps/allinssl_plugins/ctyun/openapi"
	"github.com/dtapps/allinssl_plugins/ctyun/types"
)

type Client struct {
	openapi  *openapi.Client // 请求
	endpoint string          // 终端节点
}

// 证书管理服务
// https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=152&data=204&vid=283
func NewClient(accessKeyId string, secretAccessKey string) (c *Client, err error) {
	c = &Client{
		endpoint: "https://ccms-global.ctapi.ctyun.cn",
	}
	c.openapi, err = openapi.NewClient(c.endpoint, accessKeyId, secretAccessKey)
	return c, err
}

// 查询用户证书列表
// https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=152&api=17233&data=204&isNormal=1&vid=283
func (c *Client) GetQueryCertList() (response types.CommonResponse[types.CcmsCertificateListResponse], err error) {
	paramsStr, err := json.Marshal(types.CcmsQueryCertListRequest{
		PageSize: 1,
		PageNum:  100,
		Origin:   "UPLOAD",
	})
	if err != nil {
		return
	}
	resp, err := c.openapi.R().
		SetContentType("application/json").
		SetBody(string(paramsStr)).
		SetResult(&response).
		Post("/v1/certificate/list")
	if err != nil {
		return
	}
	if response.StatusCode != 200 {
		err = fmt.Errorf("%s", resp.String())
		return
	}
	return
}

// 上传证书
// name: 证书名称
// certificate: 证书字符串
// privateKey: 私钥字符串
// https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=152&api=17243&data=204&isNormal=1&vid=283
func (c *Client) PostUpdateCertInfo(name string, certificate string, privateKey string) (response types.CommonResponse[any], err error) {
	paramsStr, err := json.Marshal(types.CcmsUpdateCertInfoRequest{
		Name:               name,            // 证书名称
		Certificate:        certificate,     // 证书字符串
		PrivateKey:         privateKey,      // 私钥字符串
		EncryptionStandard: "INTERNATIONAL", // 加密标准
	})
	if err != nil {
		return
	}
	resp, err := c.openapi.R().
		SetContentType("application/json").
		SetBody(string(paramsStr)).
		SetResult(&response).
		Post("/v1/certificate/upload")
	if err != nil {
		return
	}
	if response.StatusCode != 200 {
		err = fmt.Errorf("%s", resp.String())
		return
	}
	return
}

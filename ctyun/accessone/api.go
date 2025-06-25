package accessone

import (
	"encoding/json"
	"fmt"

	"github.com/dtapps/allinssl_plugins/ctyun/openapi"
	"github.com/dtapps/allinssl_plugins/ctyun/types"
)

type Client struct {
	openapi     *openapi.Client // 请求
	endpoint    string          // 终端节点
	productCode string          // 产品类型
}

// 边缘安全加速平台
func NewClient(accessKeyId string, secretAccessKey string) (c *Client, err error) {
	c = &Client{
		endpoint:    "accessone-global.ctapi.ctyun.cn",
		productCode: "020",
	}
	c.openapi, err = openapi.NewClient(c.endpoint, accessKeyId, secretAccessKey)
	return c, err
}

// 域名基础及加速配置查询
// domain: 域名
// https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=113&api=13412&data=174&isNormal=1&vid=167
func (c *Client) GetQueryDomainInfo(domain string) (response types.CommonResponse[types.AccessoneQueryDomainInfoResponse], err error) {
	resp, err := c.openapi.R().
		SetContentType("application/json").
		SetQueryParam("domain", domain).              // 域名
		SetQueryParam("product_code", c.productCode). // 产品类型
		SetResult(&response).
		Get("/ctapi/v1/accessone/domain/config")
	if err != nil {
		return
	}
	if response.StatusCode != types.StatusCodeSuccess {
		err = fmt.Errorf("%s", resp.String())
		return
	}
	return
}

// 域名基础及加速配置修改
// domain：域名
// cert_name：证书备注名
// https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=113&api=13413&data=174&isNormal=1&vid=167
func (c *Client) PostUpdateDomainInfo(domain string, cert_name string) (response types.CommonResponse[any], err error) {
	paramsStr, err := json.Marshal(types.AccessoneUpdateDomainInfoRequest{
		Domain:      domain,        // 域名
		ProductCode: c.productCode, // 产品类型
		CertName:    cert_name,     // 证书备注名
	})
	if err != nil {
		return
	}
	resp, err := c.openapi.R().
		SetContentType("application/json").
		SetBody(string(paramsStr)).
		SetResult(&response).
		Post("/ctapi/v1/accessone/domain/modify_config")
	if err != nil {
		return
	}
	if response.StatusCode != types.StatusCodeSuccess {
		err = fmt.Errorf("%s", resp.String())
		return
	}
	return
}

// 查询证书详情
// name：证书备注名
// https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=113&api=13015&data=174&isNormal=1&vid=167
func (c *Client) GetQueryCertInfo(name string) (response types.CommonResponse[types.AccessoneQueryCertInfoResponse], err error) {
	resp, err := c.openapi.R().
		SetContentType("application/json").
		SetQueryParam("name", name). // 证书备注名
		SetResult(&response).
		Get("/ctapi/v1/accessone/cert/query")
	if err != nil {
		return
	}
	if response.StatusCode != types.StatusCodeSuccess {
		err = fmt.Errorf("%s", resp.String())
		return
	}
	return
}

// 创建证书
// name：证书备注名
// key：证书私钥
// certs：证书公钥
// https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=113&api=13014&data=174&isNormal=1&vid=167
func (c *Client) PostUpdateCertInfo(name string, key string, certs string) (response types.CommonResponse[types.AccessoneUpdateCertInfoResponse], err error) {
	paramsStr, err := json.Marshal(types.AccessoneUpdateCertInfoRequest{
		Name:  name,  // 证书备注名
		Key:   key,   // 证书私钥
		Certs: certs, // 证书公钥
	})
	if err != nil {
		return
	}
	resp, err := c.openapi.R().
		SetContentType("application/json").
		SetBody(string(paramsStr)).
		SetResult(&response).
		Post("/ctapi/v1/accessone/cert/create")
	if err != nil {
		return
	}
	if response.StatusCode != types.StatusCodeSuccess {
		err = fmt.Errorf("%s", resp.String())
		return
	}
	return
}

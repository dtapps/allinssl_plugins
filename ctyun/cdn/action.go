package cdn

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

// CDN加速
// https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=108&data=161&vid=154
func NewClient(accessKeyId string, secretAccessKey string) (c *Client, err error) {
	c = &Client{
		endpoint:    "https://ctcdn-global.ctapi.ctyun.cn",
		productCode: "008",
	}
	c.openapi, err = openapi.NewClient(c.endpoint, accessKeyId, secretAccessKey)
	return c, err
}

// 查询域名配置信息
// domain: 域名
// https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=108&api=11304&data=161&isNormal=1&vid=154
func (c *Client) GetQueryDomainInfo(domain string) (response types.CommonResponse[types.CdnQueryDomainInfoResponse], err error) {
	resp, err := c.openapi.R().
		SetContentType("application/json").
		SetQueryParam("product_code", c.productCode). // 产品类型
		SetQueryParam("domain", domain).              // 域名
		SetResult(&response).
		Get("/v1/domain/query-domain-detail")
	if err != nil {
		return
	}
	if response.StatusCode != types.StatusCodeSuccess {
		err = fmt.Errorf("%s", resp.String())
		return
	}
	return
}

// 修改域名配置
// domain: 域名
// cert_name: 证书备注名
// https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=108&api=11308&data=161&isNormal=1&vid=154
func (c *Client) PostUpdateDomainInfo(domain string, cert_name string) (response types.CommonResponse[any], err error) {
	paramsStr, err := json.Marshal(types.CdnUpdateDomainInfoRequest{
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
		Post("/v1/domain/update-domain")
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
// name: 证书备注名
// https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=108&api=10899&data=161&isNormal=1&vid=154
func (c *Client) GetQueryCertInfo(name string) (response types.CommonResponse[types.CdnQueryCertInfoResponse], err error) {
	resp, err := c.openapi.R().
		SetContentType("application/json").
		SetQueryParam("name", name). // 证书备注名
		SetResult(&response).
		Get("/v1/cert/query-cert-detail")
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
// name: 证书备注名
// key: 证书私钥
// certs: 证书公钥
// https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=108&api=10893&data=161&isNormal=1&vid=154
func (c *Client) PostUpdateCertInfo(name string, key string, certs string) (response types.CommonResponse[types.CdnUpdateCertInfoResponse], err error) {
	paramsStr, err := json.Marshal(types.CdnUpdateCertInfoRequest{
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
		Post("/v1/cert/creat-cert")
	if err != nil {
		return
	}
	if response.StatusCode != types.StatusCodeSuccess {
		err = fmt.Errorf("%s", resp.String())
		return
	}
	return
}

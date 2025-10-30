package teo

import (
	"context"
	"fmt"
	"time"

	"github.com/dtapps/allinssl_plugins/core"
	"github.com/dtapps/allinssl_plugins/edgeone/openapi"
	"github.com/dtapps/allinssl_plugins/edgeone/ssl"
	"github.com/dtapps/allinssl_plugins/edgeone/types"
	"go.dtapp.net/library/utils/gotime"
)

// 参数
type Params struct {
	ZoneID string // 加速域名所属站点 ID
	Domain string // 域名
}

// 返回
type Return struct {
	Bound  bool   // 是否已绑定
	CertID string // 证书 ID
}

// 域名绑定证书
// domain: 域名
// isBind: 是否已绑定
// 查询加速域名列表 https://cloud.tencent.com/document/product/1552/86336
// 配置域名证书 https://cloud.tencent.com/document/product/1552/80764
func Action(ctx context.Context, openapiClient *openapi.Client, certBundle *core.CertBundle, par *Params) (res *Return, err error) {

	if par == nil {
		return nil, fmt.Errorf("参数不能为空")
	}

	// 查询加速域名列表
	var dmainListResp types.TeoDescribeAccelerationDomainsResponse
	_, err = openapiClient.R().
		SetXTCEndpoint(Endpoint).
		SetXTCAction("DescribeAccelerationDomains").
		SetXTCVersion("2022-09-01").
		SetBodyMap(map[string]any{
			"ZoneId": par.ZoneID, // 加速域名所属站点 ID
			"Filters": []map[string]any{
				{
					"Name":   "domain-name", // 按照加速域名进行过滤
					"Values": []string{par.Domain},
				},
			}, // 过滤条件
			"Offset": 0,   // 分页偏移量
			"Limit":  200, // 每页数量
		}).
		SetResult(&dmainListResp).
		Post(Endpoint)
	if err != nil {
		return nil, fmt.Errorf("[%s]查询加速域名列表 错误: %w", par.Domain, err)
	}

	// 检查响应
	if dmainListResp.Response.Error.Code != "" {
		return nil, fmt.Errorf("[%s]查询加速域名列表 失败: %s", par.Domain, dmainListResp.Response.Error.Message)
	}

	// 检查域名证书是否配置了
	for _, dmainInfo := range dmainListResp.Response.AccelerationDomains {
		if dmainInfo.DomainName == par.Domain {
			for _, certInfo := range dmainInfo.Certificate.List {
				// 检查域名是否配置了现存证书
				if certBundle.IsGeneratedNote(certInfo.Alias) {
					expireTime := gotime.SetCurrentParse(certInfo.ExpireTime).Time
					if expireTime.After(time.Now()) {
						// 证书已存在且未过期
						return &Return{
							Bound:  true,
							CertID: certInfo.CertID,
						}, nil
					}
				}
			}
		}
	}

	// 上传证书
	sslResp, err := ssl.Action(ctx, openapiClient, certBundle, &ssl.Params{
		Domain: par.Domain,
	})
	if err != nil {
		return nil, err
	}

	// 配置域名证书
	var dmainConfigSslResp types.TeoModifyHostsCertificateResponse
	_, err = openapiClient.R().
		SetXTCEndpoint(Endpoint).
		SetXTCAction("ModifyHostsCertificate").
		SetXTCVersion("2022-09-01").
		SetBodyMap(map[string]any{
			"ZoneId": par.ZoneID,           // 站点 ID
			"Hosts":  []string{par.Domain}, // 需要修改证书配置的加速域名
			"Mode":   "sslcert",            // 配置服务端证书的模式
			"ServerCertInfo": []map[string]any{
				{
					"CertId": sslResp.CertID,
				},
			}, // SSL 证书配置
		}).
		SetResult(&dmainConfigSslResp).
		Post(Endpoint)
	if err != nil {
		return nil, fmt.Errorf("[%s]配置域名证书 错误: %w", par.Domain, err)
	}

	// 检查响应
	if dmainConfigSslResp.Response.Error.Code != "" {
		return nil, fmt.Errorf("[%s]配置域名证书 失败: %s", par.Domain, dmainConfigSslResp.Response.Error.Message)
	}

	return &Return{
		Bound:  true,
		CertID: sslResp.CertID,
	}, nil
}

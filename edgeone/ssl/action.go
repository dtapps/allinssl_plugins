package ssl

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/dtapps/allinssl_plugins/core"
	"github.com/dtapps/allinssl_plugins/edgeone/openapi"
	"github.com/dtapps/allinssl_plugins/edgeone/types"
	"go.dtapp.net/library/utils/gotime"
)

// 参数
type Params struct {
	Debug  bool   // 是否调试模式
	Domain string // 域名
}

// 返回
type Return struct {
	IsExist bool   // 表示原先存在证书
	CertID  string // 证书 ID
}

// 上传证书
// 获取证书列表 https://cloud.tencent.com/document/product/400/41671
// 上传证书 https://cloud.tencent.com/document/product/400/41665
func Action(ctx context.Context, openapiClient *openapi.Client, certBundle *core.CertBundle, par *Params) (res *Return, err error) {

	// 初始化请求参数
	var certListReq = map[string]any{
		"Offset": 0,    // 分页偏移量
		"Limit":  1000, // 每页数量
	}
	if par != nil {
		if par.Domain != "" {
			certListReq = map[string]any{
				"SearchKey": par.Domain, // 搜索关键词
				"Offset":    0,          // 分页偏移量
				"Limit":     1000,       // 每页数量
			}
		}
	}

	// 获取证书列表
	var certListResp types.SslDescribeCertificatesResponse
	_, err = openapiClient.R().
		SetXTCEndpoint(Endpoint).
		SetXTCAction("DescribeCertificates").
		SetXTCVersion("2019-12-05").
		SetBodyMap(certListReq).
		SetResult(&certListResp).
		SetContext(ctx).
		Post(Endpoint)
	if err != nil {
		return nil, fmt.Errorf("获取证书列表 错误: %w", err)
	}

	// 检查响应
	if certListResp.Response.Error.Code != "" {
		return nil, fmt.Errorf("获取证书列表 失败: %s", certListResp.Response.Error.Message)
	}

	// 检查证书是否已存在
	for _, certInfo := range certListResp.Response.Certificates {
		if par.Debug {
			slog.Info("[ssl] 比较",
				slog.Any("当前域名", certBundle.DNSNames),
				slog.String("接口域名", certInfo.Domain),
				slog.String("当前备注", certBundle.GetNoteShort()),
				slog.String("接口备注", certInfo.Alias),
				slog.Bool("结果", certBundle.IsSameCertificateNote(certBundle.GetNoteShort(), certInfo.Alias)),
			)
		}
		if certBundle.IsSameCertificateNote(certBundle.GetNoteShort(), certInfo.Alias) {
			expireTime := gotime.SetCurrentParse(certInfo.CERTEndTime).Time
			if expireTime.After(time.Now()) {
				// 证书已存在且未过期
				return &Return{
					IsExist: true,
					CertID:  certInfo.CertificateID,
				}, nil
			}
		}
	}

	// 上传证书（使用格式优化后的函数，确保证书链完整）
	privateKey, certificate := core.BuildCertsForAPIFormat(certBundle)
	certUpdateResp := types.SslUploadCertificateResponse{}
	_, err = openapiClient.R().
		SetXTCEndpoint(Endpoint).
		SetXTCAction("UploadCertificate").
		SetXTCVersion("2019-12-05").
		SetBodyMap(map[string]any{
			"CertificatePublicKey":  certificate,               // 证书内容（包含完整证书链）
			"CertificatePrivateKey": privateKey,                // 私钥内容
			"Alias":                 certBundle.GetNoteShort(), // 证书名称
			"CertificateType":       "SVR",                     // 证书类型，默认为服务器证书
			"Repeatable":            true,                      // 允许上传相同指纹的证书
		}).
		SetResult(&certUpdateResp).
		SetContext(ctx).
		Post(Endpoint)
	if err != nil {
		return nil, fmt.Errorf("上传证书 错误: %w", err)
	}

	// 检查响应
	if certUpdateResp.Response.Error.Code != "" {
		return nil, fmt.Errorf("上传证书 失败: %s", certUpdateResp.Response.Error.Message)
	}

	return &Return{
		CertID: certUpdateResp.Response.CertificateID,
	}, nil
}

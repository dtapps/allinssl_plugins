package ccms

import (
	"fmt"
	"strings"
	"time"

	"github.com/dtapps/allinssl_plugins/core"
	"github.com/dtapps/allinssl_plugins/ctyun/openapi"
	"github.com/dtapps/allinssl_plugins/ctyun/types"
)

// 上传证书
// isExist: 是否已存在
// 查询用户证书列表 https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=152&api=17233&data=204&isNormal=1&vid=283
// 上传证书 https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=152&api=17243&data=204&isNormal=1&vid=283
func Action(openapiClient *openapi.Client, certBundle *core.CertBundle) (isExist bool, err error) {

	// 1. 获取证书列表
	var certListResp types.CommonResponse[types.CcmsCertificateListResponse]
	_, err = openapiClient.R().
		SetBodyMap(map[string]any{
			"pageSize": 1,        // 当前页码
			"pageNum":  100,      // 每页记录数
			"origin":   "UPLOAD", // 证书来源
		}).
		SetResult(&certListResp).
		Post("/v1/certificate/list")
	if err != nil {
		return false, fmt.Errorf("获取证书列表错误: %w", err)
	}
	// 检查证书列表响应
	if certListResp.StatusCode != 200 {
		return false, fmt.Errorf("获取证书列表失败: %s", certListResp.Message)
	}
	for _, certInfo := range certListResp.ReturnObj.List {
		if strings.EqualFold(certInfo.Fingerprint, certBundle.GetFingerprintSHA1()) {
			var expireTime time.Time
			expireTime, err = time.Parse(time.RFC3339, certInfo.ExpireTime)
			if err != nil {
				return false, fmt.Errorf("Ccms 解析过期时间失败: %w", err)
			}
			if expireTime.After(time.Now()) {
				// 证书已存在且未过期
				return true, nil
			}
		}
	}

	// 2. 上传证书
	var certUpdateResp types.CommonResponse[any]
	_, err = openapiClient.R().
		SetBodyMap(map[string]any{
			"name":               certBundle.GetNoteShort(),   // 证书名称
			"certificate":        certBundle.Certificate,      // 证书字符串
			"privateKey":         certBundle.PrivateKey,       // 私钥字符串
			"certificateChain":   certBundle.CertificateChain, // 证书链字符串
			"encryptionStandard": "INTERNATIONAL",             // 加密标准
		}).
		SetResult(&certUpdateResp).
		Post("/v1/certificate/upload")
	if err != nil {
		err = fmt.Errorf("上传证书错误: %w", err)
		return
	}
	if certUpdateResp.StatusCode != 200 {
		return false, fmt.Errorf("上传证书失败: %s", certUpdateResp.Message)
	}

	return false, nil
}

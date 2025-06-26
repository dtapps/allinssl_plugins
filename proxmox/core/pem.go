package core

import (
	"crypto/sha1"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"time"
)

// CertBundle 表示从PEM文件中提取的证书和私钥
type CertBundle struct {
	Certificate      string `json:"-"` // 证书字符串
	PrivateKey       string `json:"-"` // 私钥字符串
	CertificateChain string `json:"-"` // 证书链字符串

	SerialNumber       string    `json:"serialNumber"`       // 证书序列号
	FingerprintSHA1    string    `json:"fingerprintSHA1"`    // 证书SHA1指纹
	FingerprintSHA256  string    `json:"fingerprintSHA256"`  // 证书SHA256指纹
	NotBefore          time.Time `json:"notBefore"`          // 证书生效时间
	NotAfter           time.Time `json:"notAfter"`           // 证书过期时间
	Subject            string    `json:"subject"`            // 证书主题
	Issuer             string    `json:"issuer"`             // 颁发者
	DNSNames           []string  `json:"dnsNames"`           // 域名列表
	EmailAddresses     []string  `json:"emailAddresses"`     // 邮箱地址
	IPAddresses        []string  `json:"ipAddresses"`        // IP地址
	SignatureAlgorithm string    `json:"signatureAlgorithm"` // 签名算法
}

func ParseCertBundle(certPEMData, keyPEMData []byte) (*CertBundle, error) {
	// 解析主证书
	block, rest := pem.Decode(certPEMData)
	if block == nil || block.Type != "CERTIFICATE" {
		return nil, fmt.Errorf("invalid certificate PEM")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}

	// 计算证书SHA1指纹
	sha1Hash := sha1.Sum(cert.Raw)                     // cert.Raw 包含证书的 DER 编码字节
	fingerprintSHA1 := hex.EncodeToString(sha1Hash[:]) // 转换为十六进制字符串

	// 计算证书SHA256指纹
	sha256Hash := sha256.Sum256(cert.Raw)                  // 使用 sha256.Sum256
	fingerprintSHA256 := hex.EncodeToString(sha256Hash[:]) // 转换为十六进制字符串

	// 提取主证书字符串（第一个证书）
	mainCertPEM := string(pem.EncodeToMemory(block))

	// 提取证书链（剩下的部分）
	var chainPEM string
	for len(rest) > 0 {
		block, rest = pem.Decode(rest)
		if block == nil || block.Type != "CERTIFICATE" {
			continue // 跳过非证书内容
		}
		chainPEM += string(pem.EncodeToMemory(block))
	}

	// 转换 IP 地址为字符串
	ipStrings := make([]string, 0, len(cert.IPAddresses))
	for _, ip := range cert.IPAddresses {
		ipStrings = append(ipStrings, ip.String())
	}

	return &CertBundle{
		Certificate:        mainCertPEM,
		PrivateKey:         string(keyPEMData),
		CertificateChain:   chainPEM,
		SerialNumber:       cert.SerialNumber.String(),
		FingerprintSHA1:    fingerprintSHA1,
		FingerprintSHA256:  fingerprintSHA256,
		NotBefore:          cert.NotBefore,
		NotAfter:           cert.NotAfter,
		Subject:            cert.Subject.String(),
		Issuer:             cert.Issuer.String(),
		DNSNames:           cert.DNSNames,
		EmailAddresses:     cert.EmailAddresses,
		IPAddresses:        ipStrings,
		SignatureAlgorithm: cert.SignatureAlgorithm.String(),
	}, nil
}

const notePrefix = "allinssl-"

// 获取证书名字
func (cb *CertBundle) GetNote() string {
	return fmt.Sprintf("allinssl-%s", cb.FingerprintSHA256)
}

// 获取证书名字（缩短的，天翼云证书管理在使用）
func (cb *CertBundle) GetNoteShort() string {
	return fmt.Sprintf("allinssl-%s", cb.FingerprintSHA256[:6])
}

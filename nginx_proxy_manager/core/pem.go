package core

import (
	"crypto/x509"
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

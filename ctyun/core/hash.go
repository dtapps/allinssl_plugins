package core

import (
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
)

func GetSHA256(certStr string) (string, error) {
	certPEM := []byte(certStr)
	block, _ := pem.Decode(certPEM)
	if block == nil {
		return "", fmt.Errorf("无法解析证书 PEM")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("解析证书失败: %v", err)
	}

	sha256Hash := sha256.Sum256(cert.Raw)
	return hex.EncodeToString(sha256Hash[:]), nil
}

package core

import (
	"crypto/sha1"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"reflect"
	"regexp"
	"sort"
	"strings"
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
	EmailAddresses     []string  `json:"-"`                  // 邮箱地址
	IPAddresses        []string  `json:"-"`                  // IP地址
	SignatureAlgorithm string    `json:"signatureAlgorithm"` // 签名算法

	certRaw           []byte
	fingerprintSHA1   string // 证书SHA1指纹
	fingerprintSHA256 string // 证书SHA256指纹
}

// NewCertBundle 构造一个 CertBundle，支持主证书 + 私钥 + 可选中间证书链
func NewCertBundle(certPEMData, keyPEMData []byte, chainPEMData ...[]byte) (*CertBundle, error) {
	var fullCertPEM []byte
	fullCertPEM = append(fullCertPEM, certPEMData...)
	for _, chain := range chainPEMData {
		fullCertPEM = append(fullCertPEM, '\n')
		fullCertPEM = append(fullCertPEM, chain...)
	}
	return ParseCertBundle(fullCertPEM, keyPEMData)
}

// ParseCertBundle 从PEM编码的证书和私钥数据中解析证书和私钥
func ParseCertBundle(certPEMData, keyPEMData []byte) (*CertBundle, error) {

	// 解析主证书
	block, rest := pem.Decode(certPEMData)
	if block == nil || block.Type != "CERTIFICATE" {
		return nil, fmt.Errorf("invalid certificate PEM")
	}

	// 解析证书
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
		certRaw:            cert.Raw,
	}, nil
}

// FingerprintSHA1 计算证书的 SHA1 指纹
func (cb *CertBundle) GetFingerprintSHA1() string {
	if cb.fingerprintSHA1 != "" {
		return cb.fingerprintSHA1
	}
	hash := sha1.Sum(cb.certRaw)
	return hex.EncodeToString(hash[:])
}

// FingerprintSHA256 计算证书的 SHA256 指纹
func (cb *CertBundle) GetFingerprintSHA256() string {
	if cb.fingerprintSHA256 != "" {
		return cb.fingerprintSHA256
	}
	hash := sha256.Sum256(cb.certRaw)
	return hex.EncodeToString(hash[:])
}

// IsWildcard 判断该证书是否包含泛域名
func (cb *CertBundle) IsWildcard() bool {
	for _, name := range cb.DNSNames {
		if isWildcardHost(name) {
			return true
		}
	}
	return false
}

// isWildcardHost 检查证书域名是否为通配符域名
func isWildcardHost(host string) bool {
	return len(host) >= 3 && host[0] == '*' && host[1] == '.'
}

// CanDomainsUseCert 判断指定域名是否被证书覆盖
func (cb *CertBundle) CanDomainsUseCert(domains []string) bool {
	return CanDomainsUseCert(domains, cb.DNSNames)
}

// CanDomainsUseCert 判断指定域名是否能被 DNSNames 覆盖
func CanDomainsUseCert(domains, dnsNames []string) bool {
	for _, domain := range domains {
		matched := false
		for _, certDomain := range dnsNames {
			if domain == certDomain || matchWildcard(certDomain, domain) {
				matched = true
				break
			}
		}
		if !matched {
			return false
		}
	}
	return true
}

// matchWildcard 判断证书中的泛域名（如 *.example.com）是否能匹配用户域名
func matchWildcard(certDomain, userDomain string) bool {
	if strings.HasPrefix(certDomain, "*.") {
		parent := certDomain[2:] // 去掉 "*.example.com" 的 "*."
		return isSubDomainOf(userDomain, parent)
	}
	return false
}

// isSubDomainOf 判断 userDomain 是否是 parent 的子域名
func isSubDomainOf(domain, parent string) bool {
	if !strings.HasSuffix(domain, "."+parent) {
		return false
	}
	// 确保 domain 比 parent 多至少一个字符（即前面有 .）
	return len(domain) > len(parent)+1
}

// IsExpired 判断证书是否已过期
func (cb *CertBundle) IsExpired() bool {
	now := time.Now()
	return now.Before(cb.NotBefore) || now.After(cb.NotAfter)
}

// IsDNSNamesMatch 判断传入的域名列表是否与证书的 DNSNames 完全一致（顺序可不同）
func (cb *CertBundle) IsDNSNamesMatch(domains []string) bool {
	if len(cb.DNSNames) != len(domains) {
		return false // 长度不一致，直接返回 false
	}

	// 复制两个切片并排序比较
	cbCopy := make([]string, len(cb.DNSNames))
	copy(cbCopy, cb.DNSNames)
	sort.Strings(cbCopy)

	domainsCopy := make([]string, len(domains))
	copy(domainsCopy, domains)
	sort.Strings(domainsCopy)

	for i := 0; i < len(cbCopy); i++ {
		if cbCopy[i] != domainsCopy[i] {
			return false
		}
	}
	return true
}

const notePrefix = "allinssl-"

// GetNote 获取证书名字
// 旧的
func (cb *CertBundle) GetNote() string {
	fp := cb.GetFingerprintSHA256()
	return fmt.Sprintf("%s%s", notePrefix, fp)
}

// GetNoteShort 获取证书名字（缩短的，天翼云、南墙WEB应用防火墙、Lucky证书管理在使用）
// 新的
func (cb *CertBundle) GetNoteShort() string {
	fp := cb.GetFingerprintSHA256()
	if len(fp) < 6 {
		return fmt.Sprintf("%s%s", notePrefix, fp)
	}
	return fmt.Sprintf("%s%s", notePrefix, fp[:6])
}

// IsGeneratedNote 判断传入的 note 字符串是否由 GetNote 或 GetNoteShort 方法生成。
// 它会自动忽略字符串首尾的空格。
func (cb *CertBundle) IsGeneratedNote(note string) bool {
	// 1. 去除首尾空格
	note = strings.TrimSpace(note)

	// 2. 检查是否以 notePrefix 开头
	if !strings.HasPrefix(note, notePrefix) {
		return false
	}

	// 3. 提取指纹部分
	fpPart := note[len(notePrefix):]

	// 4. 定义一个正则表达式，用于匹配十六进制字符串
	// SHA256 指纹是 64 位，缩短版是 6 位，都必须是小写的 [0-9a-f]
	hexRegex := regexp.MustCompile(`^[0-9a-f]+$`)

	// 5. 验证：长度必须是 64 或 6，并且内容必须是合法的十六进制
	return (len(fpPart) == 64 || len(fpPart) == 6) && hexRegex.MatchString(fpPart)
}

// IsSameCertificateNote 判断两个备注是否属于同一个证书
// 即使一个是长版本（64位）一个是短版本（6位），只要短版本是长版本的前6位，就认为是同一个证书
func (cb *CertBundle) IsSameCertificateNote(note1, note2 string) bool {
	// 1. 去除首尾空格
	note1 = strings.TrimSpace(note1)
	note2 = strings.TrimSpace(note2)

	// 2. 验证两个备注都符合格式要求
	if !cb.IsGeneratedNote(note1) || !cb.IsGeneratedNote(note2) {
		return false
	}

	// 3. 提取指纹部分
	fp1 := note1[len(notePrefix):]
	fp2 := note2[len(notePrefix):]

	// 4. 比较：如果一个是6位，一个是64位，并且6位是64位的前6位，则认为相同
	if len(fp1) == 6 && len(fp2) == 64 {
		return fp2[:6] == fp1
	}
	if len(fp2) == 6 && len(fp1) == 64 {
		return fp1[:6] == fp2
	}

	// 5. 否则直接比较是否完全相同
	return fp1 == fp2
}

// VerifyChain 检查证书链是否完整、有效
func (cb *CertBundle) VerifyChain() error {
	// 解析主证书
	cert, err := x509.ParseCertificate(cb.certRaw)
	if err != nil {
		return fmt.Errorf("failed to parse main certificate: %w", err)
	}

	// 构造中间证书池
	intermediates := x509.NewCertPool()
	if cb.CertificateChain != "" {
		rest := []byte(cb.CertificateChain)
		for {
			var block *pem.Block
			block, rest = pem.Decode(rest)
			if block == nil {
				break
			}
			if block.Type != "CERTIFICATE" {
				continue
			}
			ic, err := x509.ParseCertificate(block.Bytes)
			if err != nil {
				continue
			}
			intermediates.AddCert(ic)
		}
	}

	// 系统根证书
	roots, err := x509.SystemCertPool()
	if err != nil {
		return fmt.Errorf("failed to load system root CAs: %w", err)
	}

	opts := x509.VerifyOptions{
		Intermediates: intermediates,
		Roots:         roots,
	}

	// 验证链
	if _, err := cert.Verify(opts); err != nil {
		return fmt.Errorf("certificate chain is invalid: %w", err)
	}
	return nil
}

// IsChainValid 检查证书链是否完整、有效，返回 true 或 false
func (cb *CertBundle) IsChainValid() bool {
	return cb.VerifyChain() == nil
}

// BuildCertsForAPI 组合第三方接口需要的 key 和 certs
func BuildCertsForAPI(certBundle *CertBundle) (key, certs string) {
	key = certBundle.PrivateKey

	// certs = 主证书 + 中间证书链
	if certBundle.CertificateChain != "" {
		certs = certBundle.Certificate + "\n" + certBundle.CertificateChain
	} else {
		certs = certBundle.Certificate
	}

	return key, certs
}

// BuildCertsForAPIFormat 组合第三方接口需要的 key 和 certs，确保格式完全符合API要求
// 严格按照用户要求的格式：网站证书 -> CA中间证书 -> CA根证书，每个证书块之间使用两个换行符
func BuildCertsForAPIFormat(certBundle *CertBundle) (key, certs string) {
	// 严格处理私钥格式，确保完全符合PEM标准
	key = processPrivateKey(certBundle.PrivateKey)

	// 按照用户要求的格式直接构建证书链
	var result strings.Builder

	// 1. 首先添加网站证书（主证书）
	mainCert := processSingleCertificate(certBundle.Certificate)
	if mainCert != "" {
		result.WriteString(mainCert)
	}

	// 2. 然后处理证书链，分离中间证书和根证书
	if certBundle.CertificateChain != "" {
		// 分割证书链中的多个证书
		chainCerts := splitAndProcessCertificates(certBundle.CertificateChain)

		var intermediateCertFound bool
		var rootCertFound bool

		// 先查找所有证书，识别中间证书和根证书
		var intermediateCert string
		var rootCert string

		for _, cert := range chainCerts {
			// 尝试解析证书以确定类型
			block, _ := pem.Decode([]byte(cert))
			if block != nil {
				x509Cert, err := x509.ParseCertificate(block.Bytes)
				if err == nil {
					// 如果颁发者和主题相同，是根证书
					if reflect.DeepEqual(x509Cert.Issuer, x509Cert.Subject) {
						rootCert = cert
						rootCertFound = true
					} else {
						// 否则是中间证书
						intermediateCert = cert
						intermediateCertFound = true
					}
				}
			}
		}

		// 按照用户要求的顺序和格式添加证书
		// 2. 添加CA中间证书，前面加两个换行符
		if intermediateCertFound {
			result.WriteString("\n\n")
			result.WriteString(intermediateCert)
		}

		// 3. 添加CA根证书，前面加两个换行符
		if rootCertFound {
			result.WriteString("\n\n")
			result.WriteString(rootCert)
		}
	}

	certs = result.String()
	return key, certs
}

// processPrivateKey 严格处理私钥格式，确保符合PEM标准
func processPrivateKey(privateKey string) string {
	// 去除所有空白字符
	privateKey = strings.TrimSpace(privateKey)
	if privateKey == "" {
		return ""
	}

	// 识别私钥类型
	var beginMark, endMark string
	if strings.Contains(privateKey, "RSA PRIVATE KEY") {
		beginMark = "-----BEGIN RSA PRIVATE KEY-----"
		endMark = "-----END RSA PRIVATE KEY-----"
	} else if strings.Contains(privateKey, "PRIVATE KEY") {
		beginMark = "-----BEGIN PRIVATE KEY-----"
		endMark = "-----END PRIVATE KEY-----"
	} else {
		// 保留原始格式
		return privateKey
	}

	// 提取私钥内容（去掉BEGIN和END标记）
	content := privateKey
	if strings.HasPrefix(content, beginMark) {
		content = content[len(beginMark):]
	}
	if strings.HasSuffix(content, endMark) {
		content = content[:len(content)-len(endMark)]
	}

	// 去除内容中的空白字符，只保留Base64编码部分
	content = strings.TrimSpace(content)
	content = strings.ReplaceAll(content, "\n", "")
	content = strings.ReplaceAll(content, "\r", "")
	content = strings.ReplaceAll(content, "\t", "")
	content = strings.ReplaceAll(content, " ", "")

	// 重新格式化私钥，确保每行64个字符
	var formattedContent strings.Builder
	for i := 0; i < len(content); i += 64 {
		end := i + 64
		if end > len(content) {
			end = len(content)
		}
		formattedContent.WriteString(content[i:end])
		formattedContent.WriteString("\n")
	}

	// 构建完整的私钥格式，确保格式完全正确
	var result strings.Builder
	result.WriteString(beginMark)
	result.WriteString("\n")
	result.WriteString(formattedContent.String())
	result.WriteString(endMark)

	return result.String()
}

// processSingleCertificate 处理单个证书，确保格式完全正确
func processSingleCertificate(cert string) string {
	// 去除所有空白字符和换行符，然后重新格式化
	cert = strings.TrimSpace(cert)
	if cert == "" {
		return ""
	}

	// 确保证书有正确的BEGIN和END标记
	beginMark := "-----BEGIN CERTIFICATE-----"
	endMark := "-----END CERTIFICATE-----"

	// 提取证书内容（去掉BEGIN和END标记）
	content := cert
	if strings.HasPrefix(content, beginMark) {
		content = content[len(beginMark):]
	}
	if strings.HasSuffix(content, endMark) {
		content = content[:len(content)-len(endMark)]
	}

	// 去除内容中的空白字符，只保留Base64编码部分
	content = strings.TrimSpace(content)
	content = strings.ReplaceAll(content, "\n", "")
	content = strings.ReplaceAll(content, "\r", "")
	content = strings.ReplaceAll(content, "\t", "")
	content = strings.ReplaceAll(content, " ", "")

	// 重新格式化证书，确保每行64个字符
	var formattedContent strings.Builder
	for i := 0; i < len(content); i += 64 {
		end := i + 64
		if end > len(content) {
			end = len(content)
		}
		if i > 0 {
			formattedContent.WriteString("\n")
		}
		formattedContent.WriteString(content[i:end])
	}

	// 构建完整的证书格式，确保格式完全正确
	var result strings.Builder
	result.WriteString(beginMark)
	result.WriteString("\n")
	result.WriteString(formattedContent.String())
	result.WriteString("\n")
	result.WriteString(endMark)

	return result.String()
}

// splitAndProcessCertificates 分割并处理多个证书
func splitAndProcessCertificates(certs string) []string {
	var result []string

	// 使用正则表达式分割证书
	certPattern := regexp.MustCompile(`-----BEGIN CERTIFICATE-----[\s\S]*?-----END CERTIFICATE-----`)
	matches := certPattern.FindAllString(certs, -1)

	for _, match := range matches {
		processed := processSingleCertificate(match)
		if processed != "" {
			result = append(result, processed)
		}
	}

	return result
}

// 证书信息
func (cb *CertBundle) ResultInfo() map[string]any {
	return map[string]any{
		"serialNumber": cb.SerialNumber,                    // 证书序列号
		"notBefore":    cb.NotBefore.Format(time.DateTime), // 证书生效时间
		"notAfter":     cb.NotAfter.Format(time.DateTime),  // 证书过期时间
		"dnsNames":     cb.DNSNames,                        // 域名列表
		"verifyChain":  cb.IsChainValid(),                  // 证书链是否完整
	}
}

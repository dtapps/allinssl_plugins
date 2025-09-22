package core

// APIBundle 表示一对证书（本地证书 + API证书）
type APICertBundle struct {
	Local *CertBundle // 本地证书
	API   *CertBundle // API证书
}

// LoadApiCert 在已有的 CertBundle 基础上，加载 API 证书并返回 APIBundle
func (cb *CertBundle) LoadApiCert(certPEMData, keyPEMData []byte) (*APICertBundle, error) {
	apiBundle, err := ParseCertBundle(certPEMData, keyPEMData)
	if err != nil {
		return nil, err
	}
	return &APICertBundle{
		Local: cb,
		API:   apiBundle,
	}, nil
}

// LoadApiCertWithChain 在已有的 CertBundle 基础上，加载 API 证书（可带中间证书链）并返回 APICertBundle
func (cb *CertBundle) LoadApiCertWithChain(certPEMData, keyPEMData, chainPEMData []byte) (*APICertBundle, error) {
	// 拼接证书 PEM + 中间链，确保有换行分隔
	fullCertPEM := make([]byte, 0, len(certPEMData)+len(chainPEMData)+1)
	fullCertPEM = append(fullCertPEM, certPEMData...)
	if len(chainPEMData) > 0 {
		fullCertPEM = append(fullCertPEM, '\n')
		fullCertPEM = append(fullCertPEM, chainPEMData...)
	}

	apiBundle, err := ParseCertBundle(fullCertPEM, keyPEMData)
	if err != nil {
		return nil, err
	}

	return &APICertBundle{
		Local: cb,
		API:   apiBundle,
	}, nil
}

// 两个证书是否一致
func (ab *APICertBundle) IsCertsEqual() bool {
	if ab.Local == nil || ab.API == nil {
		return false
	}
	return ab.Local.GetFingerprintSHA256() == ab.API.GetFingerprintSHA256()
}

// API 证书是否过期
func (ab *APICertBundle) IsAPIExpired() bool {
	if ab.API == nil {
		return true
	}
	return ab.API.IsExpired()
}

// API 证书是否包含指定域名（支持 *.abc.com）
func (ab *APICertBundle) IsAPICertContainsDomain(domains []string) bool {
	if ab.API == nil {
		return false
	}
	return ab.API.CanDomainsUseCert(domains)
}

package core

import "strings"

// ParseDomainsFixedSeparator 按固定分隔符解析域名字符串，判断是否为多个域名
func ParseDomainsFixedSeparator(input, separator string) ([]string, bool) {
	if input == "" || separator == "" {
		return nil, false
	}

	parts := strings.Split(input, separator)
	var domains []string

	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			domains = append(domains, trimmed)
		}
	}

	isMultiple := len(domains) > 1
	return domains, isMultiple
}

package types

import "time"

// 获取域名列表 响应参数
type SiteQueryDomainListResponse struct {
	Data  []SiteQueryDomainListDataResponse `json:"data"`  // 域名列表
	Total int                               `json:"total"` // 总数
}

type SiteQueryDomainListDataResponse struct {
	ID          int      `json:"id"`           // 网站ID
	CertID      int      `json:"cert_id"`      // 证书ID
	ServerNames []string `json:"server_names"` // 网站域名
	Upstreams   []string `json:"upstreams"`    // 上游服务器
	Ports       []string `json:"ports"`        // 端口列表

	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	GroupID     int       `json:"group_id"`
	Comment     string    `json:"comment"`
	IsEnabled   bool      `json:"is_enabled"`
	LoadBalance struct {
		BalanceType int `json:"balance_type"`
	} `json:"load_balance"`
	ExcludePaths             any    `json:"exclude_paths"`
	ExcludeContentType       any    `json:"exclude_content_type"`
	Mode                     int    `json:"mode"`
	Static                   bool   `json:"static"`
	Type                     int    `json:"type"`
	Index                    string `json:"index"`
	StaticDefault            int    `json:"static_default"`
	Init                     bool   `json:"init"`
	RedirectStatusCode       int    `json:"redirect_status_code"`
	CertType                 int    `json:"cert_type"`
	CertFilename             string `json:"cert_filename"`
	KeyFilename              string `json:"key_filename"`
	Email                    string `json:"email"`
	Title                    string `json:"title"`
	Icon                     string `json:"icon"`
	ACLResponseStatusCode    int    `json:"acl_response_status_code"`
	ACLResponseHTMLPath      string `json:"acl_response_html_path"`
	ForbiddenStatusCode      int    `json:"forbidden_status_code"`
	ForbiddenHTMLPath        string `json:"forbidden_html_path"`
	NotFoundStatusCode       int    `json:"not_found_status_code"`
	NotFoundHTMLPath         string `json:"not_found_html_path"`
	OfflineStatusCode        int    `json:"offline_status_code"`
	OfflineHTMLPath          string `json:"offline_html_path"`
	BadGatewayStatusCode     int    `json:"bad_gateway_status_code"`
	BadGatewayHTMLPath       string `json:"bad_gateway_html_path"`
	GatewayTimeoutStatusCode int    `json:"gateway_timeout_status_code"`
	GatewayTimeoutHTMLPath   string `json:"gateway_timeout_html_path"`
	AuthDefenseID            int    `json:"auth_defense_id"`
	ChallengeID              int    `json:"challenge_id"`
	ChaosID                  int    `json:"chaos_id"`
	ChaosIsEnabled           bool   `json:"chaos_is_enabled"`
	AccessLogLimit           int    `json:"access_log_limit"`
	ErrorLogLimit            int    `json:"error_log_limit"`
	ACLEnabled               bool   `json:"acl_enabled"`
	TamperRefresh            int    `json:"tamper_refresh"`
	TamperRefreshState       string `json:"tamper_refresh_state"`
	WRID                     int    `json:"wr_id"`
	CustomLocation           any    `json:"custom_location"`
	HealthCheck              bool   `json:"health_check"`
	Portal                   bool   `json:"portal"`
	PortalRedirect           string `json:"portal_redirect"`
	Position                 int    `json:"position"`
	StatEnabled              bool   `json:"stat_enabled"`
	ReqValue                 int    `json:"req_value"`
	DeniedValue              int    `json:"denied_value"`
	HealthState              map[string]struct {
		State int    `json:"state"`
		Error string `json:"error"`
	} `json:"health_state"`
	CCBot     bool `json:"cc_bot"`
	Semantics bool `json:"semantics"`
}

// 获取证书列表 响应参数
type SiteQueryCertListResponse struct {
	Nodes []struct {
		ID           int      `json:"id"`            // 证书ID
		Domains      []string `json:"domains"`       // 证书域名
		RelatedSites []string `json:"related_sites"` // 关联域名
	} `json:"nodes"` // 域名列表
	Total int `json:"total"` // 总数
}

// 获取证书详情 响应参数
type SiteQueryCertInfoResponse struct {
	Acme struct {
		Domains []string `json:"domains"` // 证书域名
	} `json:"acme"`
	Manual struct {
		Crt string `json:"crt"`
		Key string `json:"key"`
	} `json:"manual"`
}

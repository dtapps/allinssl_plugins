package types

// 获取所有网站列表 响应参数
// https://www.yuque.com/xiaopimianban/iqqphe/scrqpgw5snhoxg7a
type SiteQueryDomainListResponse struct {
	ID       int    `json:"id"`       // 网站ID，操作网站相关时用到
	Name     string `json:"name"`     // 网站名称
	Lang     string `json:"lang"`     // 网站语言类型
	NodeDemo string `json:"nodeDemo"` // 网站备注
}

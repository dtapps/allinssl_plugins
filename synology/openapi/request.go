package openapi

import (
	"crypto/tls"
	"fmt"
	"net/url"

	"resty.dev/v3"
)

type Client struct {
	*resty.Client
	username string // 用户名
	password string // 密码
	token    string // 令牌
}

// NewClient 创建请求客户端
// https://kb.synology.cn/zh-cn/DG/DSM_Login_Web_API_Guide/1
func NewClient(baseURL, username, password string) (*Client, error) {
	if baseURL == "" {
		return nil, fmt.Errorf("check baseURL")
	}
	if _, err := url.Parse(baseURL); err != nil {
		return nil, fmt.Errorf("check baseURL: %w", err)
	}

	// 安全地确保 baseURL 末尾是 /webapi/entry.cgi
	baseURL, err := ensureAPIPath(baseURL)
	if err != nil {
		return nil, err
	}

	client := resty.New().SetBaseURL(baseURL)
	client.SetResponseMiddlewares(
		Ensure2xxResponseMiddleware,       // 先调用，判断状态是不是请求成功
		resty.AutoParseResponseMiddleware, // 必须放后面，才能先判断状态码再解析
	)

	return &Client{
		Client:   client,
		username: username,
		password: password,
	}, nil
}

// WithDebug 开启调试模式
func (c *Client) WithDebug() *Client {
	c.EnableDebug()
	return c
}

// WithSkipVerify 跳过验证
func (c *Client) WithSkipVerify() *Client {
	c.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	return c
}

// WithLogin 登录
// https://kb.synology.cn/zh-cn/DG/DSM_Login_Web_API_Guide/2#x_anchor_iddbcc293edb
func (c *Client) WithLogin() (*Client, error) {
	if c.token != "" {
		return c, nil
	}
	var loginResp struct {
		Data struct {
			Did          string `json:"did"`            // 会话ID
			IsPortalPort bool   `json:"is_portal_port"` // 是否通过公网端口访问
			Sid          string `json:"sid"`            // 设备唯一标识
			Synotoken    string `json:"synotoken"`      // 防止 CSRF 攻击的安全 Token
		} `json:"data"`
		Success bool `json:"success"` // 是否成功
	}
	_, err := c.R().
		SetQueryParam("api", "SYNO.API.Auth").
		SetQueryParam("version", "6").
		SetQueryParam("method", "login").
		SetQueryParam("account", c.username).
		SetQueryParam("passwd", c.password).
		SetQueryParam("enable_syno_token", "yes").
		SetResult(&loginResp).
		Get("")
	if err != nil {
		return nil, fmt.Errorf("login failed: %v", err)
	}

	// 检查令牌是否为空
	if loginResp.Data.Synotoken == "" {
		return nil, fmt.Errorf("token is empty")
	}
	c.token = loginResp.Data.Synotoken

	return c, nil
}

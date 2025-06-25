package openapi

import (
	"fmt"
	"net/url"
	"time"

	"resty.dev/v3"
)

type Client struct {
	*resty.Client
	email    string    // 邮箱
	password string    // 密码
	token    string    // 令牌
	tokenExp time.Time // 令牌 过期时间
}

// NewClient 创建请求客户端
func NewClient(baseURL string, email string, password string) (*Client, error) {
	if baseURL == "" {
		return nil, fmt.Errorf("check baseURL")
	}
	if _, err := url.Parse(baseURL); err != nil {
		return nil, fmt.Errorf("check baseURL: %w", err)
	}
	if email == "" {
		return nil, fmt.Errorf("check email")
	}
	if password == "" {
		return nil, fmt.Errorf("check password")
	}

	// 安全地确保 baseURL 末尾是 /api
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
		email:    email,
		password: password,
		token:    "",          // 初始化为空字符串
		tokenExp: time.Time{}, // 初始化为零值时间
	}, nil
}

// WithDebug 开启调试模式
func (c *Client) WithDebug() *Client {
	c.EnableDebug()
	return c
}

// WithLogin 登录
func (c *Client) WithLogin() (*Client, error) {
	if c.token != "" && c.tokenExp.After(time.Now()) {
		return c, nil
	}
	var tokensResp struct {
		Token   string `json:"token"`
		Expires string `json:"expires"`
	}
	type tokensReq struct {
		Identity string `json:"identity"`
		Secret   string `json:"secret"`
	}
	_, err := c.R().
		SetContentType("application/json").
		SetBody(map[string]string{
			"identity": c.email,
			"secret":   c.password,
		}).
		SetResult(&tokensResp).
		Post("/tokens")
	if err != nil {
		return nil, fmt.Errorf("login failed: %v", err)
	}

	// 检查令牌是否为空
	if tokensResp.Token == "" {
		return nil, fmt.Errorf("token is empty")
	}
	c.token = tokensResp.Token

	// 解析过期时间
	expires, err := time.Parse(time.RFC3339, tokensResp.Expires)
	if err != nil {
		return nil, fmt.Errorf("parse expires: %v", err)
	}
	c.tokenExp = expires

	return c, nil
}

// WithDebug 设置令牌
func (c *Client) WithToken() *Client {
	c.SetHeader("Authorization", "Bearer "+c.token)
	return c
}

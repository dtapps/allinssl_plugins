package openapi

import (
	"crypto/tls"
	"fmt"
	"net/url"

	"github.com/dtapps/allinssl_plugins/uuwaf/types"
	"resty.dev/v3"
)

type Client struct {
	*resty.Client
	username string // 邮箱
	password string // 密码
	token    string // 令牌
}

// NewClient 创建请求客户端
// http://xxxx:xx/api/schema
func NewClient(baseURL string, username string, password string) (*Client, error) {
	if _, err := url.Parse(baseURL); err != nil {
		return nil, fmt.Errorf("check baseURL: %w", err)
	}

	// 安全地确保 baseURL 末尾是 /api/v1/
	baseURL, err := ensureAPIPath(baseURL)
	if err != nil {
		return nil, err
	}

	client := resty.New().SetBaseURL(baseURL)
	client.SetResponseMiddlewares(
		Ensure2xxResponseMiddleware,       // 先调用，判断状态是不是请求成功
		resty.AutoParseResponseMiddleware, // 再调用，才能先判断状态码再解析
	)

	return &Client{
		Client:   client,
		username: username,
		password: password,
		token:    "", // 初始化为空字符串
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
func (c *Client) WithLogin() (*Client, error) {
	if c.token != "" {
		return c, nil
	}
	var loginResp types.LoginResponse
	_, err := c.R().
		SetBody(map[string]string{
			"usr": c.username,
			"pwd": c.password,
		}).
		SetResult(&loginResp).
		SetContentType("application/json").
		Post("/users/login")
	if err != nil {
		return nil, fmt.Errorf("login failed: %v", err)
	}

	// 检查令牌是否为空
	if loginResp.Token == "" {
		return nil, fmt.Errorf("token is empty")
	}
	c.token = loginResp.Token

	return c, nil
}

// WithDebug 设置令牌
func (c *Client) WithV6Token() *Client {
	c.SetHeader("Authorization", c.token)
	return c
}

// WithDebug 设置令牌
func (c *Client) WithV7Token() *Client {
	c.SetAuthToken(c.token)
	return c
}

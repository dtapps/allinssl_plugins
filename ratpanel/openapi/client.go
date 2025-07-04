package openapi

import (
	"crypto/tls"
	"fmt"
	"net/url"

	"resty.dev/v3"
)

type Client struct {
	*resty.Client
}

// NewClient 创建请求客户端
func NewClient(baseURL string, userID int, userToken string) (*Client, error) {
	if _, err := url.Parse(baseURL); err != nil {
		return nil, fmt.Errorf("check baseURL: %w", err)
	}

	// 安全地确保 baseURL 末尾是 /api
	baseURL, err := ensureAPIPath(baseURL)
	if err != nil {
		return nil, err
	}

	client := resty.New().SetBaseURL(baseURL)
	client.SetRequestMiddlewares(
		resty.PrepareRequestMiddleware,          // 先调用，创建 RawRequest
		PreRequestMiddleware(userID, userToken), // 再调用，自定义中间
	)
	client.SetResponseMiddlewares(
		Ensure2xxResponseMiddleware,       // 先调用，判断状态是不是请求成功
		resty.AutoParseResponseMiddleware, // 再调用，才能先判断状态码再解析
	)

	return &Client{
		Client: client,
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

// R 返回一个自定义的 Request，以便我们可以调用 SetBodyMap() SetBodyStruct() 解决因 body 顺序不同导致 SHA256 不一样的问题
func (c *Client) R() *Request {
	return &Request{c.Client.R()}
}

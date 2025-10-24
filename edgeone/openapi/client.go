package openapi

import (
	"crypto/tls"
	"fmt"

	"resty.dev/v3"
)

type Client struct {
	*resty.Client
}

// NewClient 创建请求客户端
func NewClient(secret_id string, secret_key string) (*Client, error) {

	if secret_id == "" {
		return nil, fmt.Errorf("check secret_id")
	}
	if secret_key == "" {
		return nil, fmt.Errorf("check secret_key")
	}

	client := resty.New()
	client.SetRequestMiddlewares(
		resty.PrepareRequestMiddleware,              // 先调用，创建 RawRequest
		PreRequestMiddleware(secret_id, secret_key), // 再调用，自定义中间件安全使用 RawRequest
	)
	client.SetResponseMiddlewares(
		Ensure2xxResponseMiddleware,       // 先调用，判断状态是不是请求成功
		resty.AutoParseResponseMiddleware, // 再调用，才能先判断状态码再解析
	)

	return &Client{client}, nil
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
	return &Request{c.Client.R().SetContentType("application/json")}
}

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
func NewClient(endpoint string, accessKeyId string, secretAccessKey string) (*Client, error) {
	if endpoint == "" {
		return nil, fmt.Errorf("check endpoint")
	}
	if _, err := url.Parse(endpoint); err != nil {
		return nil, fmt.Errorf("check endpoint: %w", err)
	}
	if accessKeyId == "" {
		return nil, fmt.Errorf("check accessKeyId")
	}
	if secretAccessKey == "" {
		return nil, fmt.Errorf("check secretAccessKey")
	}

	client := resty.New().SetBaseURL(endpoint)
	client.SetRequestMiddlewares(
		resty.PrepareRequestMiddleware,                     // 先调用，创建 RawRequest
		PreRequestMiddleware(accessKeyId, secretAccessKey), // 再调用，自定义中间件安全使用 RawRequest
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
	return &Request{c.Client.R()}
}

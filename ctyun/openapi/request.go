package openapi

import (
	"fmt"
	"net/url"

	"resty.dev/v3"
)

type Client struct {
	*resty.Client
}

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
	// client.EnableDebug()

	return &Client{client}, nil
}

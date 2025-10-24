package openapi

import (
	"encoding/json"
	"fmt"

	"resty.dev/v3"
)

type Request struct {
	*resty.Request
}

func (r *Request) SetXTCEndpoint(endpoint string) *Request {
	r.SetHeader("X-TC-Endpoint", endpoint)
	return r
}

func (r *Request) SetXTCAction(action string) *Request {
	r.SetHeader("X-TC-Action", action)
	return r
}

func (r *Request) SetXTCVersion(version string) *Request {
	r.SetHeader("X-TC-Version", version)
	return r
}

func (r *Request) SetXTCRegion(region string) *Request {
	r.SetHeader("X-TC-Region", region)
	return r
}

func (r *Request) SetBodyMap(bodyMap map[string]any) *Request {
	bodyBytes, err := json.Marshal(bodyMap)
	if err != nil {
		r.SetError(fmt.Errorf("failed to canonicalize JSON body from struct: %w", err))
		return r
	}
	r.Body = string(bodyBytes)
	return r
}

func (r *Request) SetBodyStruct(bodyStruct any) *Request {
	bodyBytes, err := json.Marshal(bodyStruct)
	if err != nil {
		r.SetError(fmt.Errorf("failed to canonicalize JSON body from struct: %w", err))
		return r
	}
	r.Body = string(bodyBytes)
	return r
}

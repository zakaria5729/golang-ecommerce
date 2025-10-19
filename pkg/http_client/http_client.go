package httpclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	defaultTimeout = 5 * time.Second
	defaultRetry   = 0
)

type LoggerFunc func(format string, v ...any)

type HTTPClient struct {
	client *http.Client
	logger LoggerFunc
	retry  int
}

func New() *HTTPClient {
	return &HTTPClient{
		client: &http.Client{Timeout: defaultTimeout},
		retry:  defaultRetry,
	}
}

func (c *HTTPClient) Timeout(timeout time.Duration) *HTTPClient {
	c.client.Timeout = timeout
	return c
}

func (c *HTTPClient) Retry(retry int) *HTTPClient {
	c.retry = retry
	return c
}

func (c *HTTPClient) Logger(logger LoggerFunc) *HTTPClient {
	c.logger = logger
	return c
}

type Request struct {
	Method      string
	URL         string
	Headers     map[string]string
	QueryParams map[string]string
	RequestBody any
	Response    any
}

func (c *HTTPClient) Do(request Request) error {
	req, err := c.buildRequest(request)
	if err != nil {
		return err
	}

	return c.doRequest(req, request.Response)
}

func (c *HTTPClient) buildRequest(request Request) (*http.Request, error) {
	var bodyReader io.Reader
	if request.RequestBody != nil {
		jsonBytes, err := json.Marshal(request.RequestBody)
		if err != nil {
			if c.logger != nil {
				c.logger("[HTTP ERROR] %v", err)
			}
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewBuffer(jsonBytes)
	}

	method := request.Method
	if method == "" {
		method = http.MethodGet
	}

	requestUrl := request.URL
	if len(request.QueryParams) > 0 {
		params := url.Values{}
		for key, value := range request.QueryParams {
			params.Add(key, value)
		}
		if strings.Contains(requestUrl, "?") {
			requestUrl += "&" + params.Encode()
		} else {
			requestUrl += "?" + params.Encode()
		}
	}

	req, err := http.NewRequest(method, requestUrl, bodyReader)
	if err != nil {
		if c.logger != nil {
			c.logger("[HTTP ERROR] %v", err)
		}
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if request.RequestBody != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	for k, v := range request.Headers {
		req.Header.Set(k, v)
	}

	return req, nil
}

// func (c *HTTPClient) buildRequest(request Request) (*http.Request, error) {
// 	var bodyReader io.Reader
// 	if request.RequestBody != nil {
// 		jsonBytes, err := json.Marshal(request.RequestBody)
// 		if err != nil {
// 			if c.logger != nil {
// 				c.logger("[HTTP ERROR] %v", err)
// 			}
// 			return nil, fmt.Errorf("failed to marshal request body: %w", err)
// 		}
// 		bodyReader = bytes.NewBuffer(jsonBytes)
// 	}

// 	method := request.Method
// 	if method == "" {
// 		method = http.MethodGet
// 	}

// 	req, err := http.NewRequest(method, request.URL, bodyReader)
// 	if err != nil {
// 		if c.logger != nil {
// 			c.logger("[HTTP ERROR] %v", err)
// 		}
// 		return nil, fmt.Errorf("failed to create request: %w", err)
// 	}

// 	if request.RequestBody != nil {
// 		req.Header.Set("Content-Type", "application/json")
// 	}

// 	for k, v := range request.Headers {
// 		req.Header.Set(k, v)
// 	}

// 	return req, nil
// }

func (c *HTTPClient) doRequest(req *http.Request, result any) error {
	var lastErr error

	for attempt := 0; attempt <= c.retry; attempt++ {
		if c.logger != nil {
			c.logger("[HTTP] %s %s (attempt %d)", req.Method, req.URL, attempt+1)
		}

		resp, err := c.client.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("request failed: %w", err)
			if c.logger != nil {
				c.logger("[HTTP ERROR] %v", lastErr)
			}
		} else {
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				if c.logger != nil {
					c.logger("[HTTP ERROR] %v", err)
				}
				return fmt.Errorf("failed to read response body: %w", err)
			}

			if resp.StatusCode < 200 || resp.StatusCode >= 300 {
				lastErr = fmt.Errorf("http error %d: %s", resp.StatusCode, string(body))
				if c.logger != nil {
					c.logger("[HTTP ERROR] %s %s -> %v", req.Method, req.URL, lastErr)
				}
			} else {
				if result != nil {
					if byteResult, ok := result.(*[]byte); ok {
						*byteResult = body
					} else {
						if err := json.Unmarshal(body, result); err != nil {
							if c.logger != nil {
								c.logger("[HTTP ERROR] %v", err)
							}
							return fmt.Errorf("failed to unmarshal response: %w", err)
						}
					}
				}
				return nil
			}
		}

		if attempt < c.retry {
			time.Sleep(500 * time.Millisecond)
		}
	}

	return lastErr
}

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
	defaultRetry   = 0
	defaultTimeout = 5 * time.Second
)

type HttpClient interface {
	Timeout(timeout time.Duration) HttpClient
	Retry(retry int) HttpClient
	Logger(logger LoggerFunc) HttpClient
	Do(request Request) error
}

type LoggerFunc func(format string, v ...any)

type httpClient struct {
	client *http.Client
	logger LoggerFunc
	retry  int
}

func New() HttpClient {
	return &httpClient{
		client: &http.Client{Timeout: defaultTimeout},
		retry:  defaultRetry,
	}
}

func (c *httpClient) Timeout(timeout time.Duration) HttpClient {
	c.client.Timeout = timeout
	return c
}

func (c *httpClient) Retry(retry int) HttpClient {
	c.retry = retry
	return c
}

func (c *httpClient) Logger(logger LoggerFunc) HttpClient {
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

func (c *httpClient) Do(request Request) error {
	req, err := c.buildRequest(request)
	if err != nil {
		return err
	}

	return c.doRequest(req, request.Response)
}

func (c *httpClient) buildRequest(request Request) (*http.Request, error) {
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

func (c *httpClient) doRequest(req *http.Request, result any) error {
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

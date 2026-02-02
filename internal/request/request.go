package request

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"go-zero-template/internal/config"
	"net/http"
	"time"
)

type RequestClient struct {
	client      *http.Client
	services    *config.ServicesConfig
	userBaseURL string
}

func NewRequestClient(services *config.ServicesConfig) *RequestClient {
	return &RequestClient{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		services:    services,
		userBaseURL: fmt.Sprintf("http://%s:%d%s", services.UserService.Host, services.UserService.Port, services.UserService.Path),
	}
}

func (r *RequestClient) Request(method string, url string, body any, headers map[string]string) (any, error) {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("序列化请求体失败: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// 设置自定义 headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("请求失败: status=%d, body=%s", resp.StatusCode, string(respBody))
	}

	var result any
	if len(respBody) > 0 {
		if err := json.Unmarshal(respBody, &result); err != nil {
			return nil, fmt.Errorf("解析响应失败: %w", err)
		}
	}

	return result, nil
}

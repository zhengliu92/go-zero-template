package types

import (
	"encoding/json"
	"errors"
	"fmt"
)

// BaseResponse 通用 API 返回结构
type BaseResponse[T any] struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data T      `json:"data"`
}

// Ok 是否为成功状态（code 为 0 或 200 视为成功）
func (r *BaseResponse[T]) Ok() bool {
	return r != nil && (r.Code == 0 || r.Code == 200)
}

// Err 若为业务失败则返回包含 code、msg 的 error，成功则返回 nil
func (r *BaseResponse[T]) Err() error {
	if r == nil {
		return errors.New("response is nil")
	}
	if r.Ok() {
		return nil
	}
	return &ResponseError{Code: r.Code, Msg: r.Msg}
}

// ResponseError 业务错误，携带 code 与 msg
type ResponseError struct {
	Code int
	Msg  string
}

func (e *ResponseError) Error() string {
	return fmt.Sprintf("code=%d, msg=%s", e.Code, e.Msg)
}

// ParseBaseResponse 从 JSON 字节解析为 BaseResponse[T]
func ParseBaseResponse[T any](data []byte) (*BaseResponse[T], error) {
	var out BaseResponse[T]
	if len(data) == 0 {
		return &out, nil
	}
	if err := json.Unmarshal(data, &out); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}
	return &out, nil
}

// ParseBaseResponseFromAny 从 Request 返回的 any（如 map[string]any）解析为 BaseResponse[T]
func ParseBaseResponseFromAny[T any](data any) (*BaseResponse[T], error) {
	if data == nil {
		return &BaseResponse[T]{}, nil
	}
	bytes, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("序列化响应失败: %w", err)
	}
	return ParseBaseResponse[T](bytes)
}

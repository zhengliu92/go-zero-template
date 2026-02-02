package request

import (
	"fmt"
	"go-zero-template/internal/types"
	"net/http"
)

// CreateUserRequest 创建用户请求体（与 user_service API 一致），复用 types.UserBase 并增加 Password
type CreateUserRequest struct {
	types.UserBase
	Password string `json:"password"`
}

// UpdateUserRequest 更新用户请求体（ID 放在 path），复用 types.UserBase
type UpdateUserRequest struct {
	types.UserBase
}

// GetUserInfoResponse 获取用户信息响应（与 user_service API 一致）
type GetUserInfoResponse struct {
	User types.User `json:"user"`
}

// CreateUserResponse 创建用户响应（与 user_service API 一致）
type CreateUserResponse struct {
	User types.User `json:"user"`
}

// UpdateUserResponse 更新用户响应（与 user_service API 一致）
type UpdateUserResponse struct {
	User types.User `json:"user"`
}

func (r *RequestClient) GetUserInfo(token string) (*GetUserInfoResponse, error) {
	baseURL := r.userBaseURL + "/info"
	raw, err := r.Request(http.MethodGet, baseURL, nil, map[string]string{"Authorization": token})
	if err != nil {
		return nil, err
	}
	resp, err := types.ParseBaseResponseFromAny[GetUserInfoResponse](raw)
	if err != nil {
		return nil, err
	}
	if !resp.Ok() {
		return nil, resp.Err()
	}
	return &resp.Data, nil
}

// CreateUser 调用 user_service 创建用户接口
func (r *RequestClient) CreateUser(token string, req *CreateUserRequest) (*CreateUserResponse, error) {
	baseURL := r.userBaseURL + "/"
	raw, err := r.Request(http.MethodPost, baseURL, req, map[string]string{"Authorization": token})
	if err != nil {
		return nil, err
	}
	resp, err := types.ParseBaseResponseFromAny[CreateUserResponse](raw)
	if err != nil {
		return nil, err
	}
	if !resp.Ok() {
		return nil, resp.Err()
	}
	return &resp.Data, nil
}

// UpdateUser 调用 user_service 更新用户接口
func (r *RequestClient) UpdateUser(token string, id int, req *UpdateUserRequest) (*UpdateUserResponse, error) {
	baseURL := fmt.Sprintf("%s/%d", r.userBaseURL, id)
	raw, err := r.Request(http.MethodPut, baseURL, req, map[string]string{"Authorization": token})
	if err != nil {
		return nil, err
	}
	resp, err := types.ParseBaseResponseFromAny[UpdateUserResponse](raw)
	if err != nil {
		return nil, err
	}
	if !resp.Ok() {
		return nil, resp.Err()
	}
	return &resp.Data, nil
}

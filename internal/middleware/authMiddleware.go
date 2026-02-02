// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package middleware

import (
	"context"
	"go-zero-template/internal/config"
	"go-zero-template/internal/request"
	"go-zero-template/internal/response"
	"go-zero-template/internal/types"
	"net/http"
	"strings"
)

// contextKey 用于 context 的 key 类型
type contextKey string

const (
	// UserContextKey 用户信息在 context 中的 key
	UserContextKey contextKey = "user"
)

type AuthMiddleware struct {
	requestClient *request.RequestClient
}

func NewAuthMiddleware(services *config.ServicesConfig) *AuthMiddleware {
	return &AuthMiddleware{
		requestClient: request.NewRequestClient(services),
	}
}

func (m *AuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")

		// 检查 token 是否存在
		if token == "" {
			response.Response(w, nil, response.NewError(http.StatusUnauthorized, "缺少认证信息"))
			return
		}

		// 检查 token 格式（Bearer token）
		if !strings.HasPrefix(token, "Bearer ") {
			response.Response(w, nil, response.NewError(http.StatusUnauthorized, "无效的认证格式"))
			return
		}

		resp, err := m.requestClient.GetUserInfo(token)
		if err != nil {
			response.Response(w, nil, response.NewError(http.StatusUnauthorized, "认证失败"))
			return
		}

		// 将用户信息存入 context
		user := &resp.User
		ctx := context.WithValue(r.Context(), UserContextKey, user)
		r = r.WithContext(ctx)

		// 传递给下一个 handler
		next(w, r)
	}
}

// GetUserFromContext 从 context 中获取用户信息
func GetUserFromContext(ctx context.Context) (*types.User, bool) {
	user, ok := ctx.Value(UserContextKey).(*types.User)
	return user, ok
}

package response

import "net/http"

// 服务器错误
var InternalServerError = NewError(http.StatusInternalServerError, "服务器错误")
var ParseError = NewError(10005, "解析请求失败")

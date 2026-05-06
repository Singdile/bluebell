package controllers

import "net/http"

// 错误类型定义
// 使用字符串表示错误类型
// 符合RESTful标准,无数字错误码

// 错误类型定义
const (
	// ==================== 参数验证错误 ====================
	// HTTP 400 Bad Request
	ErrValidation  = "validation_error" // JSON格式没问题，参数不符合要求
	ErrInvalidJSON = "invalid_json"     // JSON格式有问题

	// ==================== 认证错误 ====================
	// HTTP 401 Unauthorized

	ErrUnauthorized = "unauthorized"  // 未认证,也即用户未登录
	ErrInvalidToken = "invalid_token" // Token无效
	ErrTokenExpired = "token_expired" // Token过期
	ErrLoginFailed  = "login_failed"  // 登录失败

	// ==================== 权限错误 ====================
	// HTTP 403 Forbidden

	ErrForbidden    = "forbidden"     // 权限不足
	ErrAccessDenied = "access_denied" // 拒绝访问

	// ==================== 资源错误 ====================
	// HTTP 404 Not Found

	ErrNotFound     = "not_found"      // 资源不存在
	ErrUserNotExist = "user_not_exist" // 用户不存在
	ErrPostNotExist = "post_not_exist" // 帖子不存在
	ErrIdNotExist   = "id_not_exist"   //id 不存在

	// ==================== 冲突错误 ====================
	// HTTP 409 Conflict

	ErrConflict    = "conflict"        // 资源冲突
	ErrUserExists  = "username_exists" // 用户名已存在
	ErrEmailExists = "email_exists"    // 邮箱已注册

	// ==================== 限流错误 ====================
	// HTTP 429 Too Many Requests

	ErrRateLimit = "rate_limit" // 请求频率超限

	// ==================== 服务器错误 ====================
	// HTTP 500 Internal Server Error

	ErrInternal   = "internal_error"    // 服务器内部错误
	ErrDatabase   = "database_error"    // 数据库错误
	ErrThirdParty = "third_party_error" // 第三方服务错误
)

// ==================== HTTP状态码映射 ====================

// getHTTPStatus 根据错误类型返回对应的HTTP状态码
// 这是RESTful风格的核心：HTTP状态码有明确语义
func getHTTPStatus(errorType string) int {
	switch errorType {
	// 参数验证错误 -> 400
	case ErrValidation, ErrInvalidJSON:
		return http.StatusBadRequest

	// 认证错误 -> 401
	case ErrUnauthorized, ErrInvalidToken, ErrTokenExpired, ErrLoginFailed:
		return http.StatusUnauthorized

	// 权限错误 -> 403
	case ErrForbidden, ErrAccessDenied:
		return http.StatusForbidden

	// 资源不存在 -> 404
	case ErrNotFound, ErrUserNotExist, ErrPostNotExist, ErrIdNotExist:
		return http.StatusNotFound

	// 资源冲突 -> 409
	case ErrConflict, ErrUserExists, ErrEmailExists:
		return http.StatusConflict

	// 限流 -> 429
	case ErrRateLimit:
		return http.StatusTooManyRequests

	// 服务器错误 -> 500
	case ErrInternal, ErrDatabase, ErrThirdParty:
		return http.StatusInternalServerError

	// 默认：服务器错误
	default:
		return http.StatusInternalServerError
	}
}

// ==================== 错误消息映射 ====================

// errorMsgMap 错误类型的默认消息,用于给用户提供提示信息
// 可用于快速调用，无需每次都传入message
var errorMsgMap = map[string]string{
	ErrValidation:   "参数验证失败",
	ErrInvalidJSON:  "请求参数格式错误,请检查输入",
	ErrUnauthorized: "未认证，请先登录",
	ErrInvalidToken: "认证信息无效,请重新登录",
	ErrTokenExpired: "登录已过期,请重新登录",
	ErrLoginFailed:  "用户名或密码错误",
	ErrForbidden:    "权限不足",
	ErrAccessDenied: "拒绝访问",
	ErrNotFound:     "资源不存在",
	ErrUserNotExist: "用户不存在",
	ErrPostNotExist: "帖子不存在",
	ErrIdNotExist:   "id不存在",
	ErrConflict:     "请求冲突,请稍后重试",
	ErrUserExists:   "用户名已存在",
	ErrEmailExists:  "邮箱已注册",
	ErrRateLimit:    "请求频率超限，请稍后再试",
	ErrInternal:     "服务器内部错误",
	ErrDatabase:     "服务器内部错误",
	ErrThirdParty:   "服务器内部错误",
}

// getErrorMsg 获取错误类型的默认消息
func getErrorMsg(errorType string) string {
	msg, ok := errorMsgMap[errorType]
	if !ok {
		return "服务器错误"
	}
	return msg
}

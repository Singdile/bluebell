package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 响应结构体
type SuccessResponse struct {
	Data    interface{} //业务数据
	Message string      //提示信息
}

// 单一错误响应结构
type ErrorResponse struct {
	Error   string //错误类型标识
	Message string //用户友好的错误信息
	Field   string //错误字段
}

// Success 成功响应, 比如Get,Put等成功
func Success(ctx *gin.Context, data any) {
	ctx.JSON(http.StatusOK, SuccessResponse{
		Data: data,
	})
}

// Success 成功响应, 比如Get,Put等成功
func SuccessWithMessage(ctx *gin.Context, message string, data any) {
	ctx.JSON(http.StatusOK, SuccessResponse{
		Data:    data,
		Message: message,
	})
}

// Fail 错误响应
func Fail(ctx *gin.Context, errorType string, message string) {
	//获取HTTP状态码
	httpstatus := getHTTPStatus(errorType)

	ctx.JSON(httpstatus, ErrorResponse{
		Error:   errorType,
		Message: message,
	})
}

// Fail 错误响应，使用默认的响应信息
func FailWithDefault(ctx *gin.Context, errorType string) {
	//获取HTTP状态码
	httpstatus := getHTTPStatus(errorType)

	//返回响应
	ctx.JSON(httpstatus, ErrorResponse{
		Error:   errorType,
		Message: getErrorMsg(errorType),
	})
}

// FailWithField 错误响应(带错误字段信息)
func FailWithField(ctx *gin.Context, errorType string, message string, field string) {
	//获取HTTP状态码
	httpstatus := getHTTPStatus(errorType)

	ctx.JSON(httpstatus, ErrorResponse{
		Error:   errorType,
		Message: message,
		Field:   field,
	})
}

package controllers

import (
	"bluebell/models/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ========== 成功响应辅助函数 ==========

// Success 成功响应
func Success(ctx *gin.Context, data any) {
	ctx.JSON(http.StatusOK, dto.SuccessResponse{
		Data: data,
	})
}

// SuccessWithMessage 成功响应（带消息）
func SuccessWithMessage(ctx *gin.Context, message string, data any) {
	ctx.JSON(http.StatusOK, dto.SuccessResponse{
		BaseResponse: dto.BaseResponse{Message: message},
		Data:         data,
	})
}

// ========== 错误响应辅助函数 ==========

// Fail 错误响应
func Fail(ctx *gin.Context, errorType string, message string) {
	httpstatus := getHTTPStatus(errorType)

	ctx.JSON(httpstatus, dto.ErrorResponse{
		BaseResponse: dto.BaseResponse{Message: message},
		Error:        errorType,
	})
}

// FailWithDefault 错误响应（使用默认消息）
func FailWithDefault(ctx *gin.Context, errorType string) {
	httpstatus := getHTTPStatus(errorType)

	ctx.JSON(httpstatus, dto.ErrorResponse{
		BaseResponse: dto.BaseResponse{Message: getErrorMsg(errorType)},
		Error:        errorType,
	})
}

// FailWithField 错误响应（带字段信息）
func FailWithField(ctx *gin.Context, errorType string, message string, field string) {
	httpstatus := getHTTPStatus(errorType)

	ctx.JSON(httpstatus, dto.ErrorResponse{
		BaseResponse: dto.BaseResponse{Message: message},
		Error:        errorType,
		Field:        field,
	})
}

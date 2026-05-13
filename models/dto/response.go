package dto

import "bluebell/models"

// ========== 通用响应 ==========

// BaseResponse 基础响应
type BaseResponse struct {
	Message string `json:"message"`
}

// ErrorResponse 错误响应
type ErrorResponse struct {
	BaseResponse
	Error string `json:"error"`
	Field string `json:"field,omitempty"`
}

// ========== 认证响应 ==========

// TokenResponse 登录成功响应
type TokenResponse struct {
	BaseResponse
	Data string `json:"data"` // JWT Token
}

// ========== 社区响应 ==========

// CommunityListResponse 社区列表响应
type CommunityListResponse struct {
	BaseResponse
	Data []models.CommunityDetail `json:"data"`
}

// CommunityDetailResponse 社区详情响应
type CommunityDetailResponse struct {
	BaseResponse
	Data models.CommunityDetail `json:"data"`
}

// ========== 帖子响应 ==========

// PostListResponse 帖子列表响应
type PostListResponse struct {
	BaseResponse
	Data []models.PostListItem `json:"data"`
}

// PostDetailResponse 帖子详情响应
type PostDetailResponse struct {
	BaseResponse
	Data models.PostDetail `json:"data"`
}

// ========== 通用成功响应 ==========

// SuccessResponse 通用成功响应（用于创建、投票等）
type SuccessResponse struct {
	BaseResponse
	Data interface{} `json:"data"`
}

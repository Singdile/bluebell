package dto

// ========== 认证请求 ==========

// LoginRequest 登录请求（只暴露必要字段）
type LoginRequest struct {
	Username string `json:"username" binding:"required" example:"admin"`
	Password string `json:"password" binding:"required" example:"password123"`
}

// SignUpRequest 注册请求
type SignUpRequest struct {
	Username   string `json:"username" binding:"required" example:"admin"`
	Password   string `json:"password" binding:"required" example:"password123"`
	RePassword string `json:"re_password" binding:"required" example:"password123"`
}

// ========== 帖子请求 ==========

// CreatePostRequest 创建帖子请求
type CreatePostRequest struct {
	Title       string `json:"title" binding:"required" example:"帖子标题"`
	Content     string `json:"content" binding:"required" example:"帖子内容"`
	CommunityID int64  `json:"community_id" binding:"required" example:"1"`
}

// VoteRequest 投票请求
type VoteRequest struct {
	PostID int64  `json:"post_id" binding:"required" example:"1"`
	Action string `json:"action" binding:"required,oneof=up down cancel" example:"up"`
}

// PostListRequest 帖子列表查询请求
type PostListRequest struct {
	Page        int64  `json:"page" form:"page" example:"1"`
	Pagesize    int64  `json:"pagesize" form:"pagesize" example:"10"`
	Order       string `json:"order" form:"order" example:"time"` // time 或 score
	CommunityID int64  `json:"community_id" form:"community_id" example:"1"`
}
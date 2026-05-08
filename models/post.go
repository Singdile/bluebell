package models

import "time"

// Post 表示帖子的结构体
type Post struct {
	Content     string    `json:"content" db:"content" binding:"required"`
	Title       string    `json:"title" db:"title" binding:"required"`
	ID          int64     `json:"id,string" db:"post_id"`
	AuthorID    int64     `json:"author_id,string" db:"author_id"`
	CommunityID int64     `json:"community_id,string" db:"community_id" binding:"required"`
	Status      int32     `json:"status" db:"status"`
	CreateTime  time.Time `json:"create_time" db:"create_time"`
}

// PostDetail 表示帖子的详情信息，有关联的作者名称，社区
type PostDetail struct {
	AuthorName string
	Post       *Post
	Community  *CommunityDetail
}

// PostListItem 帖子列表项，用于分页展示
// 只包含列表页所需要的的字段，数据库层面负责截取摘要
type PostListItem struct {
	PostID         int64     `json:"post_id,string" db:"post_id"`
	Title          string    `json:"title" db:"title"`
	ContentPreview string    `json:"content_preview" db:"content_preview"` // 内容摘要
	AuthorID       int64     `json:"author_id,string" db:"author_id"`
	AuthorName     string    `json:"author_name" db:"author_name"`
	CommunityID    int64     `json:"community_id,string" db:"community_id"`
	CommunityName  string    `json:"community_name" db:"community_name"`
	Status         int32     `json:"status" db:"status"`
	CreateTime     time.Time `json:"create_time" db:"create_time"`
}

// 用于post分页响应
type PostListResponse struct {
	Total      int64           `json:"total,string"`       //总记录数
	Page       int64           `json:"page,string"`        //当前页码
	PageSize   int64           `json:"page_size,string"`   //每页数量
	TotalPages int64           `json:"total_pages,string"` //总页数
	List       []*PostListItem `json:"list"`               //帖子列表
}

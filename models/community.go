package models

import "time"

// Community 表示社区的结构体
type Community struct {
	ID           int64     `json:"id,string" db:"id"`
	Name         string    `json:"name" db:"community_name"`
	Introduction string    `json:"introduction" db:"introduction"`
	CreateTime   time.Time `json:"create_time" db:"create_time"`
	UpdateTime   time.Time `json:"update_time" db:"update_time"`
	CreatorID    int64     `json:"creator_id" db:"creator_id"`
	Status       int       `json:"status" db:"status"`
}

type CommunityDetail struct {
	ID           int64  `json:"id,string" db:"id"`
	Name         string `json:"name" db:"community_name"`
	Introduction string `json:"introduction" db:"introduction"`
}

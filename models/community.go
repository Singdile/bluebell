package models

type Community struct {
	ID   int    `json:"id" db:"community_id"`
	Name string `json:"name" db:"community_name"`
}

type CommunityDetail struct {
	ID           int    `json:"id" db:"community_id"`
	Name         string `json:"name" db:"community_name"`
	Introduction string `json:"introduction" db:"introduction"`
}

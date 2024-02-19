package model

import "time"

type UserListRequest struct {
	Page         int
	PageSize     int
	SearchName   string
	SearchEmail  string
	FilterRoleID []string
}

type RoleListResponse struct {
	Name    string    `json:"name"`
	AddedAt time.Time `json:"added_at"`
}

type UserListResponse struct {
	ID       string             `json:"id"`
	Name     string             `json:"name"`
	Email    string             `json:"email"`
	Roles    []RoleListResponse `json:"roles"`
	JoinedAt time.Time          `json:"joined_at"`
}

package model

import "time"

type UserListRequest struct {
	Page         int
	PageSize     int
	SearchName   string
	SearchEmail  string
	FilterRoleID []string
	RoleIDs      string
}

type UserGetByIDRequest struct {
	ID string `validate:"required"`
}

type UserUpdateRolesRequest struct {
	UserID  string   `validate:"required"`
	RoleIDs []string `json:"role_ids" validate:"required"`
}

type UserDeleteRequest struct {
	ID string `validate:"required"`
}

type RoleResponse struct {
	ID          string    `json:"id"`
	DisplayName string    `json:"display_name"`
	AddedAt     time.Time `json:"added_at"`
}

type UserResponse struct {
	ID       string         `json:"id"`
	Name     string         `json:"name"`
	Email    string         `json:"email"`
	Roles    []RoleResponse `json:"roles"`
	JoinedAt time.Time      `json:"joined_at"`
}

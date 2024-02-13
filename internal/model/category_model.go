package model

import "time"

type CategoryListRequest struct {
	Page          int
	PageSize      int
	SortByPopular bool
	UserID        string
}

type CategoryCreateRequest struct {
	Name   string `json:"name" validate:"required"`
	UserID string `json:"-" validate:"required"`
}

type CategoryUpdateRequest struct {
	ID     string `json:"-" validate:"required"`
	Name   string `json:"name" validate:"required"`
	UserID string `json:"-" validate:"required"`
}

type CategoryResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	CreatedBy   string    `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	MemberCount uint64    `json:"member_count"`
}

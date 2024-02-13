package model

import "time"

type CategoryListRequest struct {
	Page          int
	PageSize      int
	SortByPopular bool
	UserID        string
}

type CategoryResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	CreatedBy   string    `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	MemberCount uint64    `json:"member_count"`
}

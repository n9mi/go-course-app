package model

import "bytes"

type CourseListRequest struct {
	Page               int
	PageSize           int
	CategoryID         string
	SortByMinimumPrice bool
	SortByMaximumPrice bool
	IsFree             bool
	SearchTitle        string
	UserID             string
}

type CourseListResponse struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Category    string  `json:"category"`
	PriceIdr    float64 `json:"price_idr"`
	BannerLink  string  `json:"banner_link"`
	CreatedBy   string  `json:"created_by"`
	MemberCount uint64  `json:"member_count"`
}

type CourseGetRequest struct {
	ID     string `validate:"requried"`
	UserID string
}

type CourseCreateRequest struct {
	Name        string        `json:"name" validate:"required"`
	Description string        `json:"description" validate:"required"`
	CategoryID  string        `json:"category_id" validate:"required"`
	PriceIdr    float64       `json:"price_idr" validate:"required"`
	CreatedBy   string        `json:"created_by" validate:"required"`
	Image       *bytes.Buffer `json:"-"`
}

type CourseUpdateRequest struct {
	ID            string        `json:"-" validate:"required"`
	Name          string        `json:"name" validate:"required"`
	Description   string        `json:"description" validate:"required"`
	CategoryID    string        `json:"category_id" validate:"required"`
	PriceIdr      float64       `json:"price_idr" validate:"required"`
	UserID        string        `json:"-" validate:"required"`
	Image         *bytes.Buffer `json:"-"`
	IsRemoveImage bool          `json:"is_remove_image" validate:"required"`
}

type CourseDeleteRequest struct {
	ID     string `json:"-" validate:"required"`
	UserID string `json:"-" validate:"required"`
}

type CourseResponse struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Category    string  `json:"category"`
	PriceIdr    float64 `json:"price_idr"`
	BannerLink  string  `json:"banner_link"`
	CreatedBy   string  `json:"created_by"`
	MemberCount uint64  `json:"member_count"`
}

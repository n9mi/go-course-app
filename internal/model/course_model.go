package model

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

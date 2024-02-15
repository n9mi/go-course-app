package model

type MessageResponse struct {
	Code     int      `json:"-"`
	Messages []string `json:"messages"`
}

type DataResponse[T any] struct {
	Data T `json:"data"`
}

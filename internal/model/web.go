package model

type WebResponse[T any] struct {
	Code     int      `json:"-"`
	Messages []string `json:"messages,omitempty"`
	Data     T        `json:"data,omitempty"`
}

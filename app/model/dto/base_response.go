package dto

type Response[T any] struct {
	IsSuccessful bool   `json:"is_successful"`
	Message      string `json:"message"`
	Status       string `json:"status"`
	Data         T      `json:"data"`
}

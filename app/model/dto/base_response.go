package dto

type Response[T any] struct {
	IsSuccessful bool   `json:"IsSuccessful"`
	Message      string `json:"Message"`
	Status       string `json:"Status"`
	Data         T      `json:"Data"`
}

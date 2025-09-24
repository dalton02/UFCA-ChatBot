package util_dto

import "time"

type QueryPaginationDto struct {
	Page  int `form:"page" json:"page" validate:"gte=1" example:"1"`
	Limit int `form:"limit" json:"limit" validate:"gte=1,lte=100" example:"10"`
}
type ResponsePaginatedDto[T any] struct {
	QueryPaginationDto
	Total int `json:"total"`
	Data  T   `json:"data"`
}

type TimeStampDefaultDB struct {
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type AppResponse struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message"`

	StatusCode int `json:"statusCode"`
}

package example_dto

type RouteTestValidation struct {
	Name string `json:"name" validator:"required"`
}

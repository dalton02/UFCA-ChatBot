package app

import "fmt"

type AppError struct {
	Message    string `json:"message"`
	StatusCode int    `json:"statusCode"`
}

func (e AppError) Error() string {

	return fmt.Sprintf("%s (status: %d)", e.Message, e.StatusCode)

}

func NotFound(message string) AppError {
	return AppError{
		Message:    message,
		StatusCode: 404,
	}
}

func BadRequest(message string) AppError {
	return AppError{
		Message:    message,
		StatusCode: 400,
	}
}

func Conflict(message string) AppError {
	return AppError{
		Message:    message,
		StatusCode: 409,
	}
}

func Forbidden(message string) AppError {
	return AppError{
		Message:    message,
		StatusCode: 401,
	}
}

func InternalServerError(message string) AppError {
	return AppError{
		Message:    message,
		StatusCode: 500,
	}
}

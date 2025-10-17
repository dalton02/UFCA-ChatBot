package interceptor

import (
	app "licor_model/core/util/errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response é a estrutura padrão para todas as respostas da API
type Response[T any] struct {
	Message    string `json:"message"`
	Data       T      `json:"data"`
	StatusCode int    `json:"statusCode"`
}

// AppSuccess retorna uma resposta de sucesso (200 OK)
func AppSuccess(ctx *gin.Context, message string, data any) {
	ctx.JSON(http.StatusOK, Response[any]{
		Message:    message,
		Data:       data,
		StatusCode: http.StatusOK,
	})
}

// AppCreated retorna uma resposta para criação bem-sucedida (201 Created)
func AppCreated(ctx *gin.Context, message string, data any) {
	ctx.JSON(http.StatusCreated, Response[any]{
		Message:    message,
		Data:       data,
		StatusCode: http.StatusCreated,
	})
}

// AppNoContent retorna uma resposta sem conteúdo (204 No Content)
func AppNoContent(ctx *gin.Context) {
	ctx.AbortWithStatusJSON(http.StatusNoContent, Response[any]{
		Message:    "",
		Data:       nil,
		StatusCode: http.StatusNoContent,
	})
}

// AppBadRequest retorna uma resposta de requisição inválida (400 Bad Request)
func AppBadRequest(ctx *gin.Context, message string) {
	ctx.AbortWithStatusJSON(http.StatusBadRequest, Response[any]{
		Message:    message,
		Data:       nil,
		StatusCode: http.StatusBadRequest,
	})
}

// AppUnauthorized retorna uma resposta de não autorizado (401 Unauthorized)
func AppUnauthorized(ctx *gin.Context, message string) {
	ctx.AbortWithStatusJSON(http.StatusUnauthorized, Response[any]{
		Message:    message,
		Data:       nil,
		StatusCode: http.StatusUnauthorized,
	})
}

// AppForbidden retorna uma resposta de acesso proibido (403 Forbidden)
func AppForbidden(ctx *gin.Context, message string) {
	ctx.AbortWithStatusJSON(http.StatusForbidden, Response[any]{
		Message:    message,
		Data:       nil,
		StatusCode: http.StatusForbidden,
	})
}

// AppNotFound retorna uma resposta de recurso não encontrado (404 Not Found)
func AppNotFound(ctx *gin.Context, message string) {
	ctx.AbortWithStatusJSON(http.StatusNotFound, Response[any]{
		Message:    message,
		Data:       nil,
		StatusCode: http.StatusNotFound,
	})
}

// AppNotFound retorna uma resposta de recurso não encontrado (404 Not Found)
func AppConflict(ctx *gin.Context, message string) {
	ctx.AbortWithStatusJSON(http.StatusConflict, Response[any]{
		Message:    message,
		Data:       nil,
		StatusCode: http.StatusConflict,
	})
}

// AppInternalServerError retorna uma resposta de erro interno (500 Internal Server Error)
func AppInternalServerError(ctx *gin.Context, message string) {

	ctx.AbortWithStatusJSON(http.StatusInternalServerError, Response[any]{
		Message:    message,
		Data:       nil,
		StatusCode: http.StatusInternalServerError,
	})
}

// AppCustomResponse retorna uma resposta personalizada com qualquer status code
func AppCustomResponse(ctx *gin.Context, statusCode int, message string, data any) {
	ctx.JSON(statusCode, Response[any]{
		Message:    message,
		Data:       data,
		StatusCode: statusCode,
	})
}
func AppError(ctx *gin.Context, err error) {

	errApp, valid := err.(app.AppError)

	if valid {
		ctx.AbortWithStatusJSON(errApp.StatusCode, Response[any]{
			Message:    errApp.Message,
			Data:       nil,
			StatusCode: errApp.StatusCode,
		})

		return
	} //Motivação: Caso eu retorne um erro não formatado, certeza que é burrice do codigo mesmo, e não da regra de negocio
	ctx.AbortWithStatusJSON(500, Response[any]{
		Message:    err.Error(),
		Data:       nil,
		StatusCode: 500,
	})

}

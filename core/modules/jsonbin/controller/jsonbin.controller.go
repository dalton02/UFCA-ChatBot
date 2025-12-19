package jsonbin_controller

import (
	auth_dto "licor_model/core/modules/auth/dto"
	auth_middleware "licor_model/core/modules/auth/middleware"
	jsonbin_dto "licor_model/core/modules/jsonbin/dto"
	jsonbin_service "licor_model/core/modules/jsonbin/service"
	"licor_model/core/util/interceptor"

	"github.com/gin-gonic/gin"
)

type controller struct {
	service        *jsonbin_service.JsonBinService
	authMiddleware *auth_middleware.AuthMiddleware
}

func New(service *jsonbin_service.JsonBinService, authMid *auth_middleware.AuthMiddleware) *controller {
	return &controller{
		service:        service,
		authMiddleware: authMid,
	}
}

func (controller *controller) SaveContext(ctx *gin.Context) {

	type Data struct {
		Data []jsonbin_dto.JsonFile `json:"data"`
	}

	var data Data

	err := interceptor.ValidateAndExtract(ctx, &data)
	if err != nil {
		interceptor.AppError(ctx, err)
		return
	}

	userInfo := ctx.MustGet("userInfo").(auth_dto.UserDto)

	err = controller.service.SaveContext(userInfo.Email, data.Data)
	if err != nil {
		interceptor.AppError(ctx, err)
		return
	}

	interceptor.AppSuccess(ctx, "Contexto salvo com sucesso", nil)

}

func (controller *controller) FeedIA(ctx *gin.Context) {

	type Data struct {
		Year     int    `json:"year"`
		Month    int    `json:"month"`
		Day      int    `json:"day"`
		Filename string `json:"filename"`
	}

	var data Data

	err := interceptor.ValidateAndExtract(ctx, &data)
	if err != nil {
		interceptor.AppError(ctx, err)
		return
	}

	err = controller.service.FeedFromFile(data.Year, data.Month, data.Day, data.Filename)
	if err != nil {
		interceptor.AppError(ctx, err)
		return
	}

	interceptor.AppSuccess(ctx, "Contexto salvo com sucesso", nil)

}

func (controller *controller) GetFromFilePath(ctx *gin.Context) {

	type Data struct {
		Year     int    `json:"year"`
		Month    int    `json:"month"`
		Day      int    `json:"day"`
		Filename string `json:"filename"`
	}

	var data Data

	err := interceptor.ValidateAndExtract(ctx, &data)
	if err != nil {
		interceptor.AppError(ctx, err)
		return
	}
	jsonInfo, err := controller.service.GetFromPath(data.Year, data.Month, data.Day, data.Filename)
	if err != nil {
		interceptor.AppError(ctx, err)
		return
	}

	interceptor.AppSuccess(ctx, "Encontrado com sucesso", jsonInfo)
}

func (controller *controller) GetLatestUserData(ctx *gin.Context) {
	userInfo := ctx.MustGet("userInfo").(auth_dto.UserDto)

	data, err := controller.service.GetLatestUserData(userInfo.Email)
	if err != nil {
		interceptor.AppError(ctx, err)
		return
	}

	interceptor.AppSuccess(ctx, "Último contexto do usuário", data)
}

func (controller *controller) GetTree(ctx *gin.Context) {

	data, err := controller.service.GetTree()
	if err != nil {
		interceptor.AppError(ctx, err)
		return
	}

	interceptor.AppSuccess(ctx, "ARVORE GERADA", data)
}

func (controller *controller) GetLatestData(ctx *gin.Context) {
	data, err := controller.service.GetLatestData()
	if err != nil {
		interceptor.AppError(ctx, err)
		return
	}

	interceptor.AppSuccess(ctx, "Último contexto geral", data)
}

func (c *controller) Routes(g *gin.RouterGroup) {
	group := g.Group("/admin", c.authMiddleware.HelperAccess)
	{
		group.GET("/find/latest", c.GetLatestData)
		group.GET("/find/my-latest", c.GetLatestUserData)
		group.POST("/find/filename", c.GetFromFilePath)
		group.GET("/tree", c.GetTree)
		group.POST("/save-context", c.SaveContext)

		restrict := group.Group("/restrict", c.authMiddleware.ManagerAccess)
		{
			restrict.POST("/feed-ia", c.FeedIA)

		}

	}
}

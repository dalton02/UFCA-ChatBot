package auth_controller

import (
	auth_dto "licor_model/core/modules/auth/dto"
	auth_service "licor_model/core/modules/auth/service"
	"licor_model/core/util/interceptor"

	"github.com/gin-gonic/gin"
)

// refiz o controller pq agora a gnt ta usando o gin

type AuthController struct {
	authService *auth_service.AuthService
}

func NewAuthController(authService *auth_service.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

// @Summary		Login
// @Description	Autentica um usuário com email e senha
// @Tags			Auth
// @Accept			json
// @Produce		json
// @Param			body	body		auth_dto.LoginRequestDto			true	"Credenciais de login"
// @Success		200		{object}	util_dto.AppResponse{data=auth_dto.AuthResponseDto}	"Retorna um token JWT"
// @Router			/auth/login [post]
func (c *AuthController) Login(ctx *gin.Context) {
	var loginDto auth_dto.LoginRequestDto

	if err := interceptor.ValidateAndExtract(ctx, &loginDto); err != nil {
		interceptor.AppBadRequest(ctx, err.Error())
		return
	}

	response, err := c.authService.Login(loginDto)
	if err != nil {
		interceptor.AppUnauthorized(ctx, err.Error())
		return
	}

	interceptor.AppSuccess(ctx, "Login com sucesso", response)
}

// @Summary		Cadastro de usuário
// @Description	Autentica um usuário com email e senha
// @Tags			Auth
// @Accept			json
// @Produce		json
// @Param			body	body		auth_dto.RegisterRequestDto			true	"Credenciais de login"
// @Success		200		{object}	util_dto.AppResponse{data=auth_dto.AuthResponseDto}	"Retorna um token JWT"
// @Router			/auth/register [post]
func (c *AuthController) Register(ctx *gin.Context) {
	var registerDto auth_dto.RegisterRequestDto

	if err := interceptor.ValidateAndExtract(ctx, &registerDto); err != nil {
		interceptor.AppBadRequest(ctx, err.Error())
		return
	}

	response, err := c.authService.Register(registerDto)
	if err != nil {
		interceptor.AppBadRequest(ctx, err.Error())
		return
	}

	interceptor.AppCreated(ctx, "Usuário cadastrado com sucesso", response)
}

func (c *AuthController) Routes(g *gin.RouterGroup) {
	authGroup := g.Group("/auth")
	{
		authGroup.POST("/login", c.Login)
		authGroup.POST("/register", c.Register)
	}
}

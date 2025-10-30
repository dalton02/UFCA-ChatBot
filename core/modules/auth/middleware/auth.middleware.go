package auth_middleware

import (
	"fmt"
	auth_service "licor_model/core/modules/auth/service"
	"licor_model/core/util/interceptor"
	guard_util "licor_model/core/util/jwt"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	authService *auth_service.AuthService
}

func NewAuthMiddleware(authService *auth_service.AuthService) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
	}
}

func (c *AuthMiddleware) JwtGuard(ctx *gin.Context) {

	bearerToken := ctx.Request.Header.Get("Authorization")

	token, err := guard_util.GetJwtInfo(bearerToken)

	if err != nil {
		interceptor.AppUnauthorized(ctx, "Token invalido")
	}

	userID := token["id"].(string)
	_, err = c.authService.GetUserByID(userID)
	if err != nil {
		interceptor.AppUnauthorized(ctx, "Token valido, porém usuário não existe")
	}

	ctx.Set("userID", userID)
}

func (c *AuthMiddleware) IpValidation(ctx *gin.Context) {
	ip := ctx.ClientIP()

	fmt.Println("Tentativa de acesso por: ", ip)
	// _, err := strconv.Atoi(os.Getenv("LIMIT_TRY_ACCESS"))
	// if err != nil {
	// 	limit = 10
	// }

	//Adicionar  serviço redis pra bloquear os caba depois
	// tryAcessCount, err := c.authService.GetRedisService.GetIpTryAccess(ip)
	// if err == nil && tryAcessCount > limit {
	// 	interceptor.AppForbidden(ctx, "seu ip foi bloqueado por exceder o limite de tentativas de acesso")
	// 	return
	// }
	// c.authService.GetRedisService.AddIpTryAccess(ip)

}

package auth_middleware

import (
	"fmt"
	auth_dto "licor_model/core/modules/auth/dto"
	auth_service "licor_model/core/modules/auth/service"
	redis_service "licor_model/core/modules/redis"
	"licor_model/core/util/interceptor"
	guard_util "licor_model/core/util/jwt"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	authService  *auth_service.AuthService
	redisService *redis_service.RedisService
}

func NewAuthMiddleware(authService *auth_service.AuthService, redisService *redis_service.RedisService) *AuthMiddleware {
	return &AuthMiddleware{
		authService:  authService,
		redisService: redisService,
	}
}

func (c *AuthMiddleware) JwtGuard(ctx *gin.Context) {

	bearerToken := ctx.Request.Header.Get("Authorization")

	token, err := guard_util.GetJwtInfo(bearerToken)

	if err != nil {
		interceptor.AppUnauthorized(ctx, "Token invalido")
		return
	}

	userID, exist := token["id"].(string)
	if !exist {
		interceptor.AppUnauthorized(ctx, "Token invalido")
		return
	}
	user, err := c.authService.GetUserByID(userID)
	if err != nil {
		interceptor.AppUnauthorized(ctx, "Token valido, porém usuário não existe")
		return
	}

	ctx.Set("userID", userID)
	ctx.Set("userInfo", user)
}

func (c *AuthMiddleware) HelperAccess(ctx *gin.Context) {

	user := ctx.MustGet("userInfo").(auth_dto.UserDto)

	if user.Role == auth_dto.Helper || user.Role == auth_dto.Manager {
		return
	}

	interceptor.AppForbidden(ctx, "Tu pode é nada aqui")

}

func (c *AuthMiddleware) ManagerAccess(ctx *gin.Context) {

	user := ctx.MustGet("userInfo").(auth_dto.UserDto)

	if user.Role == auth_dto.Manager {
		return
	}

	interceptor.AppForbidden(ctx, "Você não tem acesso para realizar esta operação")

}
func (c *AuthMiddleware) IpValidation(ctx *gin.Context) {
	ip := ctx.ClientIP()

	fmt.Println("Tentativa de acesso por: ", ip)
	limit, err := strconv.Atoi(os.Getenv("LIMIT_TRY_ACCESS"))
	if err != nil {
		limit = 100
	}

	tryAcessCount, err := c.redisService.GetIpTryAccess(ip)
	if err == nil && tryAcessCount > limit {
		interceptor.AppForbidden(ctx, "seu ip foi bloqueado por exceder o limite de tentativas de acesso")
		return
	}
	c.redisService.AddIpTryAccess(ip)

}

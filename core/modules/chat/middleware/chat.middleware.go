package chat_middleware

import (
	"fmt"
	redis_service "licor_model/core/modules/redis"
	"licor_model/core/util/interceptor"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ChatMiddleware struct {
	redisService *redis_service.RedisService
}

func NewChatService(redisService *redis_service.RedisService) *ChatMiddleware {
	return &ChatMiddleware{
		redisService: redisService,
	}
}

func (middleware *ChatMiddleware) ValidateAIDailyUse(ctx *gin.Context) {
	userID := ctx.MustGet("userID").(string)

	count, err := middleware.redisService.GetAICountDaily(userID)

	dailyLimit, _ := strconv.Atoi(os.Getenv("DAILY_AI_LIMIT"))

	fmt.Print(dailyLimit, count)

	if err == nil && count > dailyLimit {
		interceptor.AppUnauthorized(ctx, "Limite diário do chatbot utilizado")
		return
	}

	middleware.redisService.AddAIUserAccess(userID)

}

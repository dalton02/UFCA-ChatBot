package redis_service

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

const AI_COUNT_ACCESS = "ai_count_access:"
const IP_ACCESS_KEY = "ip_try_access:"

type RedisService struct {
	client *redis.Client
}

func NewRedisService() (*RedisService, error) {

	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	fmt.Println("Conectando ao redis em:", client.Options().Addr)

	_, err := client.Ping().Result()
	if err != nil {
		fmt.Println("Failed to connect to Redis: ", err)
	}
	return &RedisService{
		client: client,
	}, err
}

func (s *RedisService) AddAIUserAccess(userID string) {

	pipe := s.client.TxPipeline()

	key := AI_COUNT_ACCESS + userID

	pipe.Incr(key)
	pipe.Expire(key, time.Hour*24)

	pipe.Exec()

}

func (s *RedisService) AddIpTryAccess(ip string) {

	pipe := s.client.TxPipeline()

	key := IP_ACCESS_KEY + ip

	pipe.Incr(key)
	pipe.Expire(key, time.Hour*24)

	pipe.Exec()

}

func (s *RedisService) GetAICountDaily(userID string) (value int, err error) {
	key := AI_COUNT_ACCESS + userID

	get := s.client.Get(key)

	if get == nil {
		return 0, fmt.Errorf("valor não encontrado")
	}

	value, err = strconv.Atoi(get.Val())
	if err != nil {
		return 0, fmt.Errorf("valor não encontrado")
	}
	return value, err
}

func (s *RedisService) GetIpTryAccess(ip string) (value int, err error) {
	key := IP_ACCESS_KEY + ip
	get := s.client.Get(key)

	if get == nil {
		return 0, fmt.Errorf("valor não encontrado")
	}

	value, err = strconv.Atoi(get.Val())

	if err != nil {
		return 0, fmt.Errorf("valor não encontrado")
	}
	return value, err
}

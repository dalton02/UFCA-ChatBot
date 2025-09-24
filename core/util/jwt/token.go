package guard_util

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

func GenerateJwt[T any](data T, minutes int) (string, error) {
	secretKey := os.Getenv("JWT_SECRET")
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["data"] = data
	claims["exp"] = time.Now().Add(time.Duration(minutes) * time.Minute).Unix()
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func GetJwtInfo(tokenString string) (map[string]interface{}, error) {
	secretKey := os.Getenv("JWT_SECRET")
	claims := jwt.MapClaims{}

	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return make(map[string]interface{}), err
	}

	data, _ := claims["data"].(map[string]interface{})
	return data, nil
}

package utils

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
)

func GenerateJWT(id uint, role string, name string, username string) string {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = id
	claims["role"] = role
	claims["name"] = name
	claims["username"] = username
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	t, _ := token.SignedString(viper.GetString("secret"))
	return t
}

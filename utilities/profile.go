package utilities

import (
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

func GetUserProfile(c echo.Context) (id uint, role string, userName string) {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userName = claims["username"].(string)
	role = claims["role"].(string)
	id = claims["id"].(uint)
	return
}

func ValidateDateOnly(date string) (time.Time, error) {
	return time.Parse("2006-01-02", date)
}

func BindAndValidate(c echo.Context, i interface{}) error {
	err := c.Bind(&i)
	if err != nil {
		log.Println("here1")
		log.Println(err)
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}
	log.Println("here2")
	if err = c.Validate(i); err != nil {
		log.Println("here3")
		return err
	}
	return nil
}

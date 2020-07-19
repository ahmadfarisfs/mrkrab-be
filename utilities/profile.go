package utilities

import (
	"log"
	"net/http"

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

func BindAndValidate(c echo.Context, i interface{}) error {
	err := c.Bind(&i)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}
	if err = c.Validate(i); err != nil {
		return err
	}
	return nil
}

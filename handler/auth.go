package handler

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ahmadfarisfs/krab-core/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func (h *Handler) Authenticate(c echo.Context) error {
	// user := c.Get("user").(*jwt.MapClaims)
	user := c.Get("user").(*jwt.Token)
	claims := *user.Claims.(*jwt.MapClaims)
	// name := claims["name"].(string)
	fmt.Println(claims)
	timestamp := claims["exp"]
	var expireOffset = 3600
	if validity, ok := timestamp.(float64); ok {
		tm := time.Unix(int64(validity), 0)
		remainder := tm.Sub(time.Now())

		if remainder > 0 {
			return c.JSON(http.StatusOK, int(remainder.Seconds()+float64(expireOffset)))
		}
	}
	return c.JSON(http.StatusUnprocessableEntity, errors.New("Token Expired"))

}
func (h *Handler) Test(c echo.Context) error {
	return c.JSON(http.StatusOK, "Hello from mrkrab")
}
func (h *Handler) Login(c echo.Context) error {
	req := &loginRequest{}
	if err := req.bind(c); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}
	account, _, err := h.userStore.ListUser(utils.CommonRequest{
		Filter: map[string]interface{}{
			"username": req.Username,
			// "password": string(hashedPassword),
		},
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	if len(account) != 1 {
		return c.JSON(http.StatusUnauthorized, "Username not found")
	}
	// hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 8)
	// if err != nil {
	// 	return c.JSON(http.StatusInternalServerError, err.Error())
	// }
	err = bcrypt.CompareHashAndPassword([]byte(account[0].Password), []byte(req.Password))
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusUnauthorized, "Wrong password")
	}

	// if len(account) != 1 {
	// 	return c.JSON(http.StatusUnauthorized, "Wrong username/password combination")
	// }
	var accountID int
	if account[0].Account == nil {
		accountID = 0
	} else {
		accountID = int(account[0].Account.ID)
	}
	token := utils.GenerateJWT(account[0].BaseModel.ID, account[0].Role, account[0].Fullname, account[0].Username, accountID)
	return c.JSON(http.StatusOK, token)
}

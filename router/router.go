package router

import (
	"fmt"
	"strings"

	"github.com/ahmadfarisfs/krab-core/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

func New() *echo.Echo {
	e := echo.New()
	e.Logger.SetLevel(log.DEBUG)
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Logger())
	e.Use(middleware.CORS())
	e.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		Claims:     &jwt.MapClaims{},
		SigningKey: []byte("!!SECRET!!"),
		Skipper: func(c echo.Context) bool {
			// if strings.HasPrefix(c.Request().Host, "localhost") {
			// 	return true
			// }
			fmt.Println("JWT: " + c.Path())
			if strings.HasSuffix(c.Path(), "/auth/login") || strings.HasSuffix(c.Path(), "/auth/test") {
				return true
			}
			return false
		},
	}))
	/*	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:  []string{"*"},
		AllowMethods:  []string{"GET", "PUT", "POST", "DELETE", "OPTIONS"},
		AllowHeaders:  []string{"Origin", "Authorization", "Content-Type", "Bearer", "X-CSRF-Token"},
		ExposeHeaders: []string{"x-total-count", "Content-Range"},
		MaxAge:        50,

		AllowCredentials: true,
		//	ValidateHeaders:  false,
		//	AllowOrigins: []string{"*"},
		//	ExposeHeaders: []string{
		//		"Access-Control-Expose-Headers",
		//	},
		//		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		//	AllowHeaders: []string{"*"},

		//	AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))*/
	e.Use(utils.ParseCommonMiddleware)
	e.Validator = NewValidator()
	return e
}

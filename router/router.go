package router

import (
	"github.com/ahmadfarisfs/krab-core/utils"
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

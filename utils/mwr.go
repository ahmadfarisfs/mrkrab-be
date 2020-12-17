package utils

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func ParseCommonMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		payload, err := ParseCommonRequest(c)
		if err != nil {
			echo.NewHTTPError(http.StatusUnprocessableEntity, err)
		}
		c.Set("parsedQuery", payload)
		c.Response().Header().Set("Access-Control-Expose-Headers", "Content-Range")

		return next(c)
	}
}

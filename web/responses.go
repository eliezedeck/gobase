package web

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func OK(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"ok": true,
	})
}

func BadRequestError(c echo.Context, message string) error {
	return c.JSON(http.StatusBadRequest, map[string]interface{}{
		"error": message,
	})
}

func BadRequestErrorWithCode(c echo.Context, errCode int, message string) error {
	return c.JSON(http.StatusBadRequest, map[string]interface{}{
		"error": message,
		"code":  errCode,
	})
}

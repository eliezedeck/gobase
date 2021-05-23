package http

import (
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

func MustGetUrlParamInt32(c echo.Context, name string, defaultvalue int32) int32 {
	params := c.Request().URL.Query()
	if v := params.Get(name); v == "" {
		return defaultvalue
	} else {
		if i, err := strconv.ParseInt(v, 10, 32); err != nil {
			return defaultvalue
		} else {
			return int32(i)
		}
	}
}

func MustGetUrlParamStringTrimmed(c echo.Context, name string, defaultvalue string) string {
	params := c.Request().URL.Query()
	if v := strings.TrimSpace(params.Get(name)); v == "" {
		return defaultvalue
	} else {
		return v
	}
}

func MustGetUrlParamBoolean(c echo.Context, name string, defaultvalue bool) bool {
	params := c.Request().URL.Query()
	if v := strings.TrimSpace(params.Get(name)); v == "" {
		return defaultvalue
	} else {
		normalized := strings.TrimSpace(strings.ToLower(v))
		return normalized == "true"
	}
}

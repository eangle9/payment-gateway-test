package foundation

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"gorm.io/gorm/utils"
)

func InitCORS() echo.MiddlewareFunc {
	origins := viper.GetStringSlice("cors.origin")
	if len(origins) == 0 {
		origins = []string{"*"}
	}
	allowCredentials := viper.GetString("cors.allow_credentials")
	if allowCredentials == "" {
		allowCredentials = "true"
	}
	headers := viper.GetStringSlice("cors.headers")
	if len(headers) == 0 {
		headers = []string{"*"}
	}
	methods := viper.GetStringSlice("cors.methods")
	if len(methods) == 0 {
		methods = []string{
			"GET",
			"POST",
			"PUT",
			"DELETE",
			"PATCH",
			"OPTIONS",
		}
	}
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			requestHeader := c.Request().Header.Get("Origin")
			if utils.Contains(origins, requestHeader) {
				c.Response().Header().Set("Access-Control-Allow-Origin", requestHeader)
			} else {
				c.Response().Header().Set("Access-Control-Allow-Origin", origins[0])
			}
			c.Response().Header().Set("Access-Control-Allow-Credentials", allowCredentials)
			c.Response().Header().Set("Access-Control-Allow-Headers", strings.Join(headers, ","))
			c.Response().Header().Set("Access-Control-Allow-Methods", strings.Join(methods, ","))
			if c.Request().Method == "OPTIONS" {
				return c.NoContent(http.StatusNoContent)
			}
			return next(c)
		}
	}
}

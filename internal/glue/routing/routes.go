package routing

import (
	"github.com/labstack/echo/v4"
)

type Router struct {
	Method      string
	Path        string
	Handler     echo.HandlerFunc
	Middlewares []echo.MiddlewareFunc
}

func RegisterRoute(
	grp *echo.Group,
	routes []Router,
) {
	for _, route := range routes {
		var handler []echo.MiddlewareFunc
		handler = append(handler, route.Middlewares...)
		grp.Add(route.Method, route.Path, route.Handler, handler...)
	}
}

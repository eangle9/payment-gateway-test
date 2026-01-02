package routing

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func TestRoute(
	grp *echo.Group,
	//authMidd middleware.AuthMiddleware,
) {
	router := []Router{
		{
			Method: "GET",
			Path:   "/test",
			Handler: func(c echo.Context) error {
				return c.JSON(http.StatusOK, map[string]string{
					"message": "Hello, World!",
				})
			},
			Middlewares: []echo.MiddlewareFunc{},
		},
	}

	RegisterRoute(grp, router)
	fmt.Println("Route registered successfully")
}

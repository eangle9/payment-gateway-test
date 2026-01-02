package company

import (
	"net/http"
	"pg/internal/glue/routing"
	"pg/internal/handler/middleware"
	"pg/internal/handler/rest"

	"github.com/labstack/echo/v4"
)

func Route(
	grp *echo.Group,
	authMiddle middleware.AuthMiddleware,
	handler rest.Company,
) {
	router := []routing.Router{
		{
			Method:      http.MethodPost,
			Path:        "/signup-company-owner",
			Handler:     handler.RegisterCompany,
			Middlewares: []echo.MiddlewareFunc{},
		},
		{
			Method:      http.MethodPost,
			Path:        "/login",
			Handler:     handler.Login,
			Middlewares: []echo.MiddlewareFunc{},
		},
		{
			Method:  http.MethodPost,
			Path:    "/generate-secret-token",
			Handler: handler.GenerateSecretToken,
			Middlewares: []echo.MiddlewareFunc{
				authMiddle.AuthenticateUser(),
			},
		},
	}

	routing.RegisterRoute(grp, router)
}

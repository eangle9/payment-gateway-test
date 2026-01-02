package paymentintent

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
	handler rest.PaymentIntent,
) {
	router := []routing.Router{
		{
			Method:  http.MethodPost,
			Path:    "/payment-intents",
			Handler: handler.InitPaymentIntent,
			Middlewares: []echo.MiddlewareFunc{
				authMiddle.AuthenticateAdminUser(),
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/payment-intents/:id",
			Handler: handler.GetPaymentIntentDetail,
			Middlewares: []echo.MiddlewareFunc{
				authMiddle.AuthenticateAdminUser(),
			},
		},
	}

	routing.RegisterRoute(grp, router)
}

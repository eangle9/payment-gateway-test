package rest

import (
	"github.com/labstack/echo/v4"
)

type Company interface {
	RegisterCompany(c echo.Context) error
	Login(c echo.Context) error
	GenerateSecretToken(c echo.Context) error
}

type PaymentIntent interface {
	InitPaymentIntent(c echo.Context) error
	GetPaymentIntentDetail(c echo.Context) error
}

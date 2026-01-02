package initiator

import (
	"pg/docs"
	"pg/initiator/platform"
	"pg/internal/glue/routing"
	"pg/internal/glue/routing/company"
	paymentintent "pg/internal/glue/routing/payment_intent"
	"pg/internal/handler/middleware"
	"pg/platform/hcrypto"
	"pg/platform/hlog"

	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func InitRouter(
	group *echo.Group,
	handler HandlerLayer,
	log hlog.Logger,
	tokenMaket hcrypto.Maker,
	storage PersistenceLayer,
	platform platform.Layer,
) {
	md := middleware.InitAuthMiddleware(
		log.Named("auth-middleware"),
		tokenMaket,
		storage.company)
	docs.SwaggerInfo.Schemes = viper.GetStringSlice("swagger.schemes")
	docs.SwaggerInfo.Host = viper.GetString("swagger.host")
	docs.SwaggerInfo.BasePath = "/api"
	group.GET("/swagger/*", echoSwagger.WrapHandler)
	routing.TestRoute(group)
	company.Route(group, md, handler.company)
	paymentintent.Route(group, md, handler.paymentIntent)
}

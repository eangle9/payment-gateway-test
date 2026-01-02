package initiator

import (
	"pg/initiator/platform"
	"pg/internal/module"
	"pg/internal/module/company"
	paymentintent "pg/internal/module/payment_intent"
	"pg/platform/hlog"

	"github.com/spf13/viper"
)

type ModuleLayer struct {
	Company       module.Company
	PaymentIntent module.PaymentIntent
}

func InitModule(pl PersistenceLayer, log hlog.Logger,
	platform platform.Layer) ModuleLayer {
	return ModuleLayer{
		Company: company.New(
			pl.company,
			log.Named("company-module"),
			platform.Token),
		PaymentIntent: paymentintent.New(
			pl.paymentIntent,
			log.Named("payment-intent-module"),
			pl.company,
			platform.HTTPClient,
			platform.AMQP,
		),
	}
}

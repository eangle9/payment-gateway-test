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
	company       module.Company
	paymentIntent module.PaymentIntent
}

func InitModule(pl PersistenceLayer, log hlog.Logger,
	platform platform.Layer) ModuleLayer {
	return ModuleLayer{
		company: company.New(
			pl.company,
			log.Named("company-module"),
			platform.Token),
		paymentIntent: paymentintent.New(
			pl.paymentIntent,
			log.Named("payment-intent-module"),
			pl.company,
			platform.HTTPClient,
			viper.GetString("CHECKOUT_BASE_URL")+viper.GetString("ONETIME_CHECKOUT_PAGE"),
		),
	}
}

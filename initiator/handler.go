package initiator

import (
	"pg/internal/handler/rest"
	"pg/internal/handler/rest/company"
	paymentintent "pg/internal/handler/rest/payment_intent"
	"pg/platform/hlog"
	"time"
)

type HandlerLayer struct {
	company       rest.Company
	paymentIntent rest.PaymentIntent
}

func InitHandler(ml ModuleLayer, log hlog.Logger,
	timeout time.Duration) HandlerLayer {
	return HandlerLayer{
		company: company.New(
			log.Named("company-handler"),
			ml.Company,
			timeout),
		paymentIntent: paymentintent.New(
			log.Named("payment-intent-handler"),
			ml.PaymentIntent,
			timeout,
		),
	}
}

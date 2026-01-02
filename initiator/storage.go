package initiator

import (
	persistencedb "pg/internal/constant/persistenceDB"
	"pg/internal/storage"
	"pg/internal/storage/company"
	paymentintent "pg/internal/storage/payment_intent"
	"pg/platform/hlog"
)

type PersistenceLayer struct {
	db            persistencedb.PersistenceDB
	company       storage.Company
	paymentIntent storage.PaymentIntent
}

func InitPersistence(db persistencedb.PersistenceDB, log hlog.Logger) PersistenceLayer {
	return PersistenceLayer{
		db:            db,
		company:       company.NewCompanycePersistance(db, log.Named("company-persistence")),
		paymentIntent: paymentintent.NewPaymentIntentPersistance(db, log.Named("payment-intent-persistence")),
	}
}

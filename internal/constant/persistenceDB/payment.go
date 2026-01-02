package persistencedb

import (
	"context"
	"encoding/json"
	"pg/initiator/platform/amqp"
	"pg/internal/constant"
	"pg/internal/constant/errors"
	"pg/internal/constant/model/db"
	"pg/internal/constant/model/dto"
	"pg/platform/sql"

	"go.uber.org/zap"
)

func (q *PersistenceDB) CreatePaymentIntentTx(ctx context.Context,
	param dto.CreatePaymentIntent, client amqp.Client) (*db.PaymentIntent, error) {
	tx, err := q.pool.Begin(ctx)
	if err != nil {
		return nil, errors.ErrUnableToCreate.Wrap(err, "error starting transaction")
	}
	defer func() {
		if p := recover(); p != nil {
			if rbErr := tx.Rollback(ctx); rbErr != nil {
				q.log.Error(ctx, "error rolling back transaction", zap.Error(rbErr))
			}
			panic(p) // re-throw panic after Rollback
		} else if err != nil {
			if rbErr := tx.Rollback(ctx); rbErr != nil {
				q.log.Error(ctx, "error rolling back transaction", zap.Error(rbErr))
			}
		}
	}()

	tQ := q.WithTx(tx)
	customer, err := tQ.CreateCustomer(ctx,
		db.CreateCustomerParams{
			CompanyID:   param.CompanyID,
			FullName:    sql.StringOrNull(param.Customer.FullName),
			PhoneNumber: param.Customer.PhoneNumber,
			Email:       sql.StringOrNull(param.Customer.Email),
		})
	if err != nil {
		return nil, err
	}

	extra, err := json.Marshal(param.Extra)
	if err != nil {
		return nil, err
	}

	paymentIntent, err := tQ.CreatePaymentIntent(ctx,
		db.CreatePaymentIntentParams{
			CompanyID:   param.CompanyID,
			PaymentType: constant.PaymentTypeOnetime,
			Amount:      param.Amount,
			Status:      string(constant.Pending),
			Currency:    string(param.Currency),
			CallbackUrl: param.CallBackURL,
			ReturnUrl:   param.ReturnURL,
			Description: sql.StringOrNull(param.Description),
			CustomerID:  customer.ID,
			Extra:       sql.MapJSONOrNull(extra),
			BillRefNo:   sql.StringOrNull(param.BillRefNO),
		})
	if err != nil {
		return nil, err
	}

	// Publish to RabbitMQ
	payload, err := json.Marshal(map[string]string{
		"payment_intent_id": paymentIntent.ID.String(),
	})
	if err != nil {
		return nil, err
	}

	if err := client.Publish(ctx, "", "payment_processing", payload); err != nil {
		return nil, err
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, errors.ErrUnableToCreate.Wrap(err, "error committing transaction")
	}

	return &paymentIntent, nil
}

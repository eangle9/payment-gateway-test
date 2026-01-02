package paymentintent

import (
	"context"
	"encoding/json"
	"pg/initiator/platform/amqp"
	"pg/internal/constant"
	"pg/internal/constant/errors"
	"pg/internal/constant/model/db"
	"pg/internal/constant/model/dto"
	persistencedb "pg/internal/constant/persistenceDB"
	"pg/internal/storage"
	"pg/platform/hlog"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type paymentIntentPersistance struct {
	persistenceQueries persistencedb.PersistenceDB
	logger             hlog.Logger
}

func NewPaymentIntentPersistance(persistenceQueries persistencedb.PersistenceDB,
	logger hlog.Logger) storage.PaymentIntent {
	return &paymentIntentPersistance{
		persistenceQueries: persistenceQueries,
		logger:             logger,
	}
}

func (p *paymentIntentPersistance) CreatePaymentIntent(ctx context.Context,
	param dto.CreatePaymentIntent, client amqp.Client) (*dto.PaymentIntent, error) {
	pi, err := p.persistenceQueries.CreatePaymentIntentTx(ctx, param, client)
	if err != nil {
		err = errors.ErrUnableToCreate.Wrap(err, "unable to create payment intent")
		p.logger.Error(ctx,
			"unable to create payment intent",
			zap.Error(err), zap.String("company-id", param.CompanyID.String()))
		return nil, err
	}

	extraMap := make(map[string]any)
	if pi.Extra.Bytes != nil {
		if err := json.Unmarshal(pi.Extra.Bytes, &extraMap); err != nil {
			err = errors.ErrBadRequest.Wrap(err, "unable to unmarshal extra fields")
			p.logger.Error(ctx,
				"error unmarshalling extra fields",
				zap.Error(err), zap.String("extra", string(pi.Extra.Bytes)))
			return nil, err
		}
	}

	return &dto.PaymentIntent{
		ID:          pi.ID,
		CompanyID:   pi.CompanyID,
		CustomerID:  pi.CustomerID,
		PaymentType: constant.PaymentType(pi.PaymentType),
		Amount:      pi.Amount,
		Status:      constant.Status(pi.Status),
		Currency:    constant.Currency(pi.Currency),
		CallBackURL: pi.CallbackUrl,
		ReturnURL:   pi.ReturnUrl,
		Extra:       extraMap,
		BillRefNO:   pi.BillRefNo.String,
		ExpireAt:    pi.ExpireAt.Time,
		CreatedAt:   pi.CreatedAt,
		UpdatedAt:   pi.UpdatedAt,
	}, nil
}

func (p *paymentIntentPersistance) GetPaymentIntentByID(ctx context.Context,
	id uuid.UUID) (*dto.PaymentIntentDetail, error) {
	pi, err := p.persistenceQueries.GetPaymentIntentByID(ctx, id)
	if err != nil {
		err = errors.ErrUnableToGet.Wrap(err, "unable to get payment intent by id")
		p.logger.Error(ctx, "unable to get payment intent by id",
			zap.Error(err), zap.String("payment-intent-id", id.String()))
		return nil, err
	}

	extra := make(map[string]any)
	if err := json.Unmarshal(pi.Extra.Bytes, &extra); err != nil {
		err = errors.ErrBadRequest.Wrap(err, "unable to unmarshal extra fields")
		p.logger.Error(ctx, "unable to unmarshal extra fields",
			zap.Error(err), zap.String("extra", string(pi.Extra.Bytes)))
		return nil, err
	}
	customer := dto.Customer{}
	if err := json.Unmarshal(pi.Customer.Bytes, &customer); err != nil {
		err = errors.ErrBadRequest.Wrap(err, "unable to unmarshal customer data")
		p.logger.Error(ctx, "unable to unmarshal customer data",
			zap.Error(err), zap.String("customer", string(pi.Customer.Bytes)))
		return nil, err
	}
	company := dto.Company{}
	if err := json.Unmarshal(pi.Company.Bytes, &company); err != nil {
		err = errors.ErrBadRequest.Wrap(err, "unable to unmarshal company data")
		p.logger.Error(ctx, "unable to unmarshal company data",
			zap.Error(err), zap.String("company", string(pi.Company.Bytes)))
		return nil, err
	}

	return &dto.PaymentIntentDetail{
		ID:          pi.ID,
		PaymentType: constant.PaymentType(pi.PaymentType),
		Amount:      pi.Amount,
		Status:      constant.Status(pi.Status),
		Currency:    constant.Currency(pi.Currency),
		CallBackURL: pi.CallbackUrl,
		ReturnURL:   pi.ReturnUrl,
		Description: pi.Description.String,
		Extra:       extra,
		BillRefNO:   pi.BillRefNo.String,
		ExpireAt:    pi.ExpireAt.Time,
		CreatedAt:   pi.CreatedAt,
		UpdatedAt:   pi.UpdatedAt,
		Customer:    customer,
		Company:     company,
	}, nil
}

func (p *paymentIntentPersistance) UpdatePaymentIntentStatus(ctx context.Context,
	id uuid.UUID, status string) error {
	_, err := p.persistenceQueries.UpdatePaymentIntentStatus(ctx, db.UpdatePaymentIntentStatusParams{
		ID:     id,
		Status: status,
	})
	if err != nil {
		err = errors.ErrUnableToUpdate.Wrap(err, "unable to update payment intent status")
		p.logger.Error(ctx, "unable to update payment intent status",
			zap.Error(err), zap.String("payment-intent-id", id.String()))
		return err
	}

	return nil
}
func (p *paymentIntentPersistance) GetPaymentIntentByIDForUpdate(ctx context.Context,
	id uuid.UUID) (*dto.PaymentIntent, error) {
	pi, err := p.persistenceQueries.GetPaymentIntentByIDForUpdate(ctx, id)
	if err != nil {
		err = errors.ErrUnableToGet.Wrap(err, "unable to get payment intent for update")
		p.logger.Error(ctx, "unable to get payment intent for update",
			zap.Error(err), zap.String("payment-intent-id", id.String()))
		return nil, err
	}

	extraMap := make(map[string]any)
	if pi.Extra.Bytes != nil {
		if err := json.Unmarshal(pi.Extra.Bytes, &extraMap); err != nil {
			err = errors.ErrBadRequest.Wrap(err, "unable to unmarshal extra fields")
			p.logger.Error(ctx,
				"error unmarshalling extra fields",
				zap.Error(err), zap.String("extra", string(pi.Extra.Bytes)))
			return nil, err
		}
	}

	return &dto.PaymentIntent{
		ID:          pi.ID,
		CompanyID:   pi.CompanyID,
		CustomerID:  pi.CustomerID,
		PaymentType: constant.PaymentType(pi.PaymentType),
		Amount:      pi.Amount,
		Status:      constant.Status(pi.Status),
		Currency:    constant.Currency(pi.Currency),
		CallBackURL: pi.CallbackUrl,
		ReturnURL:   pi.ReturnUrl,
		Extra:       extraMap,
		BillRefNO:   pi.BillRefNo.String,
		ExpireAt:    pi.ExpireAt.Time,
		CreatedAt:   pi.CreatedAt,
		UpdatedAt:   pi.UpdatedAt,
	}, nil
}

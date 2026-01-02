package paymentintent

import (
	"context"
	"encoding/json"
	"pg/internal/constant"
	"pg/internal/constant/errors"
	"pg/internal/constant/model/db"
	"pg/internal/constant/model/dto"
	persistencedb "pg/internal/constant/persistenceDB"
	"pg/internal/storage"
	"pg/platform/hlog"
	"pg/platform/sql"
	"pg/platform/utils"

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
	param dto.CreatePaymentIntent) (*dto.PaymentIntent, error) {
	extra, err := json.Marshal(param.Extra)
	if err != nil {
		err := errors.ErrInternalServerError.Wrap(err, "unable to marshal extra fields")
		p.logger.Error(ctx, "error marshalling extra fields", zap.Error(err))
		return nil, err
	}
	pi, err := p.persistenceQueries.CreatePaymentIntent(ctx, db.CreatePaymentIntentParams{
		CompanyID:   param.CompanyID,
		CustomerID:  param.CustomerID,
		PaymentType: string(param.PaymentType),
		Amount:      param.Amount,
		Currency:    string(param.Currency),
		CallbackUrl: param.CallBackURL,
		ReturnUrl:   param.ReturnURL,
		Description: sql.StringOrNull(param.Description),
		Extra:       utils.MapJSONOrNull(extra),
		Status:      string(param.Status),
		BillRefNo:   sql.StringOrNull(param.BillRefNO),
	})
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

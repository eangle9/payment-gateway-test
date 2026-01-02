package paymentintent

import (
	"context"
	"pg/internal/constant"
	"pg/internal/constant/errors"
	"pg/internal/constant/model/dto"
	"pg/internal/module"
	"pg/internal/storage"
	"pg/platform/hlog"
	"pg/platform/httpclient"
	"pg/platform/utils"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type paymentIntent struct {
	log                  hlog.Logger
	paymentIntentStorage storage.PaymentIntent
	companyStorage       storage.Company
	httpClient           httpclient.HTTPClient
	topayURL             string
}

func New(paymentIntentStorage storage.PaymentIntent,
	log hlog.Logger,
	companyStorage storage.Company,
	httpClient httpclient.HTTPClient,
	topayURL string) module.PaymentIntent {
	return &paymentIntent{
		log:                  log,
		paymentIntentStorage: paymentIntentStorage,
		companyStorage:       companyStorage,
		httpClient:           httpClient,
		topayURL:             topayURL,
	}
}

func (p *paymentIntent) InitPaymentIntent(ctx context.Context,
	param dto.InitPaymentIntent, companyID string) (*dto.PaymentIntent, error) {
	if err := param.Validate(); err != nil {
		err = errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		p.log.Warn(ctx, "invalid input", zap.Error(err))
		return nil, err
	}

	cID, err := uuid.Parse(companyID)
	if err != nil {
		err = errors.ErrInvalidUserInput.Wrap(err, "unable to parse company id")
		p.log.Error(ctx, "error parsing company id",
			zap.Error(err), zap.String("company-id", companyID))
		return nil, err
	}
	company, err := p.companyStorage.GetCompanyByID(ctx, cID)
	if err != nil {
		return nil, err
	}
	customer, err := p.companyStorage.CreateCustomer(ctx, dto.CreateCustomer{
		CompanyID:   company.ID,
		FullName:    param.Customer.FullName,
		PhoneNumber: param.Customer.PhoneNumber,
		Email:       param.Customer.Email,
	})
	if err != nil {
		return nil, err
	}
	if param.CallBackURL == "" {
		param.CallBackURL = company.CallBackURL
	}
	if param.ReturnURL == "" {
		param.CallBackURL = company.ReturnURL
	}
	param.Currency = strings.ToUpper(param.Currency)
	billRefNO, err := utils.GenerateHash(uuid.NewString()+time.Now().String()+utils.CapitalLetters, 8)
	if err != nil {
		err = errors.ErrInternalServerError.Wrap(err, "unable to generate hash")
		p.log.Error(ctx, "unable to generate hash", zap.Error(err))
		return nil, err
	}

	paymentIntent, err := p.paymentIntentStorage.CreatePaymentIntent(ctx, dto.CreatePaymentIntent{
		CompanyID:   company.ID,
		PaymentType: constant.PaymentTypeOnetime,
		Amount:      param.Amount,
		Status:      constant.Pending,
		Currency:    constant.Currency(param.Currency),
		CallBackURL: param.CallBackURL,
		ReturnURL:   param.ReturnURL,
		Description: param.Description,
		CustomerID:  customer.ID,
		Extra:       param.Extra,
		BillRefNO:   billRefNO,
	})
	if err != nil {
		return nil, err
	}

	return paymentIntent, nil
}

func (p *paymentIntent) GetPaymentIntentDetail(ctx context.Context,
	id string) (*dto.PaymentIntentDetail, error) {
	pID, err := uuid.Parse(id)
	if err != nil {
		err = errors.ErrInternalServerError.Wrap(err, "unable to parse payment intent id")
		p.log.Error(ctx, "error parsing payment intent id",
			zap.Error(err), zap.String("payment-intent-id", id))
		return nil, err
	}

	return p.paymentIntentStorage.GetPaymentIntentByID(ctx, pID)
}

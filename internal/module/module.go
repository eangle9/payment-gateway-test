package module

import (
	"context"
	"pg/internal/constant/model/dto"
)

type Company interface {
	RegisterCompany(ctx context.Context,
		param dto.CreateCompany) (*dto.Company, error)
	Login(ctx context.Context, arg dto.LoginRequest) (*dto.SignInResponse, error)
	GenerateToken(ctx context.Context,
		userID string) (*dto.CompanyCredentialResponse, error)
}

type PaymentIntent interface {
	InitPaymentIntent(ctx context.Context,
		param dto.InitPaymentIntent, companyID string) (*dto.PaymentIntent, error)
	GetPaymentIntentDetail(ctx context.Context,
		id string) (*dto.PaymentIntentDetail, error)
}

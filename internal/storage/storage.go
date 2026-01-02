package storage

import (
	"context"
	"pg/initiator/platform/amqp"
	"pg/internal/constant/model/dto"

	"github.com/google/uuid"
)

type Company interface {
	AddCompany(ctx context.Context,
		arg dto.CreateCompany) (*dto.Company, error)
	GetCompanyByID(ctx context.Context,
		id uuid.UUID) (*dto.Company, error)
	CreateCompanyToken(ctx context.Context,
		arg dto.CreateCompanyToken) (*dto.CompanyToken, error)
	InactiveToken(ctx context.Context,
		companyID uuid.UUID) error
	CreateCustomer(ctx context.Context,
		arg dto.CreateCustomer) (*dto.Customer, error)
	GenerateCompanyCredentials(ctx context.Context,
		arg dto.CreateCompanyToken) error
	CreateUser(ctx context.Context,
		param dto.CreateUser) (*dto.User, error)
	GetUserByID(ctx context.Context,
		id uuid.UUID) (*dto.User, error)
	GetUserByPhoneOrEmail(ctx context.Context,
		phone string) (*dto.User, error)
	CreateUserToken(ctx context.Context,
		param dto.UserToken) (*dto.UserToken, error)
	GetActiveUserTokenByUserID(ctx context.Context,
		id uuid.UUID) (*dto.UserToken, error)
	ResetActiveToken(ctx context.Context, id uuid.UUID) error
	GetActiveCompanyTokenByID(ctx context.Context, id uuid.UUID) (*dto.CompanyToken, error)
}

type PaymentIntent interface {
	CreatePaymentIntent(ctx context.Context,
		param dto.CreatePaymentIntent, client amqp.Client) (*dto.PaymentIntent, error)
	GetPaymentIntentByID(ctx context.Context,
		id uuid.UUID) (*dto.PaymentIntentDetail, error)
	UpdatePaymentIntentStatus(ctx context.Context,
		id uuid.UUID, status string) error
}

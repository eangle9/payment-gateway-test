package company

import (
	"context"
	"pg/internal/constant/errors"
	"pg/internal/constant/errors/sqlcerr"
	"pg/internal/constant/model/db"
	"pg/internal/constant/model/dto"
	persistencedb "pg/internal/constant/persistenceDB"
	"pg/internal/storage"
	"pg/platform/hlog"
	"pg/platform/sql"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type companyPersistance struct {
	persistenceQueries persistencedb.PersistenceDB
	logger             hlog.Logger
}

func NewCompanycePersistance(persistenceQueries persistencedb.PersistenceDB,
	logger hlog.Logger) storage.Company {
	return &companyPersistance{
		persistenceQueries: persistenceQueries,
		logger:             logger,
	}
}

func (c *companyPersistance) AddCompany(ctx context.Context,
	arg dto.CreateCompany) (*dto.Company, error) {
	company, err := c.persistenceQueries.CreateCompany(ctx,
		db.CreateCompanyParams{
			Name:               arg.Name,
			RegistrationNumber: arg.RegistrationNumber,
			AddressStreet:      sql.StringOrNull(arg.AddressStreet),
			AddressCity:        sql.StringOrNull(arg.AddressCity),
			AddressState:       sql.StringOrNull(arg.AddressState),
			AddressPostalCode:  sql.StringOrNull(arg.AddressPostalCode),
			AddressCountry:     sql.StringOrNull(arg.AddressCountry),
			PrimaryPhone:       sql.StringOrNull(arg.PrimaryPhone),
			SecondaryPhone:     sql.StringOrNull(arg.SecondaryPhone),
			Email:              sql.StringOrNull(arg.Email),
			Website:            sql.StringOrNull(arg.Website),
			CallbackUrl:        sql.StringOrNull(arg.CallBackURL),
			ReturnUrl:          sql.StringOrNull(arg.ReturnURL),
		})
	if err != nil {
		err := errors.ErrUnableToCreate.Wrap(err, "Unable to add Company")
		c.logger.Error(ctx, "Unable to add Company",
			zap.Error(err), zap.String("name", arg.Name))
		return nil, err
	}

	return &dto.Company{
		ID:                 company.ID,
		Name:               company.Name,
		RegistrationNumber: company.RegistrationNumber,
		AddressStreet:      company.AddressStreet.String,
		AddressCity:        company.AddressCity.String,
		AddressState:       company.AddressState.String,
		AddressPostalCode:  company.AddressPostalCode.String,
		AddressCountry:     company.AddressCountry.String,
		PrimaryPhone:       company.PrimaryPhone.String,
		SecondaryPhone:     company.SecondaryPhone.String,
		Email:              company.Email.String,
		Website:            company.Website.String,
		CallBackURL:        company.CallbackUrl.String,
		ReturnURL:          company.ReturnUrl.String,
		CreatedAt:          company.CreatedAt,
		UpdatedAt:          company.UpdatedAt,
	}, nil
}
func (c *companyPersistance) GetCompanyByID(ctx context.Context, id uuid.UUID) (*dto.Company, error) {
	company, err := c.persistenceQueries.GetCompanyByID(ctx, id)
	if err != nil {
		err := errors.ErrUnableToGet.Wrap(err, "Unable to get Company")
		c.logger.Error(ctx, "Unable to get Company",
			zap.Error(err), zap.String("id", id.String()))
		return nil, err
	}

	return &dto.Company{
		ID:                 company.ID,
		Name:               company.Name,
		RegistrationNumber: company.RegistrationNumber,
		AddressStreet:      company.AddressStreet.String,
		AddressCity:        company.AddressCity.String,
		AddressState:       company.AddressState.String,
		AddressPostalCode:  company.AddressPostalCode.String,
		AddressCountry:     company.AddressCountry.String,
		PrimaryPhone:       company.PrimaryPhone.String,
		SecondaryPhone:     company.SecondaryPhone.String,
		Email:              company.Email.String,
		Status:             company.Status,
		Website:            company.Website.String,
		CallBackURL:        company.CallbackUrl.String,
		ReturnURL:          company.ReturnUrl.String,
		CreatedAt:          company.CreatedAt,
		UpdatedAt:          company.UpdatedAt,
	}, nil
}

func (c *companyPersistance) CreateCompanyToken(ctx context.Context,
	arg dto.CreateCompanyToken) (*dto.CompanyToken, error) {
	token, err := c.persistenceQueries.CreateCompanyToken(ctx,
		db.CreateCompanyTokenParams{
			CompanyID: arg.CompanyID,
			TokenID:   arg.TokenID,
		})
	if err != nil {
		err := errors.ErrUnableToCreate.Wrap(err, "Unable to create Company Token")
		c.logger.Error(ctx, "Unable to create Company Token",
			zap.Error(err))
		return nil, err
	}
	return &dto.CompanyToken{
		ID:        token.ID,
		CompanyID: token.CompanyID,
		TokenID:   token.TokenID,
		Status:    token.Status,
		CreatedAt: token.CreatedAt,
		UpdatedAt: token.UpdatedAt,
		DeletedAt: token.DeletedAt.Time,
	}, nil
}
func (c *companyPersistance) InactiveToken(ctx context.Context,
	companyID uuid.UUID) error {
	err := c.persistenceQueries.Queries.InActiveCompanyToken(ctx, companyID)
	if err != nil {
		err := errors.ErrUnableToUpdate.Wrap(err, "Unable to inactive Token")
		c.logger.Error(ctx, "Unable to inactive Token",
			zap.Error(err), zap.String("companyID", companyID.String()))
		return err
	}
	return nil
}
func (c *companyPersistance) GenerateCompanyCredentials(ctx context.Context,
	arg dto.CreateCompanyToken) error {
	return c.persistenceQueries.GenerateCompanyCredential(ctx,
		dto.CreateCompanyToken{
			CompanyID: arg.CompanyID,
			TokenID:   arg.TokenID,
		})
}

func (c *companyPersistance) GetActiveCompanyTokenByID(ctx context.Context,
	id uuid.UUID) (*dto.CompanyToken, error) {
	token, err := c.persistenceQueries.GetActiveCompanyTokenByID(ctx, id)
	if err != nil {
		if sqlcerr.Is(err, sqlcerr.ErrNoRows) {
			err := errors.ErrNoRecordFound.Wrap(err, "active company token not found")
			c.logger.Error(ctx, "active company token not found", zap.Error(err))
			return nil, err
		}
		err = errors.ErrUnableToGet.Wrap(err, "unable to get active company token")
		c.logger.Error(ctx, "unable to get active company token", zap.Error(err))
		return nil, err
	}

	return &dto.CompanyToken{
		ID:        token.ID,
		TokenID:   token.TokenID,
		CompanyID: token.CompanyID,
		Status:    token.Status,
		CreatedAt: token.CreatedAt,
		UpdatedAt: token.UpdatedAt,
	}, nil
}

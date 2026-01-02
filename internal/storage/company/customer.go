package company

import (
	"context"
	"pg/internal/constant/errors"
	"pg/internal/constant/model/db"
	"pg/internal/constant/model/dto"
	"pg/platform/sql"

	"go.uber.org/zap"
)

func (c *companyPersistance) CreateCustomer(ctx context.Context,
	arg dto.CreateCustomer) (*dto.Customer, error) {
	customer, err := c.persistenceQueries.CreateCustomer(ctx, db.CreateCustomerParams{
		CompanyID:   arg.CompanyID,
		FullName:    sql.StringOrNull(arg.FullName),
		PhoneNumber: arg.PhoneNumber,
		Email:       sql.StringOrNull(arg.Email),
	})
	if err != nil {
		err = errors.ErrUnableToCreate.Wrap(err, "Unable to create customer")
		c.logger.Error(ctx, "Unable to create customer",
			zap.Error(err), zap.String("name", arg.FullName))
		return nil, err
	}

	return &dto.Customer{
		ID:          customer.ID,
		CompanyID:   customer.CompanyID,
		FullName:    customer.FullName.String,
		PhoneNumber: customer.PhoneNumber,
		Email:       customer.Email.String,
		CreatedAt:   customer.CreatedAt,
		UpdatedAt:   customer.UpdatedAt,
	}, nil
}

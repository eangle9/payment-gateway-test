package company

import (
	"context"
	"pg/internal/constant/errors"
	"pg/internal/constant/errors/sqlcerr"
	"pg/internal/constant/model/db"
	"pg/internal/constant/model/dto"
	"pg/platform/sql"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (c *companyPersistance) CreateUser(ctx context.Context, param dto.CreateUser) (*dto.User, error) {
	user, err := c.persistenceQueries.CreateUser(ctx, db.CreateUserParams{
		CompanyID: param.CompanyID,
		FirstName: sql.StringOrNull(param.FirstName),
		Email:     param.Email,
		Phone:     param.Phone,
		Password:  param.Password,
	})
	if err != nil {
		err = errors.ErrUnableToCreate.Wrap(err, "Unable to create user")
		c.logger.Error(ctx, "Unable to create user",
			zap.Error(err), zap.String("name", param.FirstName))
		return nil, err
	}

	return &dto.User{
		ID:        user.ID,
		CompanyID: user.CompanyID,
		FirstName: user.FirstName.String,
		Email:     user.Email,
		Phone:     user.Phone,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}
func (c *companyPersistance) GetUserByID(ctx context.Context, id uuid.UUID) (*dto.User, error) {
	user, err := c.persistenceQueries.GetUserByID(ctx, id)
	if err != nil {
		if sqlcerr.Is(err, sqlcerr.ErrNoRows) {
			err := errors.ErrNoRecordFound.Wrap(err, "user not found")
			c.logger.Error(ctx, "user not found",
				zap.Error(err), zap.String("user-id", id.String()))
			return nil, err
		}
		err := errors.ErrUnableToGet.Wrap(err, "unable to get user")
		c.logger.Error(ctx, "unable to get user",
			zap.Error(err), zap.String("user-id", id.String()))
		return nil, err
	}

	return &dto.User{
		ID:                user.ID,
		CompanyID:         user.CompanyID,
		Username:          user.Username.String,
		Email:             user.Email,
		Password:          user.Password,
		Phone:             user.Phone,
		FirstName:         user.FirstName.String,
		LastName:          user.LastName.String,
		Role:              user.Role.String,
		Status:            user.Status,
		TimezoneID:        user.TimezoneID.String,
		Bio:               user.Bio.String,
		ProfilePictureUrl: user.ProfilePictureUrl.String,
		CreatedAt:         user.CreatedAt,
		UpdatedAt:         user.UpdatedAt,
	}, nil
}

func (c *companyPersistance) GetUserByPhoneOrEmail(ctx context.Context, phone string) (*dto.User, error) {
	user, err := c.persistenceQueries.GetUserByPhoneOrEmail(ctx, phone)
	if err != nil {
		err := errors.ErrUnableToGet.Wrap(err, "unable to get user")
		c.logger.Error(ctx, "unable to get user",
			zap.Error(err), zap.String("phone-or-email", phone))
		return nil, err
	}

	return &dto.User{
		ID:                user.ID,
		CompanyID:         user.CompanyID,
		Username:          user.Username.String,
		Email:             user.Email,
		Password:          user.Password,
		Phone:             user.Phone,
		FirstName:         user.FirstName.String,
		LastName:          user.LastName.String,
		Role:              user.Role.String,
		Status:            user.Status,
		TimezoneID:        user.TimezoneID.String,
		Bio:               user.Bio.String,
		ProfilePictureUrl: user.ProfilePictureUrl.String,
		CreatedAt:         user.CreatedAt,
		UpdatedAt:         user.UpdatedAt,
	}, nil
}

func (c *companyPersistance) CreateUserToken(ctx context.Context, param dto.UserToken) (*dto.UserToken, error) {
	userToken, err := c.persistenceQueries.CreateUserToken(ctx, db.CreateUserTokenParams{
		TokenID: param.TokenID,
		UserID:  param.UserID,
	})
	if err != nil {
		err = errors.ErrUnableToCreate.Wrap(err, "Unable to create user token")
		c.logger.Error(ctx, "Unable to create user token",
			zap.Error(err), zap.String("user-id", param.UserID.String()))
		return nil, err
	}

	return &dto.UserToken{
		ID:        userToken.ID,
		TokenID:   userToken.TokenID,
		UserID:    userToken.UserID,
		Status:    userToken.Status,
		CreatedAt: userToken.CreatedAt,
		UpdatedAt: userToken.UpdatedAt,
	}, nil
}

func (c *companyPersistance) GetActiveUserTokenByUserID(ctx context.Context,
	id uuid.UUID) (*dto.UserToken, error) {
	userToken, err := c.persistenceQueries.GetActiveUserTokenByUserID(ctx, id)
	if err != nil {
		err := errors.ErrUnableToGet.Wrap(err, "unable to get user token")
		c.logger.Error(ctx, "unable to get user token",
			zap.Error(err), zap.String("user-id", id.String()))
		return nil, err
	}

	return &dto.UserToken{
		ID:        userToken.ID,
		TokenID:   userToken.TokenID,
		UserID:    userToken.UserID,
		Status:    userToken.Status,
		CreatedAt: userToken.CreatedAt,
		UpdatedAt: userToken.UpdatedAt,
	}, nil
}

func (c *companyPersistance) ResetActiveToken(ctx context.Context, id uuid.UUID) error {
	if err := c.persistenceQueries.ResetActiveToken(ctx, id); err != nil {
		err = errors.ErrUnableToUpdate.Wrap(err, "unable to reset user token")
		c.logger.Error(ctx, "error while reset user token", zap.Error(err))
		return err
	}
	return nil
}

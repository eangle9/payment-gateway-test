package company

import (
	"context"
	"pg/internal/constant"
	"pg/internal/constant/errors"
	"pg/internal/constant/model/dto"
	"pg/internal/module"
	"pg/internal/storage"
	"pg/platform/hcrypto"
	"pg/platform/hlog"
	"pg/platform/utils"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type company struct {
	log            hlog.Logger
	companyStorage storage.Company
	maker          hcrypto.Maker
}

func New(storage storage.Company,
	log hlog.Logger,
	maker hcrypto.Maker) module.Company {
	return &company{
		companyStorage: storage,
		log:            log,
		maker:          maker,
	}
}
func (c *company) RegisterCompany(ctx context.Context,
	param dto.CreateCompany) (*dto.Company, error) {
	if param.Password != param.ConfirmPassword {
		err := errors.ErrInvalidUserInput.New("Password and Confirm Password does not match")
		c.log.Error(ctx, "password does not match", zap.Error(err))
		return nil, err
	}

	company, err := c.companyStorage.AddCompany(ctx, param)
	if err != nil {
		return nil, err
	}

	// Bcrypt the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(param.Password), bcrypt.DefaultCost)
	if err != nil {
		err = errors.ErrUnableToHashPassword.Wrap(err,
			"Unable to generate password hash")
		c.log.Error(ctx, "unable to generate password hash", zap.Error(err))
		return nil, err
	}
	// Parse phone number
	phone, err := utils.ParsePhoneNumber(param.AdminPhone)
	if err != nil {
		err = errors.ErrInvalidUserInput.Wrap(err, "failed to parse phone number")
		c.log.Error(ctx, "failed to parse phone number",
			zap.Error(err),
			zap.String("phone", param.AdminPhone),
		)
		return nil, err
	}
	param.AdminPhone = *phone
	if _, err = c.companyStorage.CreateUser(ctx, dto.CreateUser{
		CompanyID: company.ID,
		FirstName: param.AdminName,
		Email:     param.AdminEmail,
		Phone:     param.AdminPhone,
		Password:  string(hashedPassword),
	}); err != nil {
		return nil, err
	}

	return company, nil
}

func (c *company) Login(ctx context.Context, arg dto.LoginRequest) (*dto.SignInResponse, error) {
	if err := arg.Validate(); err != nil {
		err = errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		c.log.Error(ctx, "invalid input", zap.Error(err))
		return nil, err
	}
	user, err := c.companyStorage.GetUserByPhoneOrEmail(ctx, arg.PhoneOrEmail)
	if err != nil {
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(arg.Password)); err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "incorrect password")
		c.log.Error(ctx, "incorrect password", zap.Error(err))
		return nil, err
	}
	accessToken, _, err := c.maker.CreatePasetoToken(hcrypto.UserData{
		UserID:    user.ID.String(),
		Email:     user.Email,
		IsNewUser: false,
		Provider:  constant.Normal,
	}, constant.AccessToken)
	if err != nil {
		err = errors.ErrInternalServerError.Wrap(err, "unable to generate access token")
		c.log.Error(ctx, "unable to generate access token", zap.Error(err))
		return nil, err
	}
	refreshToken, tokenID, err := c.maker.CreatePasetoToken(hcrypto.UserData{
		UserID:    user.ID.String(),
		Email:     user.Email,
		IsNewUser: false,
		Provider:  constant.Normal,
	}, constant.RefreshToken)
	if err != nil {
		err = errors.ErrInternalServerError.Wrap(err, "unable to generate refresh token")
		c.log.Error(ctx, "unable to generate refresh token", zap.Error(err))
		return nil, err
	}
	if err := c.companyStorage.ResetActiveToken(ctx, user.ID); err != nil {
		return nil, err
	}
	if _, err := c.companyStorage.CreateUserToken(ctx, dto.UserToken{
		TokenID: tokenID,
		UserID:  user.ID,
	}); err != nil {
		return nil, err
	}

	return &dto.SignInResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (c *company) GenerateToken(ctx context.Context,
	userID string) (*dto.CompanyCredentialResponse, error) {
	id, err := uuid.Parse(userID)
	if err != nil {
		err = errors.ErrInvalidUserInput.Wrap(err, "Invalid user id")
		c.log.Error(ctx, "Invalid user id", zap.Error(err))
		return nil, err
	}
	user, err := c.companyStorage.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	company, err := c.companyStorage.GetCompanyByID(ctx, user.CompanyID)
	if err != nil {
		return nil, err
	}
	secret, tokenID, err := c.maker.CreatePasetoToken(hcrypto.UserData{
		UserID:    company.ID.String(),
		Email:     company.Email,
		IsNewUser: false,
		Provider:  constant.Normal,
	}, constant.SecretToken)
	if err != nil {
		return nil, err
	}
	if err := c.companyStorage.GenerateCompanyCredentials(ctx, dto.CreateCompanyToken{
		TokenID:   tokenID,
		CompanyID: user.CompanyID,
	}); err != nil {
		return nil, err
	}
	return &dto.CompanyCredentialResponse{
		ScretToken: secret,
	}, nil
}

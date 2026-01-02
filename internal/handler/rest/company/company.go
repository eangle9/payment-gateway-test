package company

import (
	"context"
	"net/http"
	"pg/internal/constant/errors"
	"pg/internal/constant/model/dto"
	"pg/internal/constant/model/response"
	"pg/internal/handler/rest"
	"pg/internal/module"
	"pg/platform/hlog"
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type company struct {
	log            hlog.Logger
	companyModule  module.Company
	contextTimeout time.Duration
}

func New(log hlog.Logger, companyModule module.Company,
	ctx time.Duration) rest.Company {
	return &company{
		log:            log,
		companyModule:  companyModule,
		contextTimeout: ctx,
	}

}

// RegisterCompany
//
//	@Summary		Create a new company
//	@Description	This endpoint registers a new company in the Activity Rewards application.
//	@Tags			company
//	@Accept			json
//	@Produce		json
//	@Param			create_company_request_body	body		dto.CreateCompany	true	"Company details for registration"
//	@Success		201							{object}	doc.SuccessResponse{data=dto.Company,meta_data=interface{}}
//	@Failure		400							{object}	doc.ErrorResponse	"Bad request due to invalid input"
//	@Failure		401							{object}	doc.ErrorResponse	"Unauthorized request"
//	@Failure		500							{object}	doc.ErrorResponse	"Internal server error"
//	@Router			/signup-company-owner [post]
//	@Security		BearerAuth
func (cr *company) RegisterCompany(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), cr.contextTimeout)
	defer cancel()

	param := dto.CreateCompany{}
	err := c.Bind(&param)
	if err != nil {
		er := errors.ErrBadRequest.Wrap(err, "unable to bind company data")
		cr.log.Error(ctx, "unable to bind company data", zap.Error(err))
		return er
	}

	data, err := cr.companyModule.RegisterCompany(ctx, param)
	if err != nil {
		return err
	}

	return response.SendSuccessResponse(c, http.StatusCreated, data, nil)
}

// Login
//
//	@Summary		Authenticate a company
//	@Description	This endpoint allows a company to log in to the application.
//	@Tags			company
//	@Accept			json
//	@Produce		json
//	@Param			company_login_request_body	body		dto.LoginRequest	true	"Login credentials of the company"
//	@Success		200							{object}	doc.SuccessResponse{data=dto.SignInResponse,meta_data=interface{}}
//	@Failure		400							{object}	doc.ErrorResponse	"Bad request due to invalid input"
//	@Failure		401							{object}	doc.ErrorResponse	"Unauthorized request"
//	@Failure		500							{object}	doc.ErrorResponse	"Internal server error"
//	@Router			/login [post]
func (cr *company) Login(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), cr.contextTimeout)
	defer cancel()

	param := dto.LoginRequest{}
	err := c.Bind(&param)
	if err != nil {
		er := errors.ErrBadRequest.Wrap(err, "Unable to bind login data")
		cr.log.Error(ctx, "Unable to bind login data", zap.Error(err))
		return er
	}

	data, err := cr.companyModule.Login(ctx, param)
	if err != nil {
		return err
	}

	return response.SendSuccessResponse(c, http.StatusOK, data, nil)
}

// GenerateSecretToken
//
//	@Summary		Get secret token
//	@Description	Get a company secret token.
//	@Tags			company
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	doc.SuccessResponse{data=dto.CompanyCredentialResponse,meta_data=interface{}}
//	@Failure		400	{object}	doc.ErrorResponse	"Bad request due to invalid input"
//	@Failure		401	{object}	doc.ErrorResponse	"Unauthorized request"
//	@Failure		500	{object}	doc.ErrorResponse	"Internal server error"
//	@Router			/generate-secret-token [post]
//	@Security		BearerAuth
func (cr *company) GenerateSecretToken(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), cr.contextTimeout)
	defer cancel()

	id, ok := ctx.Value("x-id").(string)
	if !ok {
		err := errors.ErrInvalidUserInput.New(
			"invalid user id, it could be type of string")
		return err
	}

	data, err := cr.companyModule.GenerateToken(ctx, id)
	if err != nil {
		return err
	}

	return response.SendSuccessResponse(c, http.StatusOK, data, nil)
}

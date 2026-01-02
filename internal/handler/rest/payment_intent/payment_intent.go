package paymentintent

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

type paymentIntent struct {
	log                 hlog.Logger
	PaymentIntentModule module.PaymentIntent
	contextTimeout      time.Duration
}

func New(log hlog.Logger, PaymentIntentModule module.PaymentIntent,
	ctx time.Duration) rest.PaymentIntent {
	return &paymentIntent{
		log:                 log,
		PaymentIntentModule: PaymentIntentModule,
		contextTimeout:      ctx,
	}

}

// Initiate PaymentIntent
//
//	@Summary		InitPaymentIntent
//	@Description	Initiate onetme paymentIntent
//	@Tags			payments
//	@Accept			json
//	@Produce		json
//	@Param			create_payment_intent_request_body	body		dto.InitPaymentIntent	true	"payment-intent details"
//	@Success		201									{object}	doc.SuccessResponse{data=dto.PaymentIntent,meta_data=interface{}}
//	@Failure		400									{object}	doc.ErrorResponse	"Bad request due to invalid input"
//	@Failure		401									{object}	doc.ErrorResponse	"Unauthorized request"
//	@Failure		500									{object}	doc.ErrorResponse	"Internal server error"
//	@Router			/payment-intents [post]
//	@Security		BearerAuth
func (p *paymentIntent) InitPaymentIntent(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), p.contextTimeout)
	defer cancel()

	id, ok := ctx.Value("x-companyID").(string)
	if !ok {
		err := errors.ErrInvalidUserInput.New("invalid company id, it could be type of string")
		p.log.Error(ctx, "invalid company id", zap.Error(err))
		return err
	}

	param := dto.InitPaymentIntent{}
	err := c.Bind(&param)
	if err != nil {
		er := errors.ErrBadRequest.Wrap(err, "unable to bind payment intent data")
		p.log.Error(ctx, "unable to bind payment intent data", zap.Error(err))
		return er
	}

	data, err := p.PaymentIntentModule.InitPaymentIntent(ctx, param, id)
	if err != nil {
		return err
	}

	return response.SendSuccessResponse(c, http.StatusCreated, data, nil)
}

// Get PaymentIntent Detail By ID
//
//	@Summary		Get PaymentIntent By ID
//	@Description	Get payment details by id
//	@Tags			payments
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"payment intent id"
//	@Success		200	{object}	doc.SuccessResponse{data=dto.PaymentIntentDetail,meta_data=interface{}}
//	@Failure		400	{object}	doc.ErrorResponse	"Bad request due to invalid input"
//	@Failure		401	{object}	doc.ErrorResponse	"Unauthorized request"
//	@Failure		500	{object}	doc.ErrorResponse	"Internal server error"
//	@Router			/payment-intents/{id} [get]
//	@Security		BearerAuth
func (p *paymentIntent) GetPaymentIntentDetail(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), p.contextTimeout)
	defer cancel()

	data, err := p.PaymentIntentModule.GetPaymentIntentDetail(ctx, c.Param("id"))
	if err != nil {
		return err
	}

	return response.SendSuccessResponse(c, http.StatusOK, data, nil)
}

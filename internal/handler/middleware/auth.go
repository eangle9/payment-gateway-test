package middleware

import (
	"context"
	"net/http"
	"pg/internal/constant"
	"pg/internal/constant/errors"
	"pg/internal/storage"
	"pg/platform/hcrypto"
	"pg/platform/hlog"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type AuthMiddleware interface {
	AuthenticateAdminUser() echo.MiddlewareFunc
	AuthenticateUser() echo.MiddlewareFunc
}

type authMiddleware struct {
	logger         hlog.Logger
	maker          hcrypto.Maker
	companyStorage storage.Company
}

func InitAuthMiddleware(
	logger hlog.Logger,
	maker hcrypto.Maker,
	companyStorage storage.Company,

) AuthMiddleware {
	return &authMiddleware{
		logger:         logger,
		maker:          maker,
		companyStorage: companyStorage,
	}
}
func (a *authMiddleware) AuthenticateAdminUser() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := c.Request().Context()
			payload, err := a.VerifyPasetoToken(c)
			if err != nil {
				err = errors.ErrInvalidAccessToken.Wrap(err, "invalid token")
				a.logger.Error(ctx, "invalid token", zap.Error(err))
				return c.JSON(http.StatusUnauthorized, err)
			}
			companyID, err := uuid.Parse(payload.UserID)
			if err != nil {
				err = errors.ErrInternalServerError.Wrap(err, "invalid company id")
				a.logger.Error(ctx, "error parsing company id", zap.Error(err))
				return c.JSON(http.StatusInternalServerError, err)
			}
			company, err := a.companyStorage.GetCompanyByID(ctx, companyID)
			if err != nil {
				err = errors.ErrInvalidAccessToken.New("access denied")
				return c.JSON(http.StatusUnauthorized, err)
			}
			if company.Status != string(constant.Active) {
				err = errors.ErrAuthError.New("access denied company status is %s", company.Status)
				return c.JSON(http.StatusUnauthorized, err)
			}
			companyToken, err := a.companyStorage.GetActiveCompanyTokenByID(ctx, companyID)
			if err != nil {
				err = errors.ErrInvalidAccessToken.New("unable to get active company token")
				a.logger.Error(ctx, "unable to get active company token")
				return c.JSON(http.StatusUnauthorized, err)
			}
			if payload.TokenID != companyToken.TokenID {
				err = errors.ErrInvalidAccessToken.New("invalid token")
				a.logger.Error(ctx, "invalid token", zap.Error(err))
				return c.JSON(http.StatusUnauthorized, err)
			}

			req := c.Request()
			req = req.WithContext(context.WithValue(req.Context(), constant.ContextKey("x-companyID"), company.ID.String()))
			req = req.WithContext(context.WithValue(req.Context(), constant.ContextKey("x-company"), *company))
			req = req.WithContext(context.WithValue(req.Context(), constant.ContextKey(constant.AuthorizationPayloadKey), *payload))
			c.SetRequest(req)

			return next(c)
		}
	}
}
func (a *authMiddleware) AuthenticateUser() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := c.Request().Context()
			payload, err := a.VerifyPasetoToken(c)
			if err != nil {
				err = errors.ErrInvalidAccessToken.Wrap(err, "invalid token")
				a.logger.Error(ctx, "invalid token", zap.Error(err))
				return c.JSON(http.StatusUnauthorized, err)
			}
			userID, err := uuid.Parse(payload.UserID)
			if err != nil {
				err = errors.ErrInternalServerError.Wrap(err, "invalid user id")
				a.logger.Error(ctx, "error parsing user id", zap.Error(err))
				return c.JSON(http.StatusInternalServerError, err)
			}
			user, err := a.companyStorage.GetUserByID(ctx, userID)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, err)
			}
			if user.Status != string(constant.Active) {
				err = errors.ErrAuthError.New("access denied")
				return c.JSON(http.StatusUnauthorized, err)
			}

			req := c.Request()
			req = req.WithContext(context.WithValue(req.Context(), constant.ContextKey("x-id"), payload.UserID))
			req = req.WithContext(context.WithValue(req.Context(), constant.ContextKey(constant.AuthorizationPayloadKey), *payload))
			c.SetRequest(req)

			return next(c)
		}
	}
}

func (a *authMiddleware) VerifyPasetoToken(c echo.Context) (*hcrypto.Payload, error) {
	ctx := c.Request().Context()
	authorizationHeader := c.Request().Header.Get(constant.AuthorizationHeaderkey)
	if len(authorizationHeader) == 0 {
		err := errors.ErrAuthError.New("authorization Header is not provided")
		a.logger.Error(ctx, "authorization Header is not provided")
		return nil, err
	}

	fields := strings.Fields(authorizationHeader)
	if len(fields) < 2 {
		err := errors.ErrAuthError.New("invalid authorization header format")
		a.logger.Error(ctx, "invalid authorization header format")
		return nil, err
	}
	authorizationType := fields[0]
	if authorizationType != constant.AuthorizationTypeBearer {
		err := errors.ErrAuthError.New("unsupported authorization type")
		a.logger.Error(ctx, "unsupported authorization type: "+authorizationType, zap.Error(err))
		return nil, err
	}
	accessToken := fields[1]
	payload, err := a.maker.VerifyPasetoToken(accessToken)
	if payload == nil {
		a.logger.Error(ctx, "the payload is nil", zap.Error(err))
		return nil, err
	}
	if err != nil {
		err = errors.ErrAuthError.Wrap(err, "unsupported authorization type")
		a.logger.Error(ctx, "unsupported authorization type", zap.Error(err))
		return nil, err
	}
	return payload, nil
}

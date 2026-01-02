package errors

import (
	"net/http"

	"github.com/joomcode/errorx"
)

// list of error namespaces
var (
	databaseError    = errorx.NewNamespace("database error").ApplyModifiers(errorx.TypeModifierOmitStackTrace)
	invalidInput     = errorx.NewNamespace("validation error").ApplyModifiers(errorx.TypeModifierOmitStackTrace)
	resourceNotFound = errorx.NewNamespace("not found").ApplyModifiers(errorx.TypeModifierOmitStackTrace)
	AccessDenied     = errorx.RegisterTrait("You are not authorized to perform the action")
	Ineligible       = errorx.RegisterTrait("You are not eligible to perform the action")
	serverError      = errorx.NewNamespace("server error")
	httpError        = errorx.NewNamespace("http error")
	badRequest       = errorx.NewNamespace("bad request error")
	Unauthenticated  = errorx.NewNamespace("user authentication failed")
	unauthorized     = errorx.NewNamespace("unauthorized").ApplyModifiers(errorx.TypeModifierOmitStackTrace)
	redisServerError = errorx.NewNamespace("redis service error")
)

var (
	ErrUnableToUploadFile       = errorx.NewType(databaseError, "unable to upload file")
	ErrDrawFailed               = errorx.NewType(serverError, "draw failed")
	ErrUserNotFound             = errorx.NewType(resourceNotFound, "user not found")
	ErrAcessError               = errorx.NewType(unauthorized, "Unauthorized", AccessDenied)
	ErrInvalidUserInput         = errorx.NewType(invalidInput, "invalid user input")
	ErrUnableToGet              = errorx.NewType(databaseError, "unable to get")
	ErrInternalServerError      = errorx.NewType(serverError, "internal server error")
	ErrUnableToUpdate           = errorx.NewType(databaseError, "unable to update")
	ErrUnableToCreate           = errorx.NewType(databaseError, "unable to create")
	ErrUnableToReset            = errorx.NewType(databaseError, "unable to reset")
	ErrUnableToSendMail         = errorx.NewType(databaseError, "unable to send mail")
	ErrUnableToHashPassword     = errorx.NewType(databaseError, "unable to hash password")
	ErrDBDelError               = errorx.NewType(databaseError, "could not delete record")
	ErrNoRecordFound            = errorx.NewType(resourceNotFound, "no record found")
	ErrHTTPRequestPrepareFailed = errorx.NewType(httpError, "couldn't prepare http request")
	ErrBadRequest               = errorx.NewType(badRequest, "bad request error")
	ErrSSOAuthenticationFailed  = errorx.NewType(Unauthenticated, "user authentication failed")
	ErrSSOError                 = errorx.NewType(serverError, "sso communication failed")

	ErrInvalidAccessToken = errorx.NewType(Unauthenticated, "invalid token").
				ApplyModifiers(errorx.TypeModifierOmitStackTrace)
	ErrAuthError                   = errorx.NewType(unauthorized, "you are not authorized.")
	ErrHashError                   = errorx.NewType(serverError, "error generating hash")
	ErrRedisPubSubUnableToGetRoute = errorx.NewType(redisServerError, "redis pubsub error")
	ErrReadError                   = errorx.NewType(redisServerError, "redis read message error")
	ErrRedisPubSubBadEvent         = errorx.NewType(redisServerError, "redis read message error")
	ErrUnsupportedPublicKeyFormat  = errorx.NewType(invalidInput, "unsupported public key format")
)

var ErrorMap = map[*errorx.Type]int{
	ErrUnableToUploadFile:          http.StatusBadRequest,
	ErrDrawFailed:                  http.StatusInternalServerError,
	ErrUnableToSendMail:            http.StatusInternalServerError,
	ErrUnableToReset:               http.StatusInternalServerError,
	ErrUserNotFound:                http.StatusNotFound,
	ErrReadError:                   http.StatusBadRequest,
	ErrRedisPubSubBadEvent:         http.StatusBadRequest,
	ErrRedisPubSubUnableToGetRoute: http.StatusInternalServerError,
	ErrAcessError:                  http.StatusForbidden,
	ErrInvalidUserInput:            http.StatusBadRequest,
	ErrHashError:                   http.StatusInternalServerError,
	ErrInternalServerError:         http.StatusInternalServerError,
	ErrUnableToGet:                 http.StatusInternalServerError,
	ErrNoRecordFound:               http.StatusNotFound,
	ErrDBDelError:                  http.StatusInternalServerError,
	ErrUnableToUpdate:              http.StatusInternalServerError,
	ErrUnableToCreate:              http.StatusInternalServerError,
	ErrUnableToHashPassword:        http.StatusInternalServerError,
	ErrHTTPRequestPrepareFailed:    http.StatusInternalServerError,
	ErrBadRequest:                  http.StatusBadRequest,
	ErrSSOAuthenticationFailed:     http.StatusUnauthorized,
	ErrSSOError:                    http.StatusInternalServerError,
	ErrInvalidAccessToken:          http.StatusUnauthorized,
	ErrAuthError:                   http.StatusUnauthorized,
	ErrUnsupportedPublicKeyFormat:  http.StatusBadRequest,
}

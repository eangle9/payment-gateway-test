package middleware

import (
	"fmt"
	"net/http"
	"pg/internal/constant/errors"
	response2 "pg/internal/constant/model/response"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/joomcode/errorx"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

func ErrorHandler(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}

	response2.SendErrorResponse(c, CastErrorResponse(err))
}
func ErrorFields(err error) []response2.FieldError {
	var errs []response2.FieldError

	if data, ok := err.(validation.Errors); ok {
		for i, v := range data {
			errs = append(errs, response2.FieldError{
				Name:        i,
				Description: v.Error(),
			},
			)
		}

		return errs
	}

	return nil
}

func CastErrorResponse(err error) *response2.ErrorResponse {
	debugMode := viper.GetBool("debug")
	er := errorx.Cast(err)
	if er == nil {
		return &response2.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Unknown server error",
		}
	}

	response := response2.ErrorResponse{}
	code, ok := errors.ErrorMap[er.Type()]
	if !ok {
		response = response2.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Unknown server error",
		}
	} else {
		response = response2.ErrorResponse{
			Code:       code,
			Message:    er.Message(),
			FieldError: ErrorFields(er.Cause()),
		}
	}

	if debugMode {
		response.Description = fmt.Sprintf("Error: %v", er)
		response.StackTrace = fmt.Sprintf("%+v", errorx.EnsureStackTrace(err))
	}

	return &response
}

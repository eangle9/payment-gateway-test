package response

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v4"
)

func SendSuccessResponse(ctx echo.Context, statusCode int, data interface{},
	metaData *MetaData) error {

	return ctx.JSON(
		statusCode,
		Response{
			Success:  true,
			MetaData: metaData,
			Data:     data,
		},
	)
}

func SendErrorResponse(ctx echo.Context, err *ErrorResponse) error {
	return ctx.JSON(err.Code, Response{
		Success: false,
		Error:   err,
	})
}

func ErrorFields(err error) []FieldError {
	var errs []FieldError

	if data, ok := err.(validation.Errors); ok {
		for i, v := range data {
			errs = append(errs, FieldError{
				Name:        i,
				Description: v.Error(),
			},
			)
		}

		return errs
	}

	return nil
}

package dto

import (
	"errors"
	"fmt"
	"pg/internal/constant"
	"time"

	"github.com/dongri/phonenumber"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type PaymentIntent struct {
	ID          uuid.UUID            `json:"id,omitempty"`
	CompanyID   uuid.UUID            `json:"company_id,omitempty"`
	CustomerID  uuid.UUID            `json:"customer_id,omitempty"`
	PaymentType constant.PaymentType `json:"payment_type,omitempty"`
	Amount      decimal.Decimal      `json:"amount,omitempty"`
	Status      constant.Status      `json:"status,omitempty"`
	Currency    constant.Currency    `json:"currency,omitempty"`
	CallBackURL string               `json:"callback_url,omitempty"`
	ReturnURL   string               `json:"return_url,omitempty"`
	Description string               `json:"description,omitempty"`
	Extra       map[string]any       `json:"extra,omitempty"`
	BillRefNO   string               `json:"bill_ref_no,omitempty"`
	TopayURL    string               `json:"topay_url,omitempty"`
	ExpireAt    time.Time            `json:"expire_at,omitempty"`
	CreatedAt   time.Time            `json:"created_at,omitempty"`
	UpdatedAt   time.Time            `json:"updated_at,omitempty"`
}

type PaymentIntentDetail struct {
	ID          uuid.UUID            `json:"id,omitempty"`
	PaymentType constant.PaymentType `json:"payment_type,omitempty"`
	Amount      decimal.Decimal      `json:"amount,omitempty"`
	Status      constant.Status      `json:"status,omitempty"`
	Currency    constant.Currency    `json:"currency,omitempty"`
	CallBackURL string               `json:"callback_url,omitempty"`
	ReturnURL   string               `json:"return_url,omitempty"`
	Description string               `json:"description,omitempty"`
	Extra       map[string]any       `json:"extra,omitempty"`
	BillRefNO   string               `json:"bill_ref_no,omitempty"`
	TopayURL    string               `json:"topay_url,omitempty"`
	ExpireAt    time.Time            `json:"expire_at,omitempty"`
	CreatedAt   time.Time            `json:"created_at,omitempty"`
	UpdatedAt   time.Time            `json:"updated_at,omitempty"`
	Customer    Customer             `json:"customer,omitempty"`
	Company     Company              `json:"company,omitempty"`
}

type InitPaymentIntent struct {
	Amount      decimal.Decimal `json:"amount,omitempty"`
	Currency    string          `json:"currency,omitempty"`
	CallBackURL string          `json:"callback_url,omitempty"`
	ReturnURL   string          `json:"return_url,omitempty"`
	Description string          `json:"description,omitempty"`
	Customer    PaymentCustomer `json:"customer,omitempty"`
	Extra       map[string]any  `json:"extra,omitempty"`
}

func (c InitPaymentIntent) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Amount, validation.Required.Error("Amount is required"),
			validation.By(ValidateDecimalMin(decimal.NewFromInt(1),
				fmt.Sprintf("value must be greater than or equal to %v", 1)))),
		validation.Field(&c.Currency, validation.Required.Error("Currency id is required")),
		validation.Field(&c.Customer, validation.By(func(value interface{}) error {
			customer, ok := value.(PaymentCustomer)
			if !ok {
				return errors.New("invalid customer details")
			}

			if customer.PhoneNumber == "" {
				return errors.New("phone number is required")
			}

			if emailErr := validation.Validate(customer.Email,
				validation.When(customer.Email != "",
					validation.Required.Error("invalid email provided"),
					is.EmailFormat.Error("invalid email provided"))); emailErr != nil {
				return emailErr
			}

			return ValidatePhone(customer.PhoneNumber)
		})),
		validation.Field(&c.ReturnURL,
			validation.When(c.ReturnURL != "", is.URL.Error("invalid return url provided"))),
		validation.Field(&c.CallBackURL,
			validation.When(c.CallBackURL != "", is.URL.Error("invalid callback url provided")),
		),
	)
}

func ValidateDecimalMin(minValue decimal.Decimal, message string) validation.RuleFunc {
	return func(value interface{}) error {
		val, ok := value.(decimal.Decimal)
		if !ok {
			return fmt.Errorf("value must be a Decimal type: %T", value)
		}

		if val.LessThan(minValue) {
			return validation.NewError("400", message)
		}

		return nil
	}
}

func ValidatePhone(phone any) error {
	str := phonenumber.Parse(fmt.Sprintf("%v", phone), "ET")
	if str == "" {
		return errors.New("invalid phone number")
	}

	return nil
}

type PaymentCustomer struct {
	FullName    string `json:"full_name,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`
	Email       string `json:"email,omitempty"`
}

type CreatePaymentIntent struct {
	CompanyID   uuid.UUID            `json:"company_id,omitempty"`
	PaymentType constant.PaymentType `json:"payment_type,omitempty"`
	Amount      decimal.Decimal      `json:"amount,omitempty"`
	Status      constant.Status      `json:"status,omitempty"`
	Currency    constant.Currency    `json:"currency,omitempty"`
	CallBackURL string               `json:"callback_url,omitempty"`
	ReturnURL   string               `json:"return_url,omitempty"`
	Description string               `json:"description,omitempty"`
	Customer    PaymentCustomer      `json:"customer,omitempty"`
	Extra       map[string]any       `json:"extra,omitempty"`
	BillRefNO   string               `json:"bill_ref_no,omitempty"`
}

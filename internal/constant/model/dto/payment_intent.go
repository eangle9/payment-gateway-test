package dto

import (
	"pg/internal/constant"
	"time"

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
	CustomerID  uuid.UUID            `json:"customer_id,omitempty"`
	Extra       map[string]any       `json:"extra,omitempty"`
	BillRefNO   string               `json:"bill_ref_no,omitempty"`
}

package hcrypto

import (
	"errors"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
)

type TokenKey struct {
	SymmetricKey       string `json:"symmetric_key"`
	Issuer             string `json:"issuer"`
	Footer             string `json:"footer"`
	KeyLength          int    `json:"key_length"`
	Audience           string `json:"audience"`
	AccessExpires      time.Duration
	RefreshExpires     time.Duration
	SecretTokenExpires time.Duration
}

func (m TokenKey) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.SymmetricKey, validation.Required.Error("Symmetric-key is required")),
		validation.Field(&m.Issuer, validation.Required.Error("issuer is required")),
		validation.Field(&m.Footer, validation.Required.Error("footer short code is required")),
		validation.Field(&m.KeyLength, validation.Required.Error("key length is required")),
		validation.Field(&m.Audience, validation.Required.Error("audience is required")),
		validation.Field(&m.AccessExpires, validation.Required.Error("access expires is required"),
			validation.By(func(value interface{}) error {
				expires, ok := value.(time.Duration)
				if !ok {
					return errors.New("access expires should be time.Duration")
				}
				if expires <= 0 {
					return errors.New("access expires should be positive")
				}
				return nil
			})),
		validation.Field(&m.RefreshExpires, validation.Required.Error("refresh expires is required"),
			validation.By(func(value interface{}) error {
				expires, ok := value.(time.Duration)
				if !ok {
					return errors.New("access expires should be time.Duration")
				}
				if expires <= 0 {
					return errors.New("access expires should be positive")
				}
				return nil
			})),
	)
}

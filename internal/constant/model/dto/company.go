package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateCompanyToken struct {
	TokenID   uuid.UUID `json:"token_id"`
	CompanyID uuid.UUID `json:"company_id"`
}

type Company struct {
	ID                 uuid.UUID `json:"id"`
	Name               string    `json:"name"`
	RegistrationNumber string    `json:"registration_number"`
	AddressStreet      string    `json:"address_street"`
	AddressCity        string    `json:"address_city"`
	AddressState       string    `json:"address_state"`
	AddressPostalCode  string    `json:"address_postal_code"`
	AddressCountry     string    `json:"address_country"`
	PrimaryPhone       string    `json:"primary_phone"`
	SecondaryPhone     string    `json:"secondary_phone"`
	Status             string    `json:"status"`
	Email              string    `json:"email"`
	Website            string    `json:"website"`
	CallBackURL        string    `json:"callback_url"`
	ReturnURL          string    `json:"return_url"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type CompanyToken struct {
	ID        uuid.UUID `json:"id"`
	TokenID   uuid.UUID `json:"token_id"`
	CompanyID uuid.UUID `json:"company_id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

type CreateCompany struct {
	Name               string `json:"name" example:"Acme Technologies Ltd"`
	RegistrationNumber string `json:"registration_number" example:"REG-123456"`
	AddressStreet      string `json:"address_street" example:"Bole Road"`
	AddressCity        string `json:"address_city" example:"Addis Ababa"`
	AddressState       string `json:"address_state" example:"Addis Ababa"`
	AddressPostalCode  string `json:"address_postal_code" example:"1000"`
	AddressCountry     string `json:"address_country" example:"Ethiopia"`
	PrimaryPhone       string `json:"primary_phone" example:"+251911234567"`
	SecondaryPhone     string `json:"secondary_phone" example:"+251922345678"`
	Email              string `json:"email" example:"info@acmetech.com"`
	Website            string `json:"website" example:"https://www.acmetech.com"`
	AdminName          string `json:"admin_name" example:"John Doe"`
	AdminEmail         string `json:"admin_email" example:"admin@acmetech.com"`
	AdminPhone         string `json:"admin_phone" example:"+251933456789"`
	Password           string `json:"password" example:"StrongPass@123"`
	ConfirmPassword    string `json:"confirm_password" example:"StrongPass@123"`
	CallBackURL        string `json:"callback_url" example:"https://www.acmetech.com/payment/callback"`
	ReturnURL          string `json:"return_url" example:"https://www.acmetech.com/payment/return"`
}

type CompanyCredentialResponse struct {
	ScretToken string `json:"scret_token"`
}

type LoginRequest struct {
	PhoneOrEmail string `json:"phone" example:"+251933456789"`
	Password     string `json:"password" example:"StrongPass@123"`
}

type SignInResponse struct {
	// Access token for the user.
	AccessToken string `json:"access" example:"access-token"`
	// Refresh token for the user.
	RefreshToken string `json:"refresh" example:"refresh-token"`
}

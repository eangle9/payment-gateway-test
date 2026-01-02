package dto

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID                uuid.UUID `json:"id,omitempty"`
	CompanyID         uuid.UUID `json:"company_id,omitempty"`
	Username          string    `json:"username,omitempty"`
	Email             string    `json:"email,omitempty"`
	Password          string    `json:"password,omitempty"`
	Phone             string    `json:"phone,omitempty"`
	FirstName         string    `json:"first_name,omitempty"`
	LastName          string    `json:"last_name,omitempty"`
	Role              string    `json:"role,omitempty"`
	Status            string    `json:"status,omitempty"`
	TimezoneID        string    `json:"timezone_id,omitempty"`
	Bio               string    `json:"bio,omitempty"`
	ProfilePictureUrl string    `json:"profile_picture_url,omitempty"`
	CreatedAt         time.Time `json:"created_at,omitempty"`
	UpdatedAt         time.Time `json:"updated_at,omitempty"`
}

type CreateUser struct {
	CompanyID uuid.UUID `json:"company_id"`
	FirstName string    `json:"first_name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Password  string    `json:"password"`
}

type UserToken struct {
	ID        uuid.UUID `json:"id,omitempty"`
	TokenID   uuid.UUID `json:"token_id,omitempty"`
	UserID    uuid.UUID `json:"user_id,omitempty"`
	Status    string    `json:"status,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	DeletedAt time.Time `json:"deleted_at,omitempty"`
}

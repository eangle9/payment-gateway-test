package dto

import (
	"time"

	"github.com/google/uuid"
)

type Customer struct {
	ID          uuid.UUID `json:"id,omitempty" example:"550e8400-e29b-41d4-a716-446655440001"`
	CompanyID   uuid.UUID `json:"company_id,omitempty" example:"550e8400-e29b-41d4-a716-446655440002"`
	FullName    string    `json:"full_name,omitempty" example:"John Doe"`
	PhoneNumber string    `json:"phone_number,omitempty" example:"+1234567890"`
	Email       string    `json:"email,omitempty" example:"john.doe@gmail.com"`
	CreatedAt   time.Time `json:"created_at,omitempty" example:"2023-09-11T14:30:00Z"`
	UpdatedAt   time.Time `json:"updated_at,omitempty" example:"2023-09-11T14:45:00Z"`
}

type CreateCustomer struct {
	CompanyID   uuid.UUID `json:"company_id,omitempty"`
	FullName    string    `json:"full_name,omitempty"`
	PhoneNumber string    `json:"phone_number,omitempty"`
	Email       string    `json:"email,omitempty"`
}

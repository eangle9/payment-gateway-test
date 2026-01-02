-- name: CreateCustomer :one
INSERT INTO customers (
  company_id,
  full_name,
  phone_number,
  email
) VALUES (
  $1, $2, $3, $4
) ON CONFLICT (company_id, phone_number) WHERE deleted_at IS NULL
DO UPDATE SET email = excluded.email, full_name = excluded.full_name
RETURNING *;
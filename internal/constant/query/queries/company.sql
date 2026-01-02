-- name: CreateCompany :one
INSERT INTO companies (
  name,
  registration_number,
  address_street,
  address_city,
  address_state,
  address_postal_code,
  address_country,
  primary_phone,
  secondary_phone,
  email,
  website,
  callback_url,
  return_url
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
)
RETURNING id, name, registration_number, address_street, address_city, address_state, address_postal_code, address_country, primary_phone, secondary_phone, email, status, website, callback_url, return_url, created_at, updated_at;

-- name: GetCompanyByID :one
SELECT id, name, registration_number, address_street, address_city, address_state, address_postal_code, address_country, primary_phone, secondary_phone, email, status, website, callback_url, return_url, created_at, updated_at
FROM companies
WHERE id = $1 AND deleted_at IS NULL;

-- name: InActiveCompanyToken :exec
UPDATE company_tokens ct 
SET status = 'INACTIVE'
WHERE ct.company_id = $1 
AND ct.deleted_at IS NULL 
AND ct.status = 'ACTIVE';

-- name: CreateCompanyToken :one
INSERT INTO company_tokens (token_id, company_id) 
VALUES ($1, $2) 
RETURNING *;

-- name: GetActiveCompanyTokenByID :one
SELECT * 
FROM company_tokens
WHERE company_id = $1 AND status = 'ACTIVE' AND deleted_at IS NULL;

-- name: CreateUser :one
INSERT INTO users (
  company_id,
  first_name,
  email,
  phone,
  password
)
VALUES (
  $1, $2, $3, $4, $5
)
RETURNING *;
-- name: GetUserByID :one
SELECT *
FROM users 
WHERE id = $1 AND deleted_at IS NULL;
-- name: GetUserByPhoneOrEmail :one
SELECT *
FROM users 
WHERE (email = $1 OR phone = $1) AND deleted_at IS NULL;
-- name: UsernameExists :one
SELECT EXISTS (
  SELECT 1 
  FROM users 
  WHERE username = $1 AND deleted_at IS NULL
) AS username_exists;
-- name: PhoneOrEmailExists :one
SELECT EXISTS (
  SELECT 1 
  FROM users 
  WHERE (email = $1 OR phone = $1) AND deleted_at IS NULL
) AS email_or_phone_exists;
-- name: CreateUserToken :one
INSERT INTO user_tokens (
  token_id,
  user_id
)
VALUES (
  $1, $2
)
RETURNING *;
-- name: GetActiveUserTokenByUserID :one
SELECT * 
FROM user_tokens
WHERE user_id = $1 AND status = 'ACTIVE' AND deleted_at IS NULL;
-- name: ResetActiveToken :exec
UPDATE user_tokens
SET status = 'INACTIVE', updated_at = NOW()
WHERE user_id = $1 AND status = 'ACTIVE' AND deleted_at IS NULL;
-- name: UpdatePaymentIntentStatus :one
UPDATE payment_intents
SET status = $2, updated_at = NOW()
WHERE id = $1
RETURNING id, status, updated_at;

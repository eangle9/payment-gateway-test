-- name: CreatePaymentIntent :one
INSERT INTO payment_intents (
    company_id,
    customer_id,
    payment_type,
    amount,
    currency,
    callback_url,
    return_url,
    description,
    extra,
    status,
    bill_ref_no
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
)
RETURNING *;
-- name: GetPaymentIntentByID :one
SELECT
    pi.id,
    pi.payment_type,
    pi.amount,
    pi.status,
    pi.currency,
    pi.callback_url,
    pi.return_url,
    pi.description,
    pi.extra,
    pi.bill_ref_no,
    pi.expire_at,
    pi.created_at,
    pi.updated_at,
    json_build_object (
        'id',cu.id,
        'company_id',cu.company_id,
        'full_name',cu.full_name,
        'phone_number',cu.phone_number,
        'email',cu.email,
        'created_at',cu.created_at,
        'updated_at',cu.updated_at
    ) AS customer,
    json_build_object (
        'id',c.id,
        'name',c.name,
        'registration_number',c.registration_number,
        'address_street',c.address_street,
        'address_city',c.address_city,
        'address_state',c.address_state,
        'address_postal_code',c.address_postal_code,
        'address_country',c.address_country,
        'primary_phone',c.primary_phone,
        'secondary_phone',c.secondary_phone,
        'status',c.status,
        'email',c.email,
        'website',c.website,
        'callback_url',c.callback_url,
        'return_url',c.return_url,
        'created_at',c.created_at,
        'updated_at',c.updated_at
    ) AS company
FROM 
    payment_intents pi 
JOIN 
    customers cu ON pi.customer_id = cu.id 
JOIN 
    companies c ON pi.company_id = c.id
WHERE 
    pi.id = $1 AND pi.deleted_at IS NULL AND c.deleted_at IS NULL AND cu.deleted_at IS NULL;
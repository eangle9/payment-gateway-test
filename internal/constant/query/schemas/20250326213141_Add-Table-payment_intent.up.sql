------------------------------------------------
-- PaymentIntent Table
------------------------------------------------
CREATE TABLE IF NOT EXISTS payment_intents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    customer_id UUID NOT NULL,
    payment_type VARCHAR(100) NOT NULL,
    amount DECIMAL NOT NULL,
    currency VARCHAR(100) NOT NULL,
    callback_url VARCHAR(250) NOT NULL,
    return_url VARCHAR(250) NOT NULL,
    description TEXT NULL,
    extra JSON NULL,
    status VARCHAR(100) NOT NULL DEFAULT 'PENDING',
    bill_ref_no TEXT NULL,
    expire_at TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ NULL
);

ALTER TABLE payment_intents
    ADD CONSTRAINT fk_payment_intent_customer_id  FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE CASCADE;
ALTER TABLE payment_intents 
    ADD CONSTRAINT payment_intents_company_id FOREIGN KEY (company_id) REFERENCES companies(id) ON DELETE CASCADE;

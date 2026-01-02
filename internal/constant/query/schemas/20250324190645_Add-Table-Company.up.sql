-- Enable uuid-ossp extension if not already enabled
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

------------------------------------------------
-- Companies Table
------------------------------------------------
CREATE TABLE IF NOT EXISTS companies (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    registration_number VARCHAR(50) NOT NULL,  -- will be made unique
    address_street VARCHAR(255),
    address_city VARCHAR(100),
    address_state VARCHAR(100),
    address_postal_code VARCHAR(20),
    address_country VARCHAR(100),
    primary_phone VARCHAR(50),
    secondary_phone VARCHAR(50),
    status VARCHAR(100) NOT NULL DEFAULT 'ACTIVE',
    email VARCHAR(255),  -- will be made unique
    website VARCHAR(255),
    callback_url VARCHAR(255),
    return_url VARCHAR(255),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ NULL  -- for soft delete
);

-- Uniqueness constraints for companies
ALTER TABLE companies
    ADD CONSTRAINT uq_companies_registration_number UNIQUE (registration_number);
ALTER TABLE companies
    ADD CONSTRAINT uq_companies_email UNIQUE (email);

-- (Optional) Index on deleted_at if you expect many soft-deleted rows
CREATE INDEX idx_companies_not_deleted ON companies (id) WHERE deleted_at IS NULL;


------------------------------------------------
-- Customers Table
------------------------------------------------
CREATE TABLE IF NOT EXISTS customers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    company_id UUID NOT NULL,
    full_name TEXT NULL,
	phone_number TEXT NOT NULL,
	email TEXT NULL,
	status VARCHAR(100) NOT NULL DEFAULT 'ACTIVE',
	created_at TIMESTAMPTZ NOT NULL DEFAULT now()::TIMESTAMPTZ,
	updated_at TIMESTAMPTZ NOT NULL DEFAULT now()::TIMESTAMPTZ,
	deleted_at TIMESTAMPTZ NULL,  -- for soft delete

    CONSTRAINT customers_company_id_fkey FOREIGN KEY (company_id) REFERENCES companies(id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX unique_customer_company ON customers (company_id ASC, phone_number) WHERE deleted_at IS NULL;
CREATE INDEX idx_customers_company_phone_number ON customers (company_id ASC, phone_number) WHERE deleted_at IS NULL;
CREATE INDEX idx_customers_company ON customers (company_id) WHERE deleted_at IS NULL;

------------------------------------------------
-- Users Table (Example Definition)
------------------------------------------------
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    company_id UUID NOT NULL, 
    username VARCHAR(100),  
    email VARCHAR(255) NOT NULL,        -- will be unique
    phone VARCHAR(50) NOT NULL,
    password TEXT NOT NULL,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    role VARCHAR(50),
    status VARCHAR(50) NOT NULL DEFAULT 'ACTIVE',
    timezone_id VARCHAR(50),
    bio TEXT,
    profile_picture_url TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ NULL  -- for soft delete
);
ALTER TABLE users 
ADD CONSTRAINT users_email_or_phone_not_null 
CHECK (email IS NOT NULL OR phone IS NOT NULL);
ALTER TABLE users
    ADD CONSTRAINT uq_users_username UNIQUE (username);
ALTER TABLE users
    ADD CONSTRAINT uq_users_email UNIQUE (email);
ALTER TABLE users 
    ADD CONSTRAINT users_company_id_fk FOREIGN KEY (company_id) REFERENCES companies(id) ON DELETE CASCADE; 

CREATE INDEX idx_users_not_deleted ON users (id) WHERE deleted_at IS NULL;

------------------------------------------------
-- User Tokens Table
------------------------------------------------
CREATE TABLE IF NOT EXISTS user_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    token_id UUID NOT NULL,
    user_id UUID NOT NULL,
    status VARCHAR(255) NOT NULL DEFAULT 'ACTIVE',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ NULL  -- for soft delete
);

ALTER TABLE user_tokens 
    ADD CONSTRAINT fk_user_tokens_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

-- Only one active token per user (active tokens are those with status 'ACTIVE' and not soft-deleted)
CREATE UNIQUE INDEX idx_user_tokens_active_unique 
    ON user_tokens (user_id)
    WHERE status = 'ACTIVE' AND deleted_at IS NULL;

CREATE INDEX idx_user_tokens_user_id ON user_tokens (user_id) WHERE deleted_at IS NULL;

------------------------------------------------
-- Company Tokens Table
------------------------------------------------
CREATE TABLE IF NOT EXISTS company_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    token_id UUID NOT NULL DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    status VARCHAR(255) NOT NULL DEFAULT 'ACTIVE',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ NULL  -- for soft delete
);

ALTER TABLE company_tokens
    ADD CONSTRAINT fk_company_tokens_company FOREIGN KEY (company_id) REFERENCES companies(id) ON DELETE CASCADE;

-- Only one active token per company
CREATE UNIQUE INDEX idx_company_tokens_active_unique 
    ON company_tokens (company_id)
    WHERE status = 'ACTIVE' AND deleted_at IS NULL;

CREATE INDEX idx_company_tokens_company_id ON company_tokens (company_id) WHERE deleted_at IS NULL;

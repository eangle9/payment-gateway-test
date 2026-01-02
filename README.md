# Payment Gateway Test Project

This project is a payment gateway service that allows clients to create payments, process them asynchronously using RabbitMQ, and track their status reliably with PostgreSQL.

## Prerequisites

- Go 1.21+
- Docker and Docker Compose
- Make (optional, but recommended)

## Getting Started

### 1. Setup Environment Variables

Copy the example environment file and update the values if necessary:

```bash
cp .env.example .env
```

### 2. Start Infrastructure

Start the PostgreSQL database and RabbitMQ using Docker Compose:

```bash
make mg-up
```

### 3. Run the Application

Start the Go server and the background worker:

```bash
make dev
```

The server will be running on `http://localhost:8000`.

## Testing Flow

You can test the APIs using the built-in Swagger UI:
[http://localhost:8000/api/swagger/index.html](http://localhost:8000/api/swagger/index.html)

### Step-by-Step Testing Guide

#### 1. Register a Company
- **Endpoint**: `POST /api/signup-company-owner`
- **Payload**: Provide company and owner details.
- **Goal**: Create a company entity that will own the payment intents.

#### 2. Login
- **Endpoint**: `POST /api/login`
- **Payload**: Use the credentials from Step 1.
- **Response**: You will receive an `access_token`.
- **Action**: Copy the `access_token`.

#### 3. Generate Secret Token
- **Endpoint**: `POST /api/generate-secret-token`
- **Header**: `Authorization: Bearer <access_token>`
- **Goal**: Obtain a long-lived `secret_token` used for server-to-server payment operations.
- **Action**: Copy the `secret_token`.

#### 4. Create a Payment Intent
- **Endpoint**: `POST /api/payment-intents`
- **Header**: `Authorization: Bearer <secret_token>`
- **Payload**:
  ```json
  {
    "amount": "100.00",
    "currency": "ETB",
    "customer": {
      "full_name": "John Doe",
      "phone_number": "+251911223344",
      "email": "john@example.com"
    },
    "description": "Test Payment"
  }
  ```
- **Goal**: This will store the payment as `PENDING` and publish a message to RabbitMQ.

#### 5. Verify Asynchronous Processing
- Check the server logs. You should see the worker picking up the message and simulating processing.
- The payment status will be randomly updated to `SUCCESS` or `FAILED` after 2 seconds.

#### 6. Get Payment Intent Detail
- **Endpoint**: `GET /api/payment-intents/{id}`
- **Header**: `Authorization: Bearer <secret_token>`
- **Goal**: Verify the final status of the payment.

## Core Features Implemented

- **Idempotency**: Row-level locking (`SELECT ... FOR UPDATE`) ensures payments are never processed more than once.
- **Concurrency**: Multiple workers can safely process different payments concurrently.
- **Validation**: Strict input validation (e.g., Currency must be `ETB` or `USD`).
- **Centralized Error Handling**: Standardized JSON error responses.
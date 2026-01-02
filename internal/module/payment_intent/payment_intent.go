package paymentintent

import (
	"context"
	"encoding/json"
	"pg/initiator/platform/amqp"
	"pg/internal/constant"
	"pg/internal/constant/errors"
	"pg/internal/constant/model/db"
	"pg/internal/constant/model/dto"
	persistencedb "pg/internal/constant/persistenceDB"
	"pg/internal/module"
	"pg/internal/storage"
	"pg/platform/hlog"
	"pg/platform/httpclient"
	"pg/platform/utils"
	"strings"
	"time"

	"github.com/google/uuid"
	amqp091 "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type paymentIntent struct {
	log                  hlog.Logger
	paymentIntentStorage storage.PaymentIntent
	companyStorage       storage.Company
	httpClient           httpclient.HTTPClient
	amqpClient           amqp.Client
	persistenceDB        persistencedb.PersistenceDB
}

func New(paymentIntentStorage storage.PaymentIntent,
	log hlog.Logger,
	companyStorage storage.Company,
	httpClient httpclient.HTTPClient,
	amqpClient amqp.Client,
	persistenceDB persistencedb.PersistenceDB) module.PaymentIntent {
	return &paymentIntent{
		log:                  log,
		paymentIntentStorage: paymentIntentStorage,
		companyStorage:       companyStorage,
		httpClient:           httpClient,
		amqpClient:           amqpClient,
		persistenceDB:        persistenceDB,
	}
}

func (p *paymentIntent) InitPaymentIntent(ctx context.Context,
	param dto.InitPaymentIntent, companyID string) (*dto.PaymentIntent, error) {
	if err := param.Validate(); err != nil {
		err = errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		p.log.Warn(ctx, "invalid input", zap.Error(err))
		return nil, err
	}

	cID, err := uuid.Parse(companyID)
	if err != nil {
		err = errors.ErrInvalidUserInput.Wrap(err, "unable to parse company id")
		p.log.Error(ctx, "error parsing company id",
			zap.Error(err), zap.String("company-id", companyID))
		return nil, err
	}
	company, err := p.companyStorage.GetCompanyByID(ctx, cID)
	if err != nil {
		return nil, err
	}

	if param.CallBackURL == "" {
		param.CallBackURL = company.CallBackURL
	}
	if param.ReturnURL == "" {
		param.CallBackURL = company.ReturnURL
	}

	param.Currency = strings.ToUpper(param.Currency)

	billRefNO, err := utils.GenerateHash(uuid.NewString()+time.Now().String()+utils.CapitalLetters, 8)
	if err != nil {
		err = errors.ErrInternalServerError.Wrap(err, "unable to generate hash")
		p.log.Error(ctx, "unable to generate hash", zap.Error(err))
		return nil, err
	}

	// Parse phone number
	if param.Customer.PhoneNumber != "" {
		phone, err := utils.ParsePhoneNumber(param.Customer.PhoneNumber)
		if err != nil {
			err = errors.ErrInvalidUserInput.Wrap(err, "failed to parse phone number")
			p.log.Error(ctx, "failed to parse phone number",
				zap.Error(err),
				zap.String("phone", param.Customer.PhoneNumber),
			)
			return nil, err
		}
		param.Customer.PhoneNumber = *phone
	}

	paymentIntent, err := p.paymentIntentStorage.CreatePaymentIntent(ctx,
		dto.CreatePaymentIntent{
			CompanyID:   company.ID,
			PaymentType: constant.PaymentTypeOnetime,
			Amount:      param.Amount,
			Status:      constant.Pending,
			Currency:    constant.Currency(param.Currency),
			CallBackURL: param.CallBackURL,
			ReturnURL:   param.ReturnURL,
			Description: param.Description,
			Customer:    param.Customer,
			Extra:       param.Extra,
			BillRefNO:   billRefNO,
		}, p.amqpClient)
	if err != nil {
		return nil, err
	}

	return paymentIntent, nil
}

func (p *paymentIntent) GetPaymentIntentDetail(ctx context.Context,
	id string) (*dto.PaymentIntentDetail, error) {
	pID, err := uuid.Parse(id)
	if err != nil {
		err = errors.ErrInternalServerError.Wrap(err, "unable to parse payment intent id")
		p.log.Error(ctx, "error parsing payment intent id",
			zap.Error(err), zap.String("payment-intent-id", id))
		return nil, err
	}

	return p.paymentIntentStorage.GetPaymentIntentByID(ctx, pID)
}

func (p *paymentIntent) StartWorker(ctx context.Context) {
	// Consume messages
	ch, err := p.amqpClient.GetConnection().Channel()
	if err != nil {
		p.log.Fatal(ctx, "failed to open channel", zap.Error(err))
	}
	defer ch.Close()

	// Declare the queue
	_, err = ch.QueueDeclare(
		"payment_processing", // name
		true,                 // durable
		false,                // delete when unused
		false,                // exclusive
		false,                // no-wait
		nil,                  // arguments
	)
	if err != nil {
		p.log.Fatal(ctx, "failed to declare queue", zap.Error(err))
	}

	msgs, err := ch.Consume(
		"payment_processing", // queue
		"",                   // consumer
		true,                 // auto-ack
		false,                // exclusive
		false,                // no-local
		false,                // no-wait
		nil,                  // args
	)
	if err != nil {
		p.log.Fatal(ctx, "failed to register a consumer", zap.Error(err))
	}

	p.log.Info(ctx, " [*] Waiting for messages. To exit press CTRL+C")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			p.log.Info(ctx, "received a message", zap.String("body", string(d.Body)))
			p.processPayment(ctx, d)
		}
	}()

	<-forever
}

func (p *paymentIntent) processPayment(ctx context.Context, d amqp091.Delivery) {
	var payload map[string]string
	if err := json.Unmarshal(d.Body, &payload); err != nil {
		p.log.Error(ctx, "failed to unmarshal message", zap.Error(err))
		return
	}

	paymentIntentID := payload["payment_intent_id"]
	if paymentIntentID == "" {
		p.log.Error(ctx, "payment_intent_id is empty")
		return
	}

	pID, err := uuid.Parse(paymentIntentID)
	if err != nil {
		p.log.Error(ctx, "invalid payment_intent_id", zap.Error(err))
		return
	}

	err = p.persistenceDB.WithTransaction(ctx, func(tx persistencedb.PersistenceDB) error {
		// 1. Lock the row and get current status
		pi, err := tx.GetPaymentIntentByIDForUpdate(ctx, pID)
		if err != nil {
			return err
		}

		// 2. Check status (Idempotency)
		if constant.Status(pi.Status) != constant.Pending {
			p.log.Info(ctx, "payment already processed or in progress",
				zap.String("id", paymentIntentID), zap.String("status", string(pi.Status)))
			return nil
		}

		// 3. Simulate processing
		p.log.Info(ctx, "processing payment", zap.String("id", paymentIntentID))
		time.Sleep(2 * time.Second)

		// 4. Randomly succeed or fail
		status := constant.Success
		if time.Now().Unix()%2 == 0 {
			status = constant.Failed
		}

		// 5. Update status
		_, err = tx.UpdatePaymentIntentStatus(ctx, db.UpdatePaymentIntentStatusParams{
			ID:     pID,
			Status: string(status),
		})
		if err != nil {
			return err
		}

		p.log.Info(ctx, "payment processed", zap.String("id", paymentIntentID), zap.String("status", string(status)))
		return nil
	})

	if err != nil {
		p.log.Error(ctx, "failed to process payment", zap.Error(err), zap.String("id", paymentIntentID))
	}
}

package amqp

import (
	"context"
	"fmt"
	"pg/platform/hlog"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Client interface {
	Publish(ctx context.Context, exchange, routingKey string, body []byte) error
	Close() error
	GetConnection() *amqp.Connection
}

type client struct {
	conn *amqp.Connection
	ch   *amqp.Channel
	log  hlog.Logger
}

func New(url string, log hlog.Logger) (Client, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to rabbitmq: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	return &client{
		conn: conn,
		ch:   ch,
		log:  log,
	}, nil
}

func (c *client) Publish(ctx context.Context, exchange, routingKey string, body []byte) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return c.ch.PublishWithContext(ctx,
		exchange,   // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
}

func (c *client) Close() error {
	if err := c.ch.Close(); err != nil {
		return err
	}
	return c.conn.Close()
}

func (c *client) GetConnection() *amqp.Connection {
	return c.conn
}

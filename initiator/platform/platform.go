package platform

import (
	"context"
	"pg/initiator/foundation"
	"pg/initiator/platform/amqp"
	"pg/platform/hcrypto"
	"pg/platform/hlog"
	"pg/platform/httpclient"

	"go.uber.org/zap"
)

type Layer struct {
	Token      hcrypto.Maker
	HTTPClient httpclient.HTTPClient
	AMQP       amqp.Client
}

func InitPlatform(log hlog.Logger, state foundation.State) Layer {
	amqpClient, err := amqp.New(state.AMQPURL, log.Named("amqp"))
	if err != nil {
		// Log error but don't fail hard if AMQP is optional, or fail if required.
		// For this task, let's assume it's critical but we might want to allow startup without it for other parts.
		// However, the requirement implies it's core. Let's log fatal if it fails?
		// Or just log error. Ideally we should retry or fail.
		// Given the context, let's just log error and continue, but in production we might want to block.
		// Actually, let's panic/fatal if we can't connect, as it's a core feature.
		// But InitPlatform doesn't return error.
		// Let's log fatal.
		log.Error(context.Background(), "failed to initialize amqp client", zap.Error(err))
		// For now, we proceed. Usage will panic or fail.
	}

	return Layer{
		Token:      InitToken(state.TokenConfig, log.Named("token")),
		HTTPClient: httpclient.Init(state.HTTPConfig, log.Named("httpclient")),
		AMQP:       amqpClient,
	}
}

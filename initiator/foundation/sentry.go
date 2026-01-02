package foundation

import (
	"github.com/getsentry/sentry-go"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func InitSentry(log *zap.Logger) {
	log.Info("initializing Sentry")

	if err := sentry.Init(sentry.ClientOptions{
		Dsn:              viper.GetString("SENTRY_DSN"),
		EnableTracing:    true,
		TracesSampleRate: viper.GetFloat64("SENTRY_TRACES_SAMPLE_RATE"),
		Environment:      viper.GetString("ENVIRONMENT"),
	}); err != nil {
		log.Fatal("failed to initialize Sentry", zap.Error(err))
	}

	log.Info("Sentry initialized")
}

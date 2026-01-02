package foundation

import (
	"context"
	"log"
	"pg/internal/constant/errors"
	"pg/platform/hcrypto"
	"pg/platform/hlog"
	"pg/platform/httpclient"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type State struct {
	HTTPConfig  httpclient.HTTPTransport
	TokenConfig hcrypto.TokenKey
	AMQPURL     string
}

func InitState(logger hlog.Logger) State {
	httpconfig := httpclient.HTTPTransport{
		MaxIdleConnsPerHost: viper.GetInt("HTTP_MAX_IDLE_CONNS_PER_HOST"),
		MaxIdleConns:        viper.GetInt("HTTP_MAX_IDLE_CONNS"),
		MaxConnsPerHost:     viper.GetInt("HTTP_MAX_CONNS_PER_HOST"),
		Timeout:             viper.GetDuration("HTTP_TIMEOUT"),
	}

	if err := httpconfig.Validate(); err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid http configuration")
		log.Fatal(context.Background(), "all http fields are required", zap.Error(err))
	}

	tokenConfig := hcrypto.TokenKey{
		SymmetricKey: viper.GetString("SECURITY_CREDENTIAL_SYMMETRIC_KEY"),
		Issuer:       viper.GetString("SECURITY_CREDENTIAL_ISSUER"),
		Footer:       viper.GetString("SECURITY_CREDENTIAL_FOOTER"),
		KeyLength:    viper.GetInt("SECURITY_CREDENTIAL_KEY_LENGTH"),
		Audience:     viper.GetString("SECURITY_CREDENTIAL_AUDIENCE"),
		AccessExpires: time.Duration(
			viper.GetInt("SECURITY_CREDENTIAL_ACCESS_EXPIRES")) * time.Minute,
		RefreshExpires: time.Duration(
			viper.GetInt("SECURITY_CREDENTIAL_REFRESH_EXPIRES")) * time.Minute,
		SecretTokenExpires: time.Duration(
			viper.GetInt("SECURITY_CREDENTIAL_SECRET_EXPIRES")) * time.Minute,
	}

	if err := tokenConfig.Validate(); err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid token config configuration")
		log.Fatal(context.Background(), "all tokenKey fields are required", zap.Error(err))
	}

	return State{
		HTTPConfig:  httpconfig,
		TokenConfig: tokenConfig,
		AMQPURL:     viper.GetString("RABBITMQ_URL"),
	}
}

package platform

import (
	"pg/initiator/foundation"
	"pg/platform/hcrypto"
	"pg/platform/hlog"
	"pg/platform/httpclient"
)

type Layer struct {
	Token      hcrypto.Maker
	HTTPClient httpclient.HTTPClient
}

func InitPlatform(log hlog.Logger, state foundation.State) Layer {
	return Layer{
		Token:      InitToken(state.TokenConfig, log.Named("token")),
		HTTPClient: httpclient.Init(state.HTTPConfig, log.Named("httpclient")),
	}
}

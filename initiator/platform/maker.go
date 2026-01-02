package platform

import (
	"pg/platform/hcrypto"
	"pg/platform/hlog"
)

func InitToken(tokenconfig hcrypto.TokenKey, log hlog.Logger) hcrypto.Maker {
	return hcrypto.PasetoInit(tokenconfig, log.Named("token-platform"))
}

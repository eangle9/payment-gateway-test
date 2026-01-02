package platform

import (
	"pg/platform/hlog"
	"pg/platform/httpclient"
)

func InitHTTPClient(httpConfig httpclient.HTTPTransport, log hlog.Logger) httpclient.HTTPClient {
	return httpclient.Init(httpConfig, log)
}

package client

import (
	_http "net/http"
	"net/url"

	"github.com/go-logr/logr"

	_client "github.com/ARM-software/embedded-development-services-client/client"
	"github.com/ARM-software/golang-utils/utils/commonerrors"
	"github.com/ARM-software/golang-utils/utils/http"
)

// NewClient returns a new API client based on request configuration
func NewClient(cfg *http.RequestConfiguration, logger logr.Logger, underlyingHTTPClient *_http.Client) (c *_client.APIClient, err error) {
	if logger.IsZero() {
		err = commonerrors.ErrNoLogger
		return
	}
	if cfg == nil {
		err = commonerrors.UndefinedVariable("configuration")
		return
	}
	err = cfg.Validate()
	if err != nil {
		err = commonerrors.WrapError(commonerrors.ErrInvalid, err, "invalid client configuration")
		return
	}

	httpClientCfg := http.DefaultRobustHTTPClientConfiguration()
	httpClient := http.NewConfigurableRetryableOauthClientWithLoggerAndCustomClient(httpClientCfg, underlyingHTTPClient, logger, cfg.Authorisation.AccessToken)
	clientCfg := newClientConfiguration(cfg)
	clientCfg.HTTPClient = httpClient.StandardClient()
	c = _client.NewAPIClient(clientCfg)
	return
}

func newClientConfiguration(cfg *http.RequestConfiguration) (clientCfg *_client.Configuration) {
	clientCfg = _client.NewConfiguration()
	if cfg.Target.Host != "" {
		basePathURL, err := url.Parse(cfg.Target.Host)
		if err == nil {
			clientCfg.Host = basePathURL.Host
			clientCfg.Scheme = basePathURL.Scheme
		}
	}
	clientCfg.UserAgent = cfg.UserAgent

	return
}

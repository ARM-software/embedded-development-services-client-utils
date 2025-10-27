/*
 * Copyright (C) 2020-2025 Arm Limited or its affiliates and Contributors. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 */
package client

import (
	_http "net/http"
	"net/url"

	"github.com/go-logr/logr"

	_client "github.com/ARM-software/embedded-development-services-client/client"
	"github.com/ARM-software/golang-utils/utils/commonerrors"
	"github.com/ARM-software/golang-utils/utils/field"
	"github.com/ARM-software/golang-utils/utils/http"
	"github.com/ARM-software/golang-utils/utils/reflection"
)

// NewClient returns a new API client based on request configuration
func NewClient(cfg *http.RequestConfiguration, logger logr.Logger, underlyingHTTPClient *_http.Client) (c *_client.APIClient, err error) {
	httpClient, err := NewHTTPClient(cfg, http.DefaultBasicRetryPolicyConfiguration(), logger, underlyingHTTPClient)
	if err != nil {
		return
	}
	clientCfg := newClientConfiguration(cfg)
	clientCfg.HTTPClient = httpClient.StandardClient()
	c = _client.NewAPIClient(clientCfg)
	return
}

// NewHTTPClient returns an HTTP retryable client based on request configuration
func NewHTTPClient(cfg *http.RequestConfiguration, retryCfg *http.RetryPolicyConfiguration, logger logr.Logger, underlyingHTTPClient *_http.Client) (httpClient http.IRetryableClient, err error) {
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

	httpClientCfg := http.DefaultHTTPClientConfiguration()
	httpClientCfg.RetryPolicy = field.Optional[http.RetryPolicyConfiguration](retryCfg, *http.DefaultBasicRetryPolicyConfiguration())
	httpClient = http.NewConfigurableRetryableClientWithLoggerAndCustomClient(httpClientCfg, cfg, logger, underlyingHTTPClient)
	return
}

func newClientConfiguration(cfg *http.RequestConfiguration) (clientCfg *_client.Configuration) {
	clientCfg = _client.NewConfiguration()
	if !reflection.IsEmpty(cfg.Host) {
		basePathURL, err := url.Parse(cfg.Host)
		if err == nil {
			clientCfg.Host = basePathURL.Host
			clientCfg.Scheme = basePathURL.Scheme
		}
	}
	clientCfg.UserAgent = cfg.UserAgent

	return
}

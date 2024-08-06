package httpplus

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/hashicorp/go-retryablehttp"
)

// Config is the configuration for http client
type Config struct {
	Component         string `mapstructure:"component" json:"component"`
	RetryMax          int    `mapstructure:"retry_max" json:"retry_max"`
	RetryWaitMaxSec   int    `mapstructure:"retry_wait_max_sec" json:"retry_wait_max_sec"`
	RetryWaitMinSec   int    `mapsctructure:"retry_wait_min_sec" json:"retry_wait_min_sec"`
	LogRequest        bool   `mapstructure:"log_request" json:"log_request"`
	LogRequestDetail  bool   `mapstructure:"log_request_detail" json:"log_request_detail"`
	LogResponse       bool   `mapstructure:"log_response" json:"log_response"`
	LogResponseDetail bool   `mapstructure:"log_response_detail" json:"log_response_detail"`
}

// Parse parses the configuration
func (c *Config) Parse() {
	if c.Component == "" {
		c.Component = "http client"
	}
	if c.RetryMax == 0 {
		c.RetryMax = 3
	}
	if c.RetryWaitMaxSec == 0 {
		c.RetryWaitMaxSec = 3
	}
	if c.RetryWaitMinSec == 0 {
		c.RetryWaitMinSec = 1
	}
}

// NewHttpClient creates a new http client
func NewHttpClient(cfg *Config) *retryablehttp.Client {
	if cfg == nil {
		cfg = &Config{}
	}
	cfg.Parse()

	client := retryablehttp.NewClient()
	client.RetryWaitMin = time.Duration(cfg.RetryWaitMinSec) * time.Second
	client.RetryWaitMax = time.Duration(cfg.RetryWaitMaxSec) * time.Second
	client.RetryMax = cfg.RetryMax
	client.Logger = slog.With("component", cfg.Component)

	client.RequestLogHook = func(_ retryablehttp.Logger, req *http.Request, retryNumber int) {
		if cfg.LogRequest {
			slog.Debug("request", "retry", retryNumber, "info", RequestInfo(req, cfg.LogRequestDetail))
		}
	}

	client.ResponseLogHook = func(_ retryablehttp.Logger, resp *http.Response) {
		if resp.StatusCode >= 400 {
			slog.Warn("response", "info", ResponseInfo(resp, true))
		} else if cfg.LogResponse {
			slog.Debug("response", "info", ResponseInfo(resp, cfg.LogResponseDetail))
		}
	}
	return client
}

// RequestInfo returns request info
func RequestInfo(req *http.Request, detail bool) map[string]any {
	r := map[string]any{
		"method": req.Method,
		"url":    req.URL.String(),
		"query":  req.URL.Query(),
	}
	if detail {
		headers := map[string][]string{}
		for k, v := range req.Header {
			if k == "Authorization" {
				headers[k] = []string{"***"}
			} else {
				headers[k] = v
			}
		}
		r["header"] = headers
	}
	return r
}

// ResponseInfo returns response info
func ResponseInfo(resp *http.Response, detail bool) map[string]any {
	r := map[string]any{
		"status":  resp.Status,
		"code":    resp.StatusCode,
		"request": RequestInfo(resp.Request, false),
	}
	if detail {
		r["header"] = resp.Header
		r["body"] = resp.Body
	}
	return r
}

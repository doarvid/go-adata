package indexmarket

import (
	"time"

	"github.com/go-resty/resty/v2"
)

type IndexMarketConfig struct {
	Timeout   time.Duration
	Proxy     string
	UserAgent string
	Headers   map[string]string
	Client    *resty.Client
}

type IndexMarketOption func(*IndexMarketConfig)

func WithTimeout(d time.Duration) IndexMarketOption {
	return func(cfg *IndexMarketConfig) { cfg.Timeout = d }
}

func WithProxy(p string) IndexMarketOption {
	return func(cfg *IndexMarketConfig) { cfg.Proxy = p }
}

func WithUserAgent(ua string) IndexMarketOption {
	return func(cfg *IndexMarketConfig) { cfg.UserAgent = ua }
}

func WithHeaders(h map[string]string) IndexMarketOption {
	return func(cfg *IndexMarketConfig) { cfg.Headers = h }
}

func WithClient(c *resty.Client) IndexMarketOption {
	return func(cfg *IndexMarketConfig) { cfg.Client = c }
}

type IndexMarket struct {
	client *resty.Client
	cfg    IndexMarketConfig
}

func NewIndexMarket(opts ...IndexMarketOption) *IndexMarket {
	cfg := IndexMarketConfig{
		Timeout:   15 * time.Second,
		UserAgent: "go-adata/indexmarket",
		Headers:   map[string]string{},
	}
	for _, opt := range opts {
		opt(&cfg)
	}
	var c *resty.Client
	if cfg.Client != nil {
		c = cfg.Client
	} else {
		c = resty.New()
		c.SetTimeout(cfg.Timeout)
		if cfg.UserAgent != "" {
			c.SetHeader("User-Agent", cfg.UserAgent)
		}
		if cfg.Proxy != "" {
			c.SetProxy(cfg.Proxy)
		}
	}
	return &IndexMarket{client: c, cfg: cfg}
}


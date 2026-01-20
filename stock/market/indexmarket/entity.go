package indexmarket

import (
	"time"

	browser "github.com/EDDYCJY/fake-useragent"
	"github.com/doarvid/go-adata/common/utils"
	"github.com/go-resty/resty/v2"
)

type IndexMarketConfig struct {
	Timeout   time.Duration
	Proxy     string
	UserAgent string
	Headers   map[string]string
	Client    *resty.Client
	Wait      time.Duration
	Retries   int
	Debug     bool
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
func WithWait(d time.Duration) IndexMarketOption {
	return func(cfg *IndexMarketConfig) { cfg.Wait = d }
}
func WithRetries(n int) IndexMarketOption {
	return func(cfg *IndexMarketConfig) { cfg.Retries = n }
}
func WithDebug(enable bool) IndexMarketOption {
	return func(cfg *IndexMarketConfig) { cfg.Debug = enable }
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
		Wait:      50 * time.Millisecond,
		Retries:   2,
		Debug:     false,
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
		ua := cfg.UserAgent
		if ua == "" {
			ua = browser.Random()
		}
		c.SetHeader("User-Agent", ua)
		if cfg.Proxy != "" {
			utils.ApplyProxyResty(c, cfg.Proxy)
		}
	}
	if cfg.Debug {
		c.SetDebug(true)
	}
	return &IndexMarket{client: c, cfg: cfg}
}

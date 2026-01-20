package alist

import (
	"time"

	browser "github.com/EDDYCJY/fake-useragent"
	"github.com/doarvid/go-adata/common/utils"
	"github.com/go-resty/resty/v2"
)

type Config struct {
	Timeout   time.Duration
	Proxy     string
	UserAgent string
	Headers   map[string]string
	Client    *resty.Client
	Wait      time.Duration
	Retries   int
	Debug     bool
}

type Option func(*Config)

func WithProxy(url string) Option            { return func(cfg *Config) { cfg.Proxy = url } }
func WithTimeout(d time.Duration) Option     { return func(cfg *Config) { cfg.Timeout = d } }
func WithUserAgent(ua string) Option         { return func(cfg *Config) { cfg.UserAgent = ua } }
func WithHeaders(h map[string]string) Option { return func(cfg *Config) { cfg.Headers = h } }
func WithClient(c *resty.Client) Option      { return func(cfg *Config) { cfg.Client = c } }
func WithWait(d time.Duration) Option        { return func(cfg *Config) { cfg.Wait = d } }
func WithRetries(n int) Option               { return func(cfg *Config) { cfg.Retries = n } }
func WithDebug(enable bool) Option           { return func(cfg *Config) { cfg.Debug = enable } }

type Client struct {
	client *resty.Client
	cfg    Config
}

func New(opts ...Option) *Client {
	cfg := Config{
		Timeout:   15 * time.Second,
		UserAgent: "go-adata/alist",
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
	return &Client{client: c, cfg: cfg}
}

package alist

import (
	"time"

	"github.com/go-resty/resty/v2"
)

type Config struct {
	Timeout   time.Duration
	Proxy     string
	UserAgent string
	Headers   map[string]string
	Client    *resty.Client
}

type Option func(*Config)

func WithProxy(url string) Option      { return func(cfg *Config) { cfg.Proxy = url } }
func WithTimeout(d time.Duration) Option { return func(cfg *Config) { cfg.Timeout = d } }
func WithUserAgent(ua string) Option  { return func(cfg *Config) { cfg.UserAgent = ua } }
func WithHeaders(h map[string]string) Option { return func(cfg *Config) { cfg.Headers = h } }
func WithClient(c *resty.Client) Option      { return func(cfg *Config) { cfg.Client = c } }

type Client struct {
	client *resty.Client
	cfg    Config
}

func New(opts ...Option) *Client {
	cfg := Config{
		Timeout:   15 * time.Second,
		UserAgent: "go-adata/alist",
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
	return &Client{client: c, cfg: cfg}
}


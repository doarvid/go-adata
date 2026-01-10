package conceptflow

import (
	"time"

	"github.com/go-resty/resty/v2"
)

type Option func(*Client)

type Client struct {
	client  *resty.Client
	Wait    time.Duration
	Retries int
}

func WithTimeout(d time.Duration) Option {
	return func(c *Client) { c.client.SetTimeout(d) }
}

func WithProxy(p string) Option {
	return func(c *Client) { c.client.SetProxy(p) }
}

func WithUserAgent(ua string) Option {
	return func(c *Client) { c.client.SetHeader("User-Agent", ua) }
}

func WithClient(rc *resty.Client) Option {
	return func(c *Client) { c.client = rc }
}

func WithWait(d time.Duration) Option {
	return func(c *Client) { c.Wait = d }
}

func WithRetries(n int) Option {
	return func(c *Client) { c.Retries = n }
}

func New(opts ...Option) *Client {
	c := &Client{
		client:  resty.New(),
		Wait:    50 * time.Millisecond,
		Retries: 2,
	}
	c.client.SetHeader("User-Agent", "go-adata/conceptflow")
	for _, opt := range opts {
		opt(c)
	}
	return c
}


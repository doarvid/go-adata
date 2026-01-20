package margin

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	browser "github.com/EDDYCJY/fake-useragent"
	"github.com/doarvid/go-adata/common/utils"
	"github.com/go-resty/resty/v2"
)

type Row struct {
	TradeDate string  `json:"trade_date"`
	Rzye      float64 `json:"rzye"`
	Rqye      float64 `json:"rqye"`
	Rzrqye    float64 `json:"rzrqye"`
	Rzrqyecz  float64 `json:"rzrqyecz"`
}

type Config struct {
	Timeout   time.Duration
	Proxy     string
	UserAgent string
	Client    *resty.Client
	Wait      time.Duration
	Retries   int
	Debug     bool
}
type Option func(*Config)

func WithTimeout(d time.Duration) Option { return func(cfg *Config) { cfg.Timeout = d } }
func WithProxy(p string) Option          { return func(cfg *Config) { cfg.Proxy = p } }
func WithUserAgent(ua string) Option     { return func(cfg *Config) { cfg.UserAgent = ua } }
func WithClient(c *resty.Client) Option  { return func(cfg *Config) { cfg.Client = c } }
func WithWait(d time.Duration) Option    { return func(cfg *Config) { cfg.Wait = d } }
func WithRetries(n int) Option           { return func(cfg *Config) { cfg.Retries = n } }
func WithDebug(enable bool) Option       { return func(cfg *Config) { cfg.Debug = enable } }

type Client struct {
	client *resty.Client
	cfg    Config
}

func New(opts ...Option) *Client {
	cfg := Config{
		Timeout:   15 * time.Second,
		UserAgent: "go-adata/margin",
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

func (m *Client) History(ctx context.Context, startDate string) ([]Row, error) {
	totalPages := 1
	currPage := 1
	pageSize := 250
	startDateStr := startDate
	var start time.Time
	hasStart := false
	if startDate != "" {
		t, err := time.Parse("2006-01-02", startDate)
		if err == nil {
			start = t
			hasStart = true
		}
	}
	out := make([]Row, 0, 512)
	for currPage <= totalPages {
		url := "https://datacenter-web.eastmoney.com/api/data/v1/get?reportName=RPTA_RZRQ_LSHJ&columns=ALL&source=WEB&sortColumns=dim_date&sortTypes=-1&pageNumber=" + strconv.Itoa(currPage) + "&pageSize=" + strconv.Itoa(pageSize)
		if m.cfg.Wait > 0 {
			time.Sleep(m.cfg.Wait)
		}
		resp, err := m.client.R().SetContext(ctx).Get(url)
		if err != nil {
			break
		}
		var res struct {
			Success bool `json:"success"`
			Result  struct {
				Pages int              `json:"pages"`
				Data  []map[string]any `json:"data"`
			} `json:"result"`
		}
		if err := json.Unmarshal(resp.Body(), &res); err != nil {
			break
		}
		if !res.Success {
			break
		}
		if currPage == 1 {
			totalPages = res.Result.Pages
		}
		data := res.Result.Data
		for _, it := range data {
			dt, _ := time.Parse("2006-01-02 15:04:05", toString(it["DIM_DATE"]))
			out = append(out, Row{
				TradeDate: dt.Format("2006-01-02"),
				Rzye:      parseF(toString(it["RZYE"])),
				Rqye:      parseF(toString(it["RQYE"])),
				Rzrqye:    parseF(toString(it["RZRQYE"])),
				Rzrqyecz:  parseF(toString(it["RZRQYECZ"])),
			})
		}
		if !hasStart {
			break
		}
		if hasStart {
			last := data[len(data)-1]
			dmin, _ := time.Parse("2006-01-02 15:04:05", toString(last["DIM_DATE"]))
			if !dmin.Before(start) {
				break
			}
		}
		currPage++
	}
	if startDateStr != "" {
		out2 := make([]Row, 0, len(out))
		for _, r := range out {
			if r.TradeDate > startDateStr {
				out2 = append(out2, r)
			}
		}
		out = out2
	}
	return out, nil
}

func parseF(s string) float64 {
	s = strings.TrimSpace(strings.ReplaceAll(s, "%", ""))
	if s == "" || s == "--" {
		return 0
	}
	v, _ := strconv.ParseFloat(s, 64)
	return v
}
func toString(v any) string { return strings.TrimSpace(fmt.Sprintf("%v", v)) }

package mine

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
	StockCode string  `json:"stock_code"`
	ShortName string  `json:"short_name"`
	Score     float64 `json:"score"`
	FType     string  `json:"f_type"`
	SType     string  `json:"s_type"`
	TType     string  `json:"t_type"`
	Reason    string  `json:"reason"`
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
		UserAgent: "go-adata/mine",
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

func (mc *Client) EvaluateTDX(ctx context.Context, stockCode string) ([]Row, error) {
	if stockCode == "" {
		return nil, fmt.Errorf("stock code is empty")
	}
	url := "http://page3.tdx.com.cn:7615/site/pcwebcall_static/bxb/json/" + stockCode + ".json"
	if mc.cfg.Wait > 0 {
		time.Sleep(mc.cfg.Wait)
	}
	resp, err := mc.client.R().SetContext(ctx).Get(url)
	if err != nil {
		return []Row{{StockCode: "", ShortName: "", Score: 0, FType: "暂无数据"}}, nil
	}
	var res struct {
		Name string           `json:"name"`
		Data []map[string]any `json:"data"`
	}
	if err := json.Unmarshal(resp.Body(), &res); err != nil {
		return []Row{{StockCode: "", ShortName: "", Score: 0, FType: "暂无数据"}}, nil
	}
	name := res.Name
	data := res.Data
	out := make([]Row, 0, 128)
	score := 100.0
	sTypeDeduct := map[string]bool{}
	for _, i := range data {
		ftype := toString(i["name"])
		rows, _ := i["rows"].([]any)
		for _, k := range rows {
			kk, _ := k.(map[string]any)
			if toString(kk["trigyy"]) != "" {
				com, _ := kk["commonlxid"].([]any)
				if len(com) == 0 {
					out = append(out, Row{StockCode: stockCode, ShortName: name, FType: ftype, SType: toString(kk["lx"]), TType: "", Reason: toString(kk["trigyy"]), Score: parseF(toString(kk["fs"]))})
					if toInt(kk["trig"]) == 1 {
						score -= parseF(toString(kk["fs"]))
					}
				}
				for _, j := range com {
					jj, _ := j.(map[string]any)
					if toString(jj["trigyy"]) != "" {
						out = append(out, Row{StockCode: stockCode, ShortName: name, FType: ftype, SType: toString(kk["lx"]), TType: toString(jj["lx"]), Reason: toString(jj["trigyy"]), Score: parseF(toString(jj["fs"]))})
						if toInt(jj["trig"]) == 1 && !sTypeDeduct[toString(kk["lx"])] {
							score -= parseF(toString(jj["fs"]))
							sTypeDeduct[toString(kk["lx"])] = true
						}
					}
				}
			}
		}
	}
	if len(out) == 0 {
		if strings.HasSuffix(name, "退") {
			return []Row{{StockCode: stockCode, ShortName: name, Score: -1, FType: "已退市"}}, nil
		}
		if score < 1 {
			score = 1
		}
		return []Row{{StockCode: stockCode, ShortName: name, Score: score, FType: "暂无风险项"}}, nil
	}
	if score < 1 {
		score = 1
	}
	for i := range out {
		out[i].Score = score
	}
	return out, nil
}

func toInt(v any) int {
	s := strings.TrimSpace(fmt.Sprintf("%v", v))
	if s == "" {
		return 0
	}
	n, _ := strconv.Atoi(s)
	return n
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

package northflow

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/doarvid/go-adata/stock/info/tradecalendar"
	"github.com/go-resty/resty/v2"
)

type Daily struct {
	TradeDate string  `json:"trade_date"`
	NetHgt    float64 `json:"net_hgt"`
	BuyHgt    float64 `json:"buy_hgt"`
	SellHgt   float64 `json:"sell_hgt"`
	NetSgt    float64 `json:"net_sgt"`
	BuySgt    float64 `json:"buy_sgt"`
	SellSgt   float64 `json:"sell_sgt"`
	NetTgt    float64 `json:"net_tgt"`
	BuyTgt    float64 `json:"buy_tgt"`
	SellTgt   float64 `json:"sell_tgt"`
}

type Minute struct {
	TradeTime string  `json:"trade_time"`
	NetHgt    float64 `json:"net_hgt"`
	NetSgt    float64 `json:"net_sgt"`
	NetTgt    float64 `json:"net_tgt"`
}

type Config struct {
	Timeout   time.Duration
	Proxy     string
	UserAgent string
	Client    *resty.Client
	Wait      time.Duration
	Retries   int
}
type Option func(*Config)

func WithTimeout(d time.Duration) Option { return func(cfg *Config) { cfg.Timeout = d } }
func WithProxy(p string) Option          { return func(cfg *Config) { cfg.Proxy = p } }
func WithUserAgent(ua string) Option     { return func(cfg *Config) { cfg.UserAgent = ua } }
func WithClient(c *resty.Client) Option  { return func(cfg *Config) { cfg.Client = c } }
func WithWait(d time.Duration) Option    { return func(cfg *Config) { cfg.Wait = d } }
func WithRetries(n int) Option           { return func(cfg *Config) { cfg.Retries = n } }

type Client struct {
	client *resty.Client
	cfg    Config
}

func New(opts ...Option) *Client {
	cfg := Config{
		Timeout:   15 * time.Second,
		UserAgent: "go-adata/northflow",
		Wait:      50 * time.Millisecond,
		Retries:   2,
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

func (c *Client) History(ctx context.Context, startDate string) ([]Daily, error) {
	currPage := 1
	out := make([]Daily, 0, 1024)
	var start time.Time
	var hasStart bool
	if startDate != "" {
		t, err := time.Parse("2006-01-02", startDate)
		if err == nil {
			start = t
			hasStart = true
		}
		dateMin, _ := time.Parse("2006-01-02", "2017-01-01")
		if start.Before(dateMin) {
			start = dateMin
		}
	}
	for currPage < 18 {
		base := "https://datacenter-web.eastmoney.com/api/data/v1/get?sortColumns=TRADE_DATE&sortTypes=-1&pageSize=1000&pageNumber=" + toString(currPage) + "&reportName=RPT_MUTUAL_DEAL_HISTORY&columns=ALL&source=WEB&client=WEB&"
		sgtURL := base + "filter=(MUTUAL_TYPE=%27001%27)"
		hgtURL := base + "filter=(MUTUAL_TYPE=%27003%27)"
		if c.cfg.Wait > 0 {
			time.Sleep(c.cfg.Wait)
		}
		resp1, err1 := c.client.R().SetContext(ctx).Get(sgtURL)
		if err1 != nil {
			break
		}
		resp2, err2 := c.client.R().SetContext(ctx).Get(hgtURL)
		if err2 != nil {
			break
		}
		buf1 := new(strings.Builder)
		buf2 := new(strings.Builder)
		if _, err := io.Copy(buf1, strings.NewReader(resp1.String())); err != nil {
			break
		}
		if _, err := io.Copy(buf2, strings.NewReader(resp2.String())); err != nil {
			break
		}
		sgtText := strings.ReplaceAll(buf1.String(), "null", "0")
		hgtText := strings.ReplaceAll(buf2.String(), "null", "0")
		l1 := strings.Index(sgtText, "{")
		l2 := strings.Index(hgtText, "{")
		if l1 < 0 || l2 < 0 {
			break
		}
		var sgtRes struct {
			Result struct {
				Data []map[string]any `json:"data"`
			} `json:"result"`
		}
		var hgtRes struct {
			Result struct {
				Data []map[string]any `json:"data"`
			} `json:"result"`
		}
		if err := json.Unmarshal([]byte(sgtText), &sgtRes); err != nil {
			break
		}
		if err := json.Unmarshal([]byte(hgtText), &hgtRes); err != nil {
			break
		}
		sgtData := sgtRes.Result.Data
		hgtData := hgtRes.Result.Data
		if len(sgtData) == 0 {
			break
		}
		isEnd := false
		for i := range hgtData {
			if !hasStart && i >= 30 {
				isEnd = true
				break
			}
			dt, _ := time.Parse("2006-01-02 15:04:05", toString(hgtData[i]["TRADE_DATE"]))
			if hasStart && start.After(dt) {
				isEnd = true
				break
			}
			nh := parseF(toString(hgtData[i]["NET_DEAL_AMT"])) * 1000000
			bh := parseF(toString(hgtData[i]["BUY_AMT"])) * 1000000
			sh := parseF(toString(hgtData[i]["SELL_AMT"])) * 1000000
			ns := parseF(toString(sgtData[i]["NET_DEAL_AMT"])) * 1000000
			bs := parseF(toString(sgtData[i]["BUY_AMT"])) * 1000000
			ss := parseF(toString(sgtData[i]["SELL_AMT"])) * 1000000
			out = append(out, Daily{
				TradeDate: dt.Format("2006-01-02"),
				NetHgt:    nh, BuyHgt: bh, SellHgt: sh,
				NetSgt: ns, BuySgt: bs, SellSgt: ss,
				NetTgt: nh + ns, BuyTgt: bh + bs, SellTgt: sh + ss,
			})
		}
		if isEnd {
			break
		}
		currPage++
	}
	return out, nil
}

func (c *Client) Minute(ctx context.Context) ([]Minute, error) {
	r, _ := c.minuteEast(ctx)
	if len(r) == 0 {
		r, _ = c.minuteThs(ctx)
	}
	return r, nil
}

func (c *Client) Current(ctx context.Context) (Minute, error) {
	mins, _ := c.Minute(ctx)
	if len(mins) == 0 {
		return Minute{}, nil
	}
	return mins[len(mins)-1], nil
}

func (c *Client) minuteThs(ctx context.Context) ([]Minute, error) {
	if c.cfg.Wait > 0 {
		time.Sleep(c.cfg.Wait)
	}
	resp, err := c.client.R().SetContext(ctx).Get("https://data.hexin.cn/market/hsgtApi/method/dayChart/")
	if err != nil {
		return []Minute{}, nil
	}
	var res struct {
		Time []string  `json:"time"`
		Hgt  []float64 `json:"hgt"`
		Sgt  []float64 `json:"sgt"`
	}
	if err := json.Unmarshal(resp.Body(), &res); err != nil {
		return []Minute{}, nil
	}
	now := time.Now()
	yrs := tradecalendar.CalendarYears()
	var days []tradecalendar.Day
	for _, y := range yrs {
		if y == now.Year() {
			d, _ := tradecalendar.TradeCalendar(y)
			days = d
			break
		}
	}
	var latest string
	if len(days) > 0 {
		for i := len(days) - 1; i >= 0; i-- {
			dt, _ := time.Parse("2006-01-02", days[i].TradeDate)
			if !dt.After(now) && days[i].TradeStatus == 1 {
				latest = days[i].TradeDate
				break
			}
		}
	}
	out := make([]Minute, 0, len(res.Time))
	for i := range res.Time {
		tr := latest + " " + res.Time[i]
		out = append(out, Minute{TradeTime: tr, NetHgt: float64(int64(res.Hgt[i] * 100000000)), NetSgt: float64(int64(res.Sgt[i] * 100000000)), NetTgt: float64(int64((res.Hgt[i] + res.Sgt[i]) * 100000000))})
	}
	return out, nil
}

func (c *Client) minuteEast(ctx context.Context) ([]Minute, error) {
	url := "https://push2.eastmoney.com/api/qt/kamt.rtmin/get?fields1=f1,f3&fields2=f51,f52,f54,f56&ut=b2884a393a59ad64002292a3e90d46a5"
	if c.cfg.Wait > 0 {
		time.Sleep(c.cfg.Wait)
	}
	resp, err := c.client.R().SetContext(ctx).Get(url)
	if err != nil {
		return []Minute{}, nil
	}
	buf := new(strings.Builder)
	if _, err := io.Copy(buf, strings.NewReader(resp.String())); err != nil {
		return []Minute{}, nil
	}
	text := buf.String()
	l := strings.Index(text, "{")
	if l < 0 {
		return []Minute{}, nil
	}
	var res struct {
		Data struct {
			S2nDate string   `json:"s2nDate"`
			S2n     []string `json:"s2n"`
		} `json:"data"`
	}
	if err := json.Unmarshal([]byte(text[l:len(text)-2]), &res); err != nil {
		return []Minute{}, nil
	}
	y := time.Now().Format("2006")
	out := make([]Minute, 0, len(res.Data.S2n))
	for _, row := range res.Data.S2n {
		cols := strings.Split(row, ",")
		if len(cols) < 4 {
			continue
		}
		if cols[1] == "-" {
			continue
		}
		out = append(out, Minute{TradeTime: y + "-" + res.Data.S2nDate + " " + cols[0], NetHgt: float64(int64(parseF(cols[1]) * 10000)), NetSgt: float64(int64(parseF(cols[2]) * 10000)), NetTgt: float64(int64(parseF(cols[3]) * 10000))})
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

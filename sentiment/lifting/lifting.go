package lifting

import (
	"context"
	"regexp"
	"strconv"
	"strings"
	"time"

	browser "github.com/EDDYCJY/fake-useragent"
	"github.com/doarvid/go-adata/common/header"
	"github.com/doarvid/go-adata/common/utils"
	"github.com/go-resty/resty/v2"
)

type Row struct {
	StockCode string  `json:"stock_code"`
	ShortName string  `json:"short_name"`
	LiftDate  string  `json:"lift_date"`
	Volume    int64   `json:"volume"`
	Amount    int64   `json:"amount"`
	Ratio     float64 `json:"ratio"`
	Price     float64 `json:"price"`
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
		UserAgent: "go-adata/lifting",
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
			c.SetProxy(cfg.Proxy)
		}
	}
	if cfg.Debug {
		c.SetDebug(true)
	}
	return &Client{client: c, cfg: cfg}
}

func (sl *Client) LastMonth(ctx context.Context) ([]Row, error) {
	out := []Row{}
	for i := 1; i < 10; i++ {
		rows, err := sl.lastMonthByPage(ctx, i)
		if err != nil {
			return []Row{}, err
		}
		out = append(out, rows...)
		if len(rows) != 50 {
			break
		}
	}
	return out, nil
}

func (sl *Client) lastMonthByPage(ctx context.Context, pageNum int) ([]Row, error) {
	url := "http://data.10jqka.com.cn/market/xsjj/field/enddate/order/desc/ajax/1/free/1/"
	if pageNum > 1 {
		url = url + "page/" + strconv.Itoa(pageNum) + "/free/1/"
	}
	if sl.cfg.Wait > 0 {
		time.Sleep(sl.cfg.Wait)
	}
	headers := header.DefaultHeaders()
	headers["Host"] = "data.10jqka.com.cn"
	headers["Cookie"] = utils.ThsCookie()
	resp, err := sl.client.R().SetContext(ctx).SetHeaders(headers).Get(url)
	if err != nil {
		return []Row{}, nil
	}
	gbkBytes := resp.Body()
	utf8Bytes, err := utils.GBKToUTF8(gbkBytes)
	if err != nil {
		return []Row{}, err
	}
	text := string(utf8Bytes)
	if !(strings.Contains(text, "解禁日期") || strings.Contains(text, "解禁股")) {
		return []Row{}, nil
	}
	rows := parseTableRows(text)
	out := make([]Row, 0, len(rows))
	for _, r := range rows {
		if len(r) < 8 {
			continue
		}
		vol := toUnitInt(r[3])
		amt := toUnitInt(r[5])
		out = append(out, Row{
			StockCode: r[0],
			ShortName: r[1],
			LiftDate:  r[2],
			Volume:    vol,
			Ratio:     parseF(r[6]),
			Price:     parseF(r[4]),
			Amount:    amt,
		})
	}
	return out, nil
}

func parseTableRows(html string) [][]string {
	trs := strings.Split(html, "<tr")
	out := [][]string{}
	for _, seg := range trs {
		cols := []string{}
		re := regexp.MustCompile(`>([\s\S]+?)</td>`)
		for _, m := range re.FindAllStringSubmatch(seg, -1) {
			cols = append(cols, strings.TrimSpace(m[1]))
		}
		if len(cols) >= 8 {
			out = append(out, []string{toHrefText(cols[1]), toHrefText(cols[2]), cols[3], cols[4], cols[5], cols[6], cols[7], cols[8]})
		}
	}
	return out
}

func toHrefText(v string) string {
	re := regexp.MustCompile(`>([\s\S]+?)</a>`)
	ms := re.FindStringSubmatch(v)
	if len(ms) >= 2 {
		return ms[1]
	}
	return ""
}

func toUnitInt(s string) int64 {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0
	}
	if strings.HasSuffix(s, "万") {
		v, _ := strconv.ParseFloat(strings.TrimSuffix(s, "万"), 64)
		return int64(v * 10000)
	}
	if strings.HasSuffix(s, "亿") {
		v, _ := strconv.ParseFloat(strings.TrimSuffix(s, "亿"), 64)
		return int64(v * 100000000)
	}
	v, _ := strconv.ParseFloat(s, 64)
	return int64(v)
}

func parseF(s string) float64 {
	s = strings.TrimSpace(strings.ReplaceAll(s, "%", ""))
	if s == "" || s == "--" {
		return 0
	}
	v, _ := strconv.ParseFloat(s, 64)
	return v
}

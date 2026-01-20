package hot

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

type PopRankRow struct {
	Rank      int     `json:"rank"`
	StockCode string  `json:"stock_code"`
	ShortName string  `json:"short_name"`
	Price     float64 `json:"price"`
	Change    float64 `json:"change"`
	ChangePct float64 `json:"change_pct"`
}

type HotRankRow struct {
	Rank       int     `json:"rank"`
	StockCode  string  `json:"stock_code"`
	ShortName  string  `json:"short_name"`
	ChangePct  float64 `json:"change_pct"`
	HotValue   float64 `json:"hot_value"`
	PopTag     string  `json:"pop_tag"`
	ConceptTag string  `json:"concept_tag"`
}

type HotConceptRow struct {
	Rank        int     `json:"rank"`
	ConceptCode string  `json:"concept_code"`
	ConceptName string  `json:"concept_name"`
	ChangePct   float64 `json:"change_pct"`
	HotValue    float64 `json:"hot_value"`
	HotTag      string  `json:"hot_tag"`
}

type PlateType string

const (
	PlateTypeConcept  PlateType = "concept"
	PlateTypeIndustry PlateType = "industry"
)

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
		UserAgent: "go-adata/hot",
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

func (h *Client) Popular(ctx context.Context) ([]PopRankRow, error) {
	if h.cfg.Wait > 0 {
		time.Sleep(h.cfg.Wait)
	}
	params := map[string]any{
		"appId":      "appId01",
		"globalId":   "786e4c21-70dc-435a-93bb-38",
		"marketType": "",
		"pageNo":     1,
		"pageSize":   100,
	}
	resp, err := h.client.R().SetContext(ctx).SetBody(params).Post("https://emappdata.eastmoney.com/stockrank/getAllCurrentList")
	if err != nil {
		return nil, err
	}
	var res struct {
		Data []map[string]any `json:"data"`
	}
	if err := json.Unmarshal(resp.Body(), &res); err != nil {
		return nil, err
	}
	sc := make([]string, 0, len(res.Data))
	for _, it := range res.Data {
		sc = append(sc, toString(it["sc"]))
	}
	marks := make([]string, 0, len(sc))
	for _, item := range sc {
		if strings.HasPrefix(item, "SZ") {
			marks = append(marks, "0."+item[2:])
		} else {
			marks = append(marks, "1."+item[2:])
		}
	}
	q := strings.Join(marks, ",")
	url := "https://push2.eastmoney.com/api/qt/ulist.np/get?ut=f057cbcbce2a86e2866ab8877db1d059&fltt=2&invt=2&fields=f14,f3,f12,f2&secids=" + q
	if h.cfg.Wait > 0 {
		time.Sleep(h.cfg.Wait)
	}
	resp2, err := h.client.R().SetContext(ctx).Get(url)
	if err != nil {
		return nil, err
	}
	var res2 struct {
		Data struct {
			Diff []map[string]any `json:"diff"`
		} `json:"data"`
	}
	if err := json.Unmarshal(resp2.Body(), &res2); err != nil {
		return nil, err
	}
	out := make([]PopRankRow, 0, len(res2.Data.Diff))
	rank := 1
	for _, it := range res2.Data.Diff {
		price := parseF(toString(it["f2"]))
		pct := parseF(toString(it["f3"]))
		out = append(out, PopRankRow{Rank: rank, StockCode: toString(it["f12"]), ShortName: toString(it["f14"]), Price: price, ChangePct: pct, Change: price * pct / 100})
		rank++
	}
	return out, nil
}

func (h *Client) Stocks(ctx context.Context) ([]HotRankRow, error) {
	url := "https://dq.10jqka.com.cn/fuyao/hot_list_data/out/hot_list/v1/stock?stock_type=a&type=hour&list_type=normal"
	if h.cfg.Wait > 0 {
		time.Sleep(h.cfg.Wait)
	}
	resp, err := h.client.R().SetContext(ctx).Get(url)
	if err != nil {
		return nil, err
	}
	var res struct {
		Data struct {
			StockList []map[string]any `json:"stock_list"`
		} `json:"data"`
	}
	if err := json.Unmarshal(resp.Body(), &res); err != nil {
		return nil, err
	}
	out := make([]HotRankRow, 0, len(res.Data.StockList))
	for _, d := range res.Data.StockList {
		conceptTags := []string{}
		if v, ok := d["tag"].(map[string]any); ok {
			if ct, ok2 := v["concept_tag"].([]any); ok2 {
				for _, x := range ct {
					conceptTags = append(conceptTags, toString(x))
				}
			}
		}
		popTag := ""
		if v, ok := d["tag"].(map[string]any); ok {
			if pt, ok2 := v["popularity_tag"].(string); ok2 {
				popTag = strings.ReplaceAll(pt, "\n", "")
			}
		}
		out = append(out, HotRankRow{
			Rank:       int(parseF(toString(d["order"]))),
			StockCode:  toString(d["code"]),
			ShortName:  toString(d["name"]),
			ChangePct:  parseF(toString(d["rise_and_fall"])),
			HotValue:   parseF(toString(d["rate"])),
			PopTag:     popTag,
			ConceptTag: strings.Join(conceptTags, ";"),
		})
	}
	return out, nil
}

func (h *Client) Concepts(ctx context.Context, plateType PlateType) ([]HotConceptRow, error) {
	if plateType != PlateTypeConcept && plateType != PlateTypeIndustry {
		return nil, fmt.Errorf("invalid plate type: %s", plateType)
	}
	t := string(plateType)
	url := fmt.Sprintf("https://dq.10jqka.com.cn/fuyao/hot_list_data/out/hot_list/v1/plate?type=%s", t)
	if h.cfg.Wait > 0 {
		time.Sleep(h.cfg.Wait)
	}
	resp, err := h.client.R().SetContext(ctx).Get(url)
	if err != nil {
		return nil, err
	}
	var res struct {
		Data struct {
			PlateList []map[string]any `json:"plate_list"`
		} `json:"data"`
	}
	if err := json.Unmarshal(resp.Body(), &res); err != nil {
		return nil, err
	}
	out := make([]HotConceptRow, 0, len(res.Data.PlateList))
	for _, d := range res.Data.PlateList {
		out = append(out, HotConceptRow{
			Rank:        int(parseF(toString(d["order"]))),
			ConceptCode: toString(d["code"]),
			ConceptName: toString(d["name"]),
			ChangePct:   parseF(toString(d["rise_and_fall"])),
			HotValue:    parseF(toString(d["rate"])),
			HotTag:      toString(d["hot_tag"]),
		})
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

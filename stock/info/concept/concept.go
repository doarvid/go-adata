package concept

import (
	"context"
	_ "embed"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/doarvid/go-adata/common/codeutils"
)

//go:embed all_concept_code_east.csv
var allConceptCodeCSVEast string

type ConceptCode struct {
	ConceptCode string `json:"concept_code"`
	IndexCode   string `json:"index_code"`
	Name        string `json:"name"`
	Source      string `json:"source"`
}

type ConceptInfo struct {
	StockCode   string `json:"stock_code"`
	ConceptCode string `json:"concept_code"`
	Name        string `json:"name"`
	Source      string `json:"source"`
	Reason      string `json:"reason"`
}

type Constituent struct {
	StockCode string `json:"stock_code"`
	ShortName string `json:"short_name"`
}

func LoadAllConceptCodesFromCSV() ([]ConceptCode, error) {
	r := csv.NewReader(strings.NewReader(allConceptCodeCSVEast))
	if _, err := r.Read(); err != nil {
		return nil, err
	}
	var out []ConceptCode
	for {
		rec, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		out = append(out, ConceptCode{ConceptCode: rec[0], IndexCode: rec[1], Name: rec[2], Source: rec[3]})
	}
	return out, nil
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

type Concept struct {
	client *resty.Client
	cfg    Config
}

func NewConcept(opts ...Option) *Concept {
	cfg := Config{
		Timeout:   15 * time.Second,
		UserAgent: "go-adata/concept",
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
	return &Concept{client: c, cfg: cfg}
}

func (c *Concept) AllConceptCodesEast(ctx context.Context) ([]ConceptCode, error) {
	client := c.client
	page := 1
	size := 100
	var out []ConceptCode
	for page < 50 {
		params := map[string]string{
			"pn":     strconv.Itoa(page),
			"pz":     strconv.Itoa(size),
			"po":     "1",
			"np":     "1",
			"fields": "f12,f13,f14,f62",
			"fid":    "f62",
			"fs":     "m:90+t:3",
		}
		if c.cfg.Wait > 0 {
			time.Sleep(c.cfg.Wait)
		}
		resp, err := client.R().SetContext(ctx).SetQueryParams(params).Get("https://push2.eastmoney.com/api/qt/clist/get")
		if err != nil {
			return out, err
		}
		var data struct {
			Data struct {
				Diff []map[string]any `json:"diff"`
			} `json:"data"`
		}
		if err := json.Unmarshal(resp.Body(), &data); err != nil {
			return out, err
		}
		if len(data.Data.Diff) == 0 {
			break
		}
		for _, d := range data.Data.Diff {
			name := strings.TrimSpace(toString(d["f14"]))
			code := toString(d["f12"]) // BKxxxx
			out = append(out, ConceptCode{ConceptCode: code, IndexCode: code, Name: name, Source: "东方财富"})
		}
		if len(data.Data.Diff) < size {
			break
		}
		page++
	}
	// 合并缓存去重
	cached, _ := LoadAllConceptCodesFromCSV()
	seen := map[string]bool{}
	for _, c := range out {
		seen[c.ConceptCode] = true
	}
	for _, c := range cached {
		if !seen[c.ConceptCode] {
			out = append(out, c)
		}
	}
	return out, nil
}

func (c *Concept) GetConceptEast(ctx context.Context, stockCode string) ([]ConceptInfo, error) {
	client := c.client
	sc := codeutils.CompileExchangeByStockCode(stockCode)
	params := map[string]string{
		"reportName":   "RPT_F10_CORETHEME_BOARDTYPE",
		"columns":      "SECUCODE,SECURITY_CODE,SECURITY_NAME_ABBR,NEW_BOARD_CODE,BOARD_NAME,SELECTED_BOARD_REASON,IS_PRECISE,BOARD_RANK,BOARD_YIELD,DERIVE_BOARD_CODE",
		"quoteColumns": "f3~05~NEW_BOARD_CODE~BOARD_YIELD",
		"filter":       "(SECUCODE=\"" + sc + "\")(IS_PRECISE=\"1\")",
		"pageNumber":   "1",
		"pageSize":     "50",
		"sortTypes":    "1",
		"sortColumns":  "BOARD_RANK",
		"source":       "HSF10",
		"client":       "PC",
	}
	if c.cfg.Wait > 0 {
		time.Sleep(c.cfg.Wait)
	}
	resp, err := client.R().SetContext(ctx).SetQueryParams(params).Get("https://datacenter.eastmoney.com/securities/api/data/v1/get")
	if err != nil {
		return nil, err
	}
	var data struct {
		Result struct {
			Data []map[string]any `json:"data"`
		} `json:"result"`
	}
	if err := json.Unmarshal(resp.Body(), &data); err != nil {
		return nil, err
	}
	var out []ConceptInfo
	for _, d := range data.Result.Data {
		out = append(out, ConceptInfo{
			StockCode:   stockCode,
			ConceptCode: toString(d["NEW_BOARD_CODE"]),
			Name:        toString(d["BOARD_NAME"]),
			Source:      "东方财富",
			Reason:      toString(d["SELECTED_BOARD_REASON"]),
		})
	}
	return out, nil
}

func (c *Concept) ConstituentEast(ctx context.Context, conceptCode string) ([]Constituent, error) {
	client := c.client
	var out []Constituent
	page := 1
	for page < 100 {
		params := map[string]string{
			"fid":    "f62",
			"po":     "1",
			"pz":     "200",
			"pn":     strconv.Itoa(page),
			"np":     "1",
			"fltt":   "2",
			"invt":   "2",
			"fs":     "b:" + conceptCode,
			"fields": "f12,f14",
		}
		if c.cfg.Wait > 0 {
			time.Sleep(c.cfg.Wait)
		}
		resp, err := client.R().SetContext(ctx).SetQueryParams(params).Get("https://push2.eastmoney.com/api/qt/clist/get")
		if err != nil {
			return out, err
		}
		var data struct {
			Data struct {
				Diff []map[string]any `json:"diff"`
			} `json:"data"`
		}
		if err := json.Unmarshal(resp.Body(), &data); err != nil {
			return out, err
		}
		if len(data.Data.Diff) == 0 {
			break
		}
		for _, d := range data.Data.Diff {
			out = append(out, Constituent{StockCode: toString(d["f12"]), ShortName: toString(d["f14"])})
		}
		page++
	}
	return out, nil
}

func toString(v any) string { return strings.TrimSpace(fmt.Sprintf("%v", v)) }

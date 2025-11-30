package stockcode

import (
	_ "embed"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/doarvid/go-adata/common/codeutils"
	header "github.com/doarvid/go-adata/common/header"
	httpc "github.com/doarvid/go-adata/common/http"
)

//go:embed code.csv
var stockCodeCSV string

type StockCode struct {
	StockCode string     `json:"stock_code"`
	ShortName string     `json:"short_name"`
	Exchange  string     `json:"exchange"`
	ListDate  *time.Time `json:"list_date"`
}

func AllCode(wait time.Duration) ([]StockCode, error) {
	baidu, _ := marketRankBaidu(wait)
	var base []StockCode
	if len(baidu) >= 5000 {
		base = baidu
	} else {
		east, _ := marketRankEast(wait)
		if len(east) >= 5000 {
			base = east
		} else {
			sina, _ := marketRankSina(wait)
			base = sina
		}
	}
	ipo, _ := newSubEast(wait)
	codes := mergeUnique(ipo, base)
	if len(codes) == 0 {
		return nil, errors.New("no codes fetched")
	}
	codes = mergeListDateFromCSV(codes)
	sort.Slice(codes, func(i, j int) bool { return codes[i].StockCode < codes[j].StockCode })
	return codes, nil
}

func marketRankEast(wait time.Duration) ([]StockCode, error) {
	client := httpc.NewClient()
	url := "https://82.push2.eastmoney.com/api/qt/clist/get"
	page := 1
	pageSize := 50
	var res []StockCode
	for page < 200 {
		params := map[string]string{
			"pn":     strconv.Itoa(page),
			"pz":     strconv.Itoa(pageSize),
			"po":     "1",
			"np":     "1",
			"ut":     "bd1d9ddb04089700cf9c27f6f7426281",
			"fltt":   "2",
			"invt":   "2",
			"fid":    "f3",
			"fs":     "m:0 t:6,m:0 t:80,m:1 t:2,m:1 t:23,m:0 t:81 s:2048",
			"fields": "f12,f14",
			"_":      "1623833739532",
		}
		if wait > 0 {
			time.Sleep(wait)
		}
		resp, err := client.R().SetQueryParams(params).Get(url)
		if err != nil {
			return res, err
		}
		var data struct {
			Data struct {
				Diff []map[string]any `json:"diff"`
			} `json:"data"`
		}
		if err := json.Unmarshal(resp.Body(), &data); err != nil {
			return res, err
		}
		if len(data.Data.Diff) == 0 {
			break
		}
		for _, d := range data.Data.Diff {
			code := toString(d["f12"])
			name := toString(d["f14"])
			if code == "" {
				continue
			}
			res = append(res, StockCode{StockCode: code, ShortName: strings.ReplaceAll(name, " ", ""), Exchange: codeutils.GetExchangeByStockCode(code)})
		}
		if len(data.Data.Diff) < pageSize {
			break
		}
		page++
	}
	return res, nil
}

func marketRankBaidu(wait time.Duration) ([]StockCode, error) {
	client := httpc.NewClient()
	baseURL := "https://finance.pae.baidu.com/selfselect/getmarketrank"
	maxPageSize := 200
	out := make([]StockCode, 0, 5000)
	for page := 0; page < 49; page++ {
		if wait > 0 {
			time.Sleep(wait)
		}
		params := map[string]string{
			"sort_type":     "1",
			"sort_key":      "14",
			"from_mid":      "1",
			"group":         "pclist",
			"type":          "ab",
			"finClientType": "pc",
			"pn":            strconv.Itoa(page * maxPageSize),
			"rn":            strconv.Itoa(maxPageSize),
		}
		resp, err := client.R().SetHeaders(header.BaiduJSONHeaders()).SetQueryParams(params).Get(baseURL)
		if err != nil {
			continue
		}
		var res any
		if err := json.Unmarshal(resp.Body(), &res); err != nil {
			continue
		}
		m, ok := res.(map[string]any)
		if !ok {
			continue
		}
		if toString(m["ResultCode"]) != "0" {
			continue
		}
		result := m["Result"]
		rm, ok := result.(map[string]any)
		if !ok {
			continue
		}
		arrAny := rm["Result"].([]any)
		if len(arrAny) == 0 {
			break
		}
		disp := arrAny[0].(map[string]any)["DisplayData"].(map[string]any)
		tpl := disp["resultData"].(map[string]any)["tplData"].(map[string]any)
		ranksAny := tpl["result"].(map[string]any)["rank"]
		ranks, ok := ranksAny.([]any)
		if !ok || len(ranks) == 0 {
			break
		}
		for _, it := range ranks {
			im := it.(map[string]any)
			code := strings.TrimSpace(toString(im["code"]))
			name := strings.ReplaceAll(strings.TrimSpace(toString(im["name"])), " ", "")
			ex := strings.TrimSpace(toString(im["exchange"]))
			if code == "" {
				continue
			}
			if ex == "" {
				ex = codeutils.GetExchangeByStockCode(code)
			}
			out = append(out, StockCode{StockCode: code, ShortName: name, Exchange: ex})
		}
	}
	return out, nil
}

func marketRankSina(wait time.Duration) ([]StockCode, error) {
	client := httpc.NewClient()
	out := make([]StockCode, 0, 5000)
	for page := 1; page < 200; page++ {
		if wait > 0 {
			time.Sleep(wait)
		}
		url := "https://vip.stock.finance.sina.com.cn/quotes_service/api/json_v2.php/Market_Center.getHQNodeData"
		params := map[string]string{
			"page":   strconv.Itoa(page),
			"num":    "80",
			"sort":   "changepercent",
			"asc":    "0",
			"node":   "hs_a",
			"symbol": "",
			"_s_r_a": "page",
		}
		resp, err := client.R().SetQueryParams(params).Get(url)
		if err != nil {
			continue
		}
		var arr []map[string]any
		if err := json.Unmarshal(resp.Body(), &arr); err != nil {
			// some responses may be non-standard JSON; skip page
			continue
		}
		if len(arr) == 0 {
			break
		}
		for _, it := range arr {
			code := strings.TrimSpace(toString(it["code"]))
			name := strings.ReplaceAll(strings.TrimSpace(toString(it["name"])), " ", "")
			if code == "" {
				continue
			}
			out = append(out, StockCode{StockCode: code, ShortName: name, Exchange: codeutils.GetExchangeByStockCode(code)})
		}
		if len(arr) < 80 {
			break
		}
	}
	return out, nil
}

func newSubEast(wait time.Duration) ([]StockCode, error) {
	client := httpc.NewClient()
	var res []StockCode
	for i := 0; i < 200; i++ {
		url := "https://datacenter-web.eastmoney.com/api/data/v1/get"
		params := map[string]string{
			"sortColumns": "APPLY_DATE,SECURITY_CODE",
			"sortTypes":   "-1,-1",
			"pageSize":    "50",
			"pageNumber":  strconv.Itoa(i + 1),
			"reportName":  "RPTA_APP_IPOAPPLY",
			"columns":     "SECURITY_CODE,SECURITY_NAME,TRADE_MARKET,LISTING_DATE",
			"quoteType":   "0",
			"filter":      "(APPLY_DATE>'2010-01-01')",
			"source":      "WEB",
			"client":      "WEB",
		}
		if wait > 0 {
			time.Sleep(wait)
		}
		resp, err := client.R().SetQueryParams(params).Get(url)
		if err != nil {
			return res, err
		}
		var data struct {
			Result struct {
				Data []struct {
					SECURITY_CODE string `json:"SECURITY_CODE"`
					SECURITY_NAME string `json:"SECURITY_NAME"`
					TRADE_MARKET  string `json:"TRADE_MARKET"`
					LISTING_DATE  string `json:"LISTING_DATE"`
				} `json:"data"`
			} `json:"result"`
		}
		if err := json.Unmarshal(resp.Body(), &data); err != nil {
			return res, err
		}
		if len(data.Result.Data) == 0 {
			break
		}
		for _, d := range data.Result.Data {
			ex := marketToExchange(d.TRADE_MARKET)
			if d.LISTING_DATE != "" {
				dt, _ := time.Parse("2006-01-02", d.LISTING_DATE)
				res = append(res, StockCode{StockCode: d.SECURITY_CODE, ShortName: d.SECURITY_NAME, Exchange: ex, ListDate: &dt})
			}
		}
		if len(res) > 0 {
			last := res[len(res)-1]
			if last.ListDate != nil && last.ListDate.Before(time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local)) {
				break
			}
		}
	}
	return res, nil
}

func mergeUnique(a, b []StockCode) []StockCode {
	seen := map[string]bool{}
	out := make([]StockCode, 0, len(a)+len(b))
	for _, s := range append(a, b...) {
		if seen[s.StockCode] {
			continue
		}
		seen[s.StockCode] = true
		out = append(out, s)
	}
	return out
}

func mergeListDateFromCSV(codes []StockCode) []StockCode {
	r := csv.NewReader(strings.NewReader(stockCodeCSV))
	hdr, _ := r.Read()
	idxCode := indexOf(hdr, "stock_code")
	idxDate := indexOf(hdr, "list_date2")
	if idxDate < 0 {
		idxDate = indexOf(hdr, "list_date")
	}
	if idxCode < 0 || idxDate < 0 {
		return codes
	}
	dates := map[string]string{}
	for {
		rec, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			break
		}
		sc := strings.TrimSpace(rec[idxCode])
		ld := strings.TrimSpace(rec[idxDate])
		if sc != "" && ld != "" {
			dates[pad6(sc)] = ld
		}
	}
	for i := range codes {
		if codes[i].ListDate == nil {
			if d, ok := dates[pad6(strings.TrimSpace(codes[i].StockCode))]; ok {
				if dt, err := time.Parse("2006-01-02", d); err == nil {
					codes[i].ListDate = &dt
				}
			}
		}
	}
	return codes
}

func marketToExchange(s string) string {
	if strings.HasPrefix(s, "上海") {
		return "SH"
	}
	if strings.HasPrefix(s, "深圳") {
		return "SZ"
	}
	if strings.HasPrefix(s, "北京") {
		return "BJ"
	}
	return ""
}

func indexOf(arr []string, t string) int {
	for i, v := range arr {
		if v == t {
			return i
		}
	}
	return -1
}

func pad6(s string) string {
	if len(s) >= 6 {
		return s
	}
	return strings.Repeat("0", 6-len(s)) + s
}

func toString(v any) string { return fmt.Sprintf("%v", v) }

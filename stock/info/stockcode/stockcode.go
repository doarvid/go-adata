package stockcode

import (
    "encoding/csv"
    "encoding/json"
    "errors"
    "fmt"
    "io"
    "os"
    "sort"
    "strconv"
    "strings"
    "time"

    "go-adata/pkg/adata/common/codeutils"
    httpc "go-adata/pkg/adata/common/http"
    "go-adata/pkg/adata/stock/cache"
)

type StockCode struct {
	StockCode string     `json:"stock_code"`
	ShortName string     `json:"short_name"`
	Exchange  string     `json:"exchange"`
	ListDate  *time.Time `json:"list_date"`
}

func AllCode(wait time.Duration) ([]StockCode, error) {
	east, _ := marketRankEast(wait)
	ipo, _ := newSubEast(wait)
	codes := mergeUnique(ipo, east)
	if len(codes) == 0 {
		return nil, errors.New("no codes fetched")
	}
	codes = mergeListDateFromCSV(codes, cache.GetAllCodeCSVPath())
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
        if wait > 0 { time.Sleep(wait) }
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
        if wait > 0 { time.Sleep(wait) }
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

func mergeListDateFromCSV(codes []StockCode, csvPath string) []StockCode {
	f, err := os.Open(csvPath)
	if err != nil {
		return codes
	}
	defer f.Close()
	r := csv.NewReader(f)
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

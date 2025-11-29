package capitalflow

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	httpc "go-adata/pkg/adata/common/http"
)

type FlowMin struct {
	StockCode     string  `json:"stock_code"`
	TradeTime     string  `json:"trade_time"`
	MainNetInflow float64 `json:"main_net_inflow"`
	SmNetInflow   float64 `json:"sm_net_inflow"`
	MidNetInflow  float64 `json:"mid_net_inflow"`
	LgNetInflow   float64 `json:"lg_net_inflow"`
	MaxNetInflow  float64 `json:"max_net_inflow"`
}

type FlowDaily struct {
	StockCode     string  `json:"stock_code"`
	TradeDate     string  `json:"trade_date"`
	MainNetInflow float64 `json:"main_net_inflow"`
	SmNetInflow   float64 `json:"sm_net_inflow"`
	MidNetInflow  float64 `json:"mid_net_inflow"`
	LgNetInflow   float64 `json:"lg_net_inflow"`
	MaxNetInflow  float64 `json:"max_net_inflow"`
}

func GetStockCapitalFlowMinEast(stockCode string, wait time.Duration) ([]FlowMin, error) {
	if stockCode == "" {
		return []FlowMin{}, nil
	}
	client := httpc.NewClient()
	cid := "0"
	if strings.HasPrefix(stockCode, "6") {
		cid = "1"
	}
	params := map[string]string{
		"lmt":     "0",
		"klt":     "1",
		"fields1": "f1,f2,f3,f7",
		"fields2": "f51,f52,f53,f54,f55,f56,f57,f58,f59,f60,f61,f62,f63,f64,f65",
		"secid":   cid + "." + stockCode,
	}
	if wait > 0 {
		time.Sleep(wait)
	}
	resp, err := client.R().SetQueryParams(params).Get("https://push2.eastmoney.com/api/qt/stock/fflow/kline/get")
	if err != nil {
		return nil, err
	}
	var data struct {
		Data struct {
			Klines []string `json:"klines"`
		} `json:"data"`
	}
	if err := json.Unmarshal(resp.Body(), &data); err != nil {
		return nil, err
	}
	if len(data.Data.Klines) == 0 {
		return []FlowMin{}, nil
	}
	out := make([]FlowMin, 0, len(data.Data.Klines))
	for _, row := range data.Data.Klines {
		cols := strings.Split(row, ",")
		if len(cols) < 6 {
			continue
		}
		fm := FlowMin{StockCode: stockCode}
		fm.TradeTime = cols[0]
		fm.MainNetInflow = parseF(cols[1])
		fm.SmNetInflow = parseF(cols[2])
		fm.MidNetInflow = parseF(cols[3])
		fm.LgNetInflow = parseF(cols[4])
		fm.MaxNetInflow = parseF(cols[5])
		out = append(out, fm)
	}
	return out, nil
}

func GetStockCapitalFlowEast(stockCode string, startDate string, endDate string, wait time.Duration) ([]FlowDaily, error) {
	if stockCode == "" {
		return []FlowDaily{}, nil
	}
	client := httpc.NewClient()
	cid := "0"
	if strings.HasPrefix(stockCode, "6") {
		cid = "1"
	}
	params := map[string]string{
		"lmt":     "0",
		"klt":     "101",
		"fields1": "f1,f2,f3,f7",
		"fields2": "f51,f52,f53,f54,f55,f56,f57,f58,f59,f60,f61",
		"secid":   cid + "." + stockCode,
	}
	if wait > 0 {
		time.Sleep(wait)
	}
	resp, err := client.R().SetQueryParams(params).Get("https://push2his.eastmoney.com/api/qt/stock/fflow/daykline/get")
	if err != nil {
		return nil, err
	}
	var data struct {
		Data struct {
			Klines []string `json:"klines"`
		} `json:"data"`
	}
	if err := json.Unmarshal(resp.Body(), &data); err != nil {
		return nil, err
	}
	if len(data.Data.Klines) == 0 {
		return []FlowDaily{}, nil
	}
	out := make([]FlowDaily, 0, len(data.Data.Klines))
	for _, row := range data.Data.Klines {
		cols := strings.Split(row, ",")
		if len(cols) < 6 {
			continue
		}
		fd := FlowDaily{StockCode: stockCode}
		fd.TradeDate = cols[0]
		fd.MainNetInflow = parseF(cols[1])
		fd.SmNetInflow = parseF(cols[2])
		fd.MidNetInflow = parseF(cols[3])
		fd.LgNetInflow = parseF(cols[4])
		fd.MaxNetInflow = parseF(cols[5])
		out = append(out, fd)
	}
	// filter by date range if provided
	if startDate != "" || endDate != "" {
		s := strings.ReplaceAll(startDate, "-", "")
		e := strings.ReplaceAll(endDate, "-", "")
		filtered := make([]FlowDaily, 0, len(out))
		for _, it := range out {
			d := strings.ReplaceAll(it.TradeDate, "-", "")
			if (s == "" || d >= s) && (e == "" || d <= e) {
				filtered = append(filtered, it)
			}
		}
		out = filtered
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

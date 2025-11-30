package stockmarket

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	httpc "github.com/doarvid/go-adata/common/http"
)

// 分时 K 线结构体（百度源）
type MinuteBar struct {
	Time      int64   `json:"time"`       // 时间戳（秒）  例如：1710000000
	Price     float64 `json:"price"`      // 价格(元)     例如：9.98
	ChangePct float64 `json:"change_pct"` // 涨跌幅(%)    例如：-0.16
	Change    float64 `json:"change"`     // 涨跌额(元)   例如：-0.02
	AvgPrice  float64 `json:"avg_price"`  // 平均价(元)   例如：9.98
	Volume    int64   `json:"volume"`     // 成交量(股)   例如：64745722
	Amount    float64 `json:"amount"`     // 成交额(元)   例如：934285179.00
	Open      float64 `json:"open"`       // 开盘价(元)   例如：10.00
	Close     float64 `json:"close"`      // 收盘价(元)   例如：9.98
	High      float64 `json:"high"`       // 最高价(元)   例如：10.05
	Low       float64 `json:"low"`        // 最低价(元)   例如：9.95
	TradeTime string  `json:"trade_time"` // 交易时间     例如：2024-01-01 14:55:00
	TradeDate string  `json:"trade_date"` // 交易日期     例如：2024-01-01
	StockCode string  `json:"stock_code"` // 股票代码     例如：600001
}

// 逐笔成交结构体（百度源）
type TickBar struct {
	TradeTime string  `json:"trade_time"` // 成交时间	例如：2023-09-13 09:31:45
	Volume    int64   `json:"volume"`     // 成交量(股)	例如：34452500
	Price     float64 `json:"price"`      // 当前价格(元)	例如：12.36
	Type      string  `json:"type"`       // 类型	例如：--
	BsType    string  `json:"bs_type"`    // 买卖类型	B：买入，S：卖出
	StockCode string  `json:"stock_code"` // 代码	例如：600001
}

func GetMarketDailyBaidu(stockCode string, startDate string, kType int, wait time.Duration) ([]DailyBar, error) {
	client := httpc.NewClient()
	url := fmt.Sprintf("https://finance.pae.baidu.com/selfselect/getstockquotation?all=1&isIndex=false&isBk=false&isBlock=false&isFutures=false&isStock=true&newFormat=1&group=quotation_kline_ab&finClientType=pc&code=%s&start_time=%s 00:00:00&ktype=%d", stockCode, startDate, kType)
	var res struct {
		ResultCode string `json:"ResultCode"`
		Result     struct {
			NewMarketData struct {
				Keys       []string `json:"keys"`
				MarketData string   `json:"marketData"`
			} `json:"newMarketData"`
		} `json:"Result"`
	}
	// retry
	for i := 0; i < 3; i++ {
		if wait > 0 {
			time.Sleep(wait)
		}
		resp, err := client.R().Get(url)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(resp.Body(), &res); err != nil {
			return nil, err
		}
		if res.ResultCode == "0" {
			break
		}
		time.Sleep(2 * time.Second)
	}
	if len(res.Result.NewMarketData.Keys) == 0 || res.Result.NewMarketData.MarketData == "" {
		return []DailyBar{}, nil
	}
	keys := res.Result.NewMarketData.Keys
	raw := strings.Split(res.Result.NewMarketData.MarketData, ";")
	out := make([]DailyBar, 0, len(raw))
	for _, row := range raw {
		cols := strings.Split(row, ",")
		if len(cols) != len(keys) {
			continue
		}
		// map keys to values
		kv := map[string]string{}
		for i := range keys {
			kv[keys[i]] = cols[i]
		}
		bar := DailyBar{StockCode: stockCode}
		bar.TradeTime = kv["time"]
		bar.TradeDate = kv["time"]
		bar.Open = parseF(kv["open"])
		bar.Close = parseF(kv["close"])
		bar.High = parseF(kv["high"])
		bar.Low = parseF(kv["low"])
		bar.Volume = parseF(kv["volume"])
		bar.Amount = parseF(kv["amount"])
		bar.Change = parseF(strings.ReplaceAll(kv["range"], "+", ""))
		bar.ChangePct = parseF(strings.ReplaceAll(kv["ratio"], "+", ""))
		bar.TurnoverRatio = kv["turnoverratio"]
		bar.PreClose = parseF(kv["preClose"])
		// skip invalid rows
		if !(bar.Amount > 0 || bar.Volume > 0) {
			continue
		}
		out = append(out, bar)
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

func GetMarketMinuteBaidu(stockCode string, wait time.Duration) ([]MinuteBar, error) {
	client := httpc.NewClient()
	url := fmt.Sprintf("https://finance.pae.baidu.com/selfselect/getstockquotation?all=1&isIndex=false&isBk=false&isBlock=false&isFutures=false&isStock=true&newFormat=1&group=quotation_minute_ab&finClientType=pc&code=%s", stockCode)
	var res struct {
		ResultCode string `json:"ResultCode"`
		Result     struct {
			Priceinfo []map[string]any `json:"priceinfo"`
		} `json:"Result"`
	}
	for i := 0; i < 3; i++ {
		if wait > 0 {
			time.Sleep(wait)
		}
		resp, err := client.R().Get(url)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(resp.Body(), &res); err != nil {
			return nil, err
		}
		if res.ResultCode == "0" {
			break
		}
		time.Sleep(2 * time.Second)
	}
	list := res.Result.Priceinfo
	out := make([]MinuteBar, 0, len(list))
	for _, it := range list {
		t := toInt64(it["time"])
		vol := toInt64(it["volume"]) * 100
		bar := MinuteBar{
			Time:      t,
			Price:     parseF(toString(it["price"])),
			ChangePct: parseF(strings.ReplaceAll(toString(it["ratio"]), "%", "")),
			Change:    parseF(strings.ReplaceAll(toString(it["increase"]), "+", "")),
			AvgPrice:  parseF(toString(it["avgPrice"])),
			Volume:    vol,
			Amount:    parseF(toString(it["oriAmount"])),
			StockCode: stockCode,
		}
		tm := time.Unix(bar.Time, 0).In(time.FixedZone("CST-8", 8*3600))
		bar.TradeTime = tm.Format("2006-01-02 15:04:05")
		bar.TradeDate = tm.Format("2006-01-02")
		out = append(out, bar)
	}
	return out, nil
}

func GetMarketBarBaidu(stockCode string, wait time.Duration) ([]TickBar, error) {
	client := httpc.NewClient()
	url := fmt.Sprintf("https://finance.pae.baidu.com/vapi/v1/getquotation?srcid=5353&all=1&pointType=string&group=quotation_minute_ab&query=%s&code=%s&market_type=ab&newFormat=1&finClientType=pc", stockCode, stockCode)
	var res struct {
		ResultCode string `json:"ResultCode"`
		Result     struct {
			Detailinfos []map[string]any `json:"detailinfos"`
		} `json:"Result"`
	}
	for i := 0; i < 3; i++ {
		if wait > 0 {
			time.Sleep(wait)
		}
		resp, err := client.R().Get(url)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(resp.Body(), &res); err != nil {
			return nil, err
		}
		if res.ResultCode == "0" {
			break
		}
		time.Sleep(1 * time.Second)
	}
	list := res.Result.Detailinfos
	out := make([]TickBar, 0, len(list))
	for _, it := range list {
		tt := toInt64(it["time"])
		tm := time.Unix(tt, 0).In(time.FixedZone("CST-8", 8*3600))
		out = append(out, TickBar{
			TradeTime: tm.Format("2006-01-02 15:04:05"),
			Volume:    toInt64(it["volume"]),
			Price:     parseF(toString(it["price"])),
			Type:      toString(it["type"]),
			BsType:    toString(it["bsFlag"]),
			StockCode: stockCode,
		})
	}
	return out, nil
}

func toInt64(v any) int64 {
	s := strings.TrimSpace(fmt.Sprintf("%v", v))
	if s == "" {
		return 0
	}
	i, _ := strconv.ParseInt(s, 10, 64)
	return i
}

func toString(v any) string { return strings.TrimSpace(fmt.Sprintf("%v", v)) }

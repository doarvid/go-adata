package stockmarket

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	httpc "github.com/doarvid/go-adata/common/http"
)

func GetMarketDailyBaidu(stockCode string, startDate string, kType KType, wait time.Duration) ([]DailyBar, error) {
	client := httpc.NewClient()
	url := fmt.Sprintf("https://finance.pae.baidu.com/selfselect/getstockquotation?all=1&isIndex=false&isBk=false&isBlock=false&isFutures=false&isStock=true&newFormat=1&group=quotation_kline_ab&finClientType=pc&code=%s&start_time=%s%%2000:00:00&ktype=%d", stockCode, startDate, int(kType))
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

func GetMarketFiveBaidu(stockCode string, wait time.Duration) (Five, error) {
	client := httpc.NewClient()
	url := fmt.Sprintf("https://finance.pae.baidu.com/vapi/v1/getquotation?srcid=5353&all=1&pointType=string&group=quotation_minute_ab&query=%s&code=%s&market_type=ab&newFormat=1&finClientType=pc", stockCode, stockCode)
	var res struct {
		Result struct {
			Askinfos   []map[string]any `json:"askinfos"`
			Buyinfos   []map[string]any `json:"buyinfos"`
			Basicinfos struct {
				Name string `json:"name"`
			} `json:"basicinfos"`
		} `json:"Result"`
	}
	if wait > 0 {
		time.Sleep(wait)
	}
	resp, err := client.R().Get(url)
	if err != nil {
		return Five{}, err
	}
	if err := json.Unmarshal(resp.Body(), &res); err != nil {
		return Five{}, err
	}
	sell := res.Result.Askinfos
	buy := res.Result.Buyinfos
	f := Five{StockCode: stockCode, ShortName: res.Result.Basicinfos.Name}
	if len(sell) >= 5 && len(buy) >= 5 {
		f.S5 = parseF(toString(sell[0]["askprice"]))
		f.Sv5 = toInt64(sell[0]["askvolume"])
		f.S4 = parseF(toString(sell[1]["askprice"]))
		f.Sv4 = toInt64(sell[1]["askvolume"])
		f.S3 = parseF(toString(sell[2]["askprice"]))
		f.Sv3 = toInt64(sell[2]["askvolume"])
		f.S2 = parseF(toString(sell[3]["askprice"]))
		f.Sv2 = toInt64(sell[3]["askvolume"])
		f.S1 = parseF(toString(sell[4]["askprice"]))
		f.Sv1 = toInt64(sell[4]["askvolume"])
		f.B1 = parseF(toString(buy[0]["bidprice"]))
		f.Bv1 = toInt64(buy[0]["bidvolume"])
		f.B2 = parseF(toString(buy[1]["bidprice"]))
		f.Bv2 = toInt64(buy[1]["bidvolume"])
		f.B3 = parseF(toString(buy[2]["bidprice"]))
		f.Bv3 = toInt64(buy[2]["bidvolume"])
		f.B4 = parseF(toString(buy[3]["bidprice"]))
		f.Bv4 = toInt64(buy[3]["bidvolume"])
		f.B5 = parseF(toString(buy[4]["bidprice"]))
		f.Bv5 = toInt64(buy[4]["bidvolume"])
	}
	return f, nil
}

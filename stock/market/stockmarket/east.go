package stockmarket

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	httpc "github.com/doarvid/go-adata/common/http"
)

func GetMarketDailyEast(stockCode string, startDate string, endDate string, kType KType, adjustType int, wait time.Duration) ([]DailyBar, error) {
	client := httpc.NewClient()
	seCid := "0"
	if strings.HasPrefix(stockCode, "6") {
		seCid = "1"
	}
	beg := strings.ReplaceAll(startDate, "-", "")
	if beg == "" {
		beg = "19900101"
	}
	end := strings.ReplaceAll(endDate, "-", "")
	if end == "" {
		end = time.Now().Format("20060102")
	}
	klt := int(kType)
	if klt < 5 {
		klt = 100 + klt
	}
	params := map[string]string{
		"fields1": "f1,f2,f3,f4,f5,f6",
		"fields2": "f51,f52,f53,f54,f55,f56,f57,f58,f59,f60,f61,f116",
		"ut":      "7eea3edcaed734bea9cbfc24409ed989",
		"klt":     strconv.Itoa(klt),
		"fqt":     strconv.Itoa(adjustType),
		"secid":   seCid + "." + stockCode,
		"beg":     beg,
		"end":     end,
		"_":       "1623766962675",
	}
	if wait > 0 {
		time.Sleep(wait)
	}
	resp, err := client.R().SetQueryParams(params).Get("http://push2his.eastmoney.com/api/qt/stock/kline/get")
	if err != nil {
		return nil, err
	}
	var data struct {
		Data struct {
			Klines   []string `json:"klines"`
			PreClose float64  `json:"preClose"`
		} `json:"data"`
	}
	if err := json.Unmarshal(resp.Body(), &data); err != nil {
		return nil, err
	}
	if len(data.Data.Klines) == 0 {
		return []DailyBar{}, nil
	}
	out := make([]DailyBar, 0, len(data.Data.Klines))
	for _, row := range data.Data.Klines {
		cols := strings.Split(row, ",")
		if len(cols) < 11 {
			continue
		}
		bar := DailyBar{StockCode: stockCode}
		bar.TradeTime = cols[0]
		bar.TradeDate = cols[0]
		bar.Open = parseF(cols[1])
		bar.Close = parseF(cols[2])
		bar.High = parseF(cols[3])
		bar.Low = parseF(cols[4])
		bar.Volume = parseF(cols[5])
		bar.Amount = parseF(cols[6])
		bar.Change = parseF(cols[7])
		bar.ChangePct = parseF(cols[8])
		bar.TurnoverRatio = cols[9]
		bar.PreClose = data.Data.PreClose
		out = append(out, bar)
	}
	return out, nil
}

func GetMarketMinuteEast(stockCode string, wait time.Duration) ([]MinuteBar, error) {
	client := httpc.NewClient()
	seCid := "0"
	if strings.HasPrefix(stockCode, "6") {
		seCid = "1"
	}
	params := map[string]string{
		"fields1": "f1,f2,f3,f4,f5,f6,f7,f8,f9,f10,f11,f12,f13",
		"fields2": "f51,f52,f53,f54,f55,f56,f57,f58",
		"ut":      "fa5fd1943c7b386f172d6893dbfba10b",
		"ndays":   "1",
		"iscr":    "1",
		"iscca":   "0",
		"secid":   seCid + "." + stockCode,
		"_":       "1623766962675",
	}
	if wait > 0 {
		time.Sleep(wait)
	}
	resp, err := client.R().SetQueryParams(params).Get("https://push2.eastmoney.com/api/qt/stock/trends2/get")
	if err != nil {
		return nil, err
	}
	var res struct {
		Data struct {
			PreClose float64  `json:"preClose"`
			Trends   []string `json:"trends"`
		} `json:"data"`
	}
	if err := json.Unmarshal(resp.Body(), &res); err != nil {
		return nil, err
	}
	if res.Data.Trends == nil {
		return []MinuteBar{}, nil
	}
	out := make([]MinuteBar, 0, len(res.Data.Trends))
	for _, row := range res.Data.Trends {
		cols := strings.Split(row, ",")
		if len(cols) < 8 {
			continue
		}
		bar := MinuteBar{StockCode: stockCode}
		bar.TradeTime = cols[0]
		bar.Price = parseF(cols[7])
		bar.Open = 0
		bar.Close = parseF(cols[2])
		bar.High = parseF(cols[3])
		bar.Low = parseF(cols[4])
		vol := parseF(cols[5])
		bar.Volume = int64(vol * 100)
		bar.Amount = parseF(cols[6])
		bar.AvgPrice = bar.Price
		bar.Change = bar.Price - res.Data.PreClose
		if res.Data.PreClose != 0 {
			bar.ChangePct = bar.Change / res.Data.PreClose * 100
		}
		out = append(out, bar)
	}
	return out, nil
}

package indexmarket

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func (im *IndexMarket) GetDailyEast(ctx context.Context, indexCode string, startDate string, kType int) ([]IndexDailyBar, error) {
	if indexCode == "" {
		return []IndexDailyBar{}, nil
	}
	secId := "0"
	if strings.HasPrefix(indexCode, "93") {
		secId = "2"
	} else if strings.HasPrefix(indexCode, "0") {
		secId = "1"
	}
	beg := strings.ReplaceAll(startDate, "-", "")
	if beg == "" {
		beg = "19900101"
	}
	params := map[string]string{
		"fields1": "f1,f2,f3,f4,f5,f6",
		"fields2": "f51,f52,f53,f54,f55,f56,f57,f58,f59,f60,f61",
		"ut":      "fa5fd1943c7b386f172d6893dbfba10b",
		"klt":     strconv.Itoa(100 + kType),
		"fqt":     "1",
		"secid":   secId + "." + indexCode,
		"beg":     beg,
		"end":     "20500101",
		"lmt":     "1000000",
	}
	if im.cfg.Wait > 0 {
		time.Sleep(im.cfg.Wait)
	}
	resp, err := im.client.R().SetContext(ctx).SetQueryParams(params).Get("https://push2his.eastmoney.com/api/qt/stock/kline/get")
	if err != nil {
		return nil, err
	}
	var data struct {
		Data struct {
			Code     string   `json:"code"`
			Klines   []string `json:"klines"`
			PreClose float64  `json:"preClose"`
		} `json:"data"`
	}
	if err := json.Unmarshal(resp.Body(), &data); err != nil {
		return nil, err
	}
	if data.Data.Code != indexCode || len(data.Data.Klines) == 0 {
		return []IndexDailyBar{}, nil
	}
	out := make([]IndexDailyBar, 0, len(data.Data.Klines))
	for _, row := range data.Data.Klines {
		cols := strings.Split(row, ",")
		if len(cols) < 10 {
			continue
		}
		bar := IndexDailyBar{IndexCode: indexCode}
		bar.TradeTime = cols[0]
		bar.TradeDate = cols[0]
		bar.Open = parseF(cols[1])
		bar.Close = parseF(cols[2])
		bar.High = parseF(cols[3])
		bar.Low = parseF(cols[4])
		bar.Volume = parseF(cols[5])
		bar.Amount = parseF(cols[6])
		bar.Change = parseF(cols[9])
		bar.ChangePct = parseF(cols[8])
		if len(cols) > 9 {
			bar.TurnoverRatio = cols[9]
		}
		bar.PreClose = data.Data.PreClose
		out = append(out, bar)
	}
	return out, nil
}

func (im *IndexMarket) GetMinuteEast(ctx context.Context, indexCode string) ([]IndexMinuteBar, error) {
	if indexCode == "" {
		return []IndexMinuteBar{}, nil
	}
	secId := "0"
	if strings.HasPrefix(indexCode, "93") {
		secId = "2"
	} else if strings.HasPrefix(indexCode, "0") {
		secId = "1"
	}
	params := map[string]string{
		"fields1": "f1,f2,f3,f4,f5,f6,f7,f8,f9,f10,f11,f12,f13",
		"fields2": "f51,f52,f53,f54,f55,f56,f57,f58",
		"ndays":   "1",
		"iscr":    "0",
		"secid":   secId + "." + indexCode,
	}
	if im.cfg.Wait > 0 {
		time.Sleep(im.cfg.Wait)
	}
	resp, err := im.client.R().SetContext(ctx).SetQueryParams(params).Get("https://push2his.eastmoney.com/api/qt/stock/trends2/get")
	if err != nil {
		return nil, err
	}
	var res struct {
		Data struct {
			Code     string   `json:"code"`
			PrePrice float64  `json:"prePrice"`
			Trends   []string `json:"trends"`
		} `json:"data"`
	}
	if err := json.Unmarshal(resp.Body(), &res); err != nil {
		return nil, err
	}
	if res.Data.Code != indexCode || res.Data.Trends == nil {
		return []IndexMinuteBar{}, nil
	}
	out := make([]IndexMinuteBar, 0, len(res.Data.Trends))
	for _, row := range res.Data.Trends {
		cols := strings.Split(row, ",")
		if len(cols) < 8 {
			continue
		}
		bar := IndexMinuteBar{IndexCode: indexCode}
		bar.TradeTime = cols[0]
		bar.TradeDate = cols[0][:10]
		bar.Open = parseF(cols[1])
		bar.Price = parseF(cols[2])
		bar.High = parseF(cols[3])
		bar.Low = parseF(cols[4])
		vol := parseF(cols[5])
		bar.Volume = int64(vol)
		bar.Amount = parseF(cols[6])
		bar.AvgPrice = parseF(cols[7])
		bar.Change = bar.Price - res.Data.PrePrice
		if res.Data.PrePrice != 0 {
			bar.ChangePct = bar.Change / res.Data.PrePrice * 100
		}
		out = append(out, bar)
	}
	return out, nil
}

func (im *IndexMarket) GetCurrentEast(ctx context.Context, indexCode string) (IndexCurrent, error) {
	if indexCode == "" {
		return IndexCurrent{}, nil
	}
	secId := "0"
	if strings.HasPrefix(indexCode, "93") {
		secId = "2"
	} else {
		if strings.HasPrefix(indexCode, "0") {
			secId = "1"
		}
	}
	params := map[string]string{
		"invt":   "2",
		"fltt":   "1",
		"fields": "f58,f107,f57,f43,f59,f169,f170,f152,f46,f60,f44,f45,f47,f48,f19,f532,f39,f161,f49,f171,f50,f86,f600,f601,f154,f84,f85,f168,f108,f116,f167,f164,f92,f71,f117,f292,f113,f114,f115,f119,f120,f121,f122,f296",
		"secid":  secId + "." + indexCode,
		"wbp2u":  "|0|0|0|web",
	}
	if im.cfg.Wait > 0 {
		time.Sleep(im.cfg.Wait)
	}
	resp, err := im.client.R().SetContext(ctx).SetQueryParams(params).Get("https://push2.eastmoney.com/api/qt/stock/get")
	if err != nil {
		return IndexCurrent{}, err
	}
	var data struct {
		Data map[string]any `json:"data"`
	}
	if err := json.Unmarshal(resp.Body(), &data); err != nil {
		return IndexCurrent{}, err
	}
	j := data.Data
	if j == nil {
		return IndexCurrent{}, nil
	}
	code := toString(j["f57"])
	if code != indexCode {
		return IndexCurrent{}, nil
	}
	preClose := parseF(toString(j["f60"]))
	cur := IndexCurrent{IndexCode: indexCode}
	cur.Open = parseF(toString(j["f46"])) / 100
	cur.High = parseF(toString(j["f44"])) / 100
	cur.Low = parseF(toString(j["f45"])) / 100
	cur.Price = parseF(toString(j["f43"])) / 100
	cur.Volume = parseF(toString(j["f47"]))
	cur.Amount = parseF(toString(j["f48"]))
	cur.Change = cur.Price - preClose/100
	if preClose != 0 {
		cur.ChangePct = cur.Change / (preClose / 100) * 100
	}
	cur.TradeTime = time.Now().Format("2006-01-02 15:04:05")
	cur.TradeDate = time.Now().Format("2006-01-02")
	return cur, nil
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

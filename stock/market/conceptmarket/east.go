package conceptmarket

import (
    "fmt"
    "encoding/json"
    "strconv"
    "strings"
    "time"

)

type ConceptDailyBar struct {
    TradeTime string  `json:"trade_time"`
    TradeDate string  `json:"trade_date"`
    Open      float64 `json:"open"`
    Close     float64 `json:"close"`
    High      float64 `json:"high"`
    Low       float64 `json:"low"`
    Volume    float64 `json:"volume"`
    Amount    float64 `json:"amount"`
    Change    float64 `json:"change"`
    ChangePct float64 `json:"change_pct"`
    IndexCode string  `json:"index_code"`
}

type ConceptMinuteBar struct {
    TradeTime string  `json:"trade_time"`
    TradeDate string  `json:"trade_date"`
    Price     float64 `json:"price"`
    Change    float64 `json:"change"`
    ChangePct float64 `json:"change_pct"`
    Volume    int64   `json:"volume"`
    AvgPrice  float64 `json:"avg_price"`
    Amount    float64 `json:"amount"`
    Open      float64 `json:"open"`
    Close     float64 `json:"close"`
    High      float64 `json:"high"`
    Low       float64 `json:"low"`
    IndexCode string  `json:"index_code"`
}

type ConceptCurrent struct {
    TradeTime string  `json:"trade_time"`
    TradeDate string  `json:"trade_date"`
    Open      float64 `json:"open"`
    High      float64 `json:"high"`
    Low       float64 `json:"low"`
    Price     float64 `json:"price"`
    Change    float64 `json:"change"`
    ChangePct float64 `json:"change_pct"`
    Volume    float64 `json:"volume"`
    Amount    float64 `json:"amount"`
    IndexCode string  `json:"index_code"`
}

func GetConceptDailyEast(indexCode string, kType int, wait time.Duration) ([]ConceptDailyBar, error) {
    if indexCode == "" { return []ConceptDailyBar{}, nil }
    client := getHTTPClient()
    params := map[string]string{
        "secid":   "90." + indexCode,
        "fields1": "f1,f2,f3,f4,f5,f6",
        "fields2": "f51,f52,f53,f54,f55,f56,f57,f58,f59,f60,f61",
        "klt":     strconv.Itoa(100 + kType),
        "fqt":     "1",
        "end":     "20500101",
        "lmt":     "1000000",
    }
    if wait > 0 { time.Sleep(wait) }
    resp, err := client.R().SetQueryParams(params).Get("https://push2his.eastmoney.com/api/qt/stock/kline/get")
    if err != nil { return nil, err }
    var data struct{ Data struct{ Code string `json:"code"`; Klines []string `json:"klines"` } `json:"data"` }
    if err := json.Unmarshal(resp.Body(), &data); err != nil { return nil, err }
    if data.Data.Code != indexCode || len(data.Data.Klines) == 0 { return []ConceptDailyBar{}, nil }
    out := make([]ConceptDailyBar, 0, len(data.Data.Klines))
    for _, row := range data.Data.Klines {
        cols := strings.Split(row, ",")
        if len(cols) < 10 { continue }
        bar := ConceptDailyBar{IndexCode: indexCode}
        bar.TradeDate = cols[0]
        bar.TradeTime = cols[0]
        bar.Open = parseF(cols[1])
        bar.Close = parseF(cols[2])
        bar.High = parseF(cols[3])
        bar.Low = parseF(cols[4])
        bar.Volume = parseF(cols[5])
        bar.Amount = parseF(cols[6])
        bar.Change = parseF(cols[9])
        bar.ChangePct = parseF(cols[8])
        out = append(out, bar)
    }
    return out, nil
}

func GetConceptMinuteEast(indexCode string, wait time.Duration) ([]ConceptMinuteBar, error) {
    if indexCode == "" { return []ConceptMinuteBar{}, nil }
    client := getHTTPClient()
    params := map[string]string{
        "fields1": "f1,f2,f3,f4,f5,f6,f7,f8,f9,f10,f11,f12,f13",
        "fields2": "f51,f52,f53,f54,f55,f56,f57,f58",
        "ndays":   "1",
        "iscr":    "0",
        "secid":   "90." + indexCode,
    }
    if wait > 0 { time.Sleep(wait) }
    resp, err := client.R().SetQueryParams(params).Get("https://push2his.eastmoney.com/api/qt/stock/trends2/get")
    if err != nil { return nil, err }
    var res struct{ Data struct{ Code string `json:"code"`; PrePrice float64 `json:"prePrice"`; Trends []string `json:"trends"` } `json:"data"` }
    if err := json.Unmarshal(resp.Body(), &res); err != nil { return nil, err }
    if res.Data.Code != indexCode || res.Data.Trends == nil { return []ConceptMinuteBar{}, nil }
    out := make([]ConceptMinuteBar, 0, len(res.Data.Trends))
    for _, row := range res.Data.Trends {
        cols := strings.Split(row, ",")
        if len(cols) < 8 { continue }
        bar := ConceptMinuteBar{IndexCode: indexCode}
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
        if res.Data.PrePrice != 0 { bar.ChangePct = bar.Change / res.Data.PrePrice * 100 }
        out = append(out, bar)
    }
    return out, nil
}

func GetConceptCurrentEast(indexCode string, wait time.Duration) (ConceptCurrent, error) {
    if indexCode == "" { return ConceptCurrent{}, nil }
    client := getHTTPClient()
    params := map[string]string{
        "secid":  "90." + indexCode,
        "fields": "f57,f58,f106,f59,f43,f46,f60,f44,f45,f47,f48,f49,f113,f114,f115,f117,f85,f50,f119,f120,f121,f122,f135,f136,f137,f138,f139,f140,f141,f142,f143,f144,f145,f146,f147,f148,f149",
    }
    if wait > 0 { time.Sleep(wait) }
    resp, err := client.R().SetQueryParams(params).Get("https://push2.eastmoney.com/api/qt/stock/get")
    if err != nil { return ConceptCurrent{}, err }
    var data struct{ Data map[string]any `json:"data"` }
    if err := json.Unmarshal(resp.Body(), &data); err != nil { return ConceptCurrent{}, err }
    j := data.Data
    if j == nil { return ConceptCurrent{}, nil }
    code := toString(j["f57"])
    if code != indexCode { return ConceptCurrent{}, nil }
    preClose := parseF(toString(j["f60"]))
    cur := ConceptCurrent{IndexCode: indexCode}
    cur.Open = parseF(toString(j["f46"]))
    cur.High = parseF(toString(j["f44"]))
    cur.Low = parseF(toString(j["f45"]))
    cur.Price = parseF(toString(j["f43"]))
    cur.Volume = parseF(toString(j["f47"]))
    cur.Amount = parseF(toString(j["f48"]))
    cur.Change = cur.Price - preClose
    if preClose != 0 { cur.ChangePct = cur.Change / preClose * 100 }
    cur.TradeTime = time.Now().Format("2006-01-02 15:04:05")
    cur.TradeDate = time.Now().Format("2006-01-02")
    return cur, nil
}

func parseF(s string) float64 {
    s = strings.TrimSpace(strings.ReplaceAll(s, "%", ""))
    if s == "" || s == "--" { return 0 }
    v, _ := strconv.ParseFloat(s, 64)
    return v
}

func toString(v any) string { return strings.TrimSpace(fmt.Sprintf("%v", v)) }

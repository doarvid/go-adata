package indexmarket

import (
    "encoding/json"
    "fmt"
    "strings"
    "time"

)

func GetIndexDailyBaidu(indexCode string, startDate string, kType int, wait time.Duration) ([]IndexDailyBar, error) {
    if indexCode == "" { return []IndexDailyBar{}, nil }
    client := getHTTPClient()
    url := fmt.Sprintf("https://finance.pae.baidu.com/vapi/v1/getquotation?srcid=5353&all=1&pointType=string&group=quotation_index_kline&query=%s&code=%s&market_type=ab&newFormat=1&is_kc=0&ktype=day&finClientType=pc", indexCode, indexCode)
    var res struct {
        ResultCode string `json:"ResultCode"`
        Result     struct {
            NewMarketData struct {
                Keys       []string `json:"keys"`
                MarketData string   `json:"marketData"`
            } `json:"newMarketData"`
        } `json:"Result"`
    }
    for i := 0; i < 3; i++ {
        if wait > 0 { time.Sleep(wait) }
        resp, err := client.R().Get(url)
        if err != nil { return nil, err }
        if err := json.Unmarshal(resp.Body(), &res); err != nil { return nil, err }
        if res.ResultCode == "0" { break }
        time.Sleep(2 * time.Second)
    }
    if len(res.Result.NewMarketData.Keys) == 0 || res.Result.NewMarketData.MarketData == "" {
        return []IndexDailyBar{}, nil
    }
    keys := res.Result.NewMarketData.Keys
    raw := strings.Split(res.Result.NewMarketData.MarketData, ";")
    out := make([]IndexDailyBar, 0, len(raw))
    for _, row := range raw {
        cols := strings.Split(row, ",")
        if len(cols) != len(keys) { continue }
        kv := map[string]string{}
        for i := range keys { kv[keys[i]] = cols[i] }
        bar := IndexDailyBar{IndexCode: indexCode}
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
        bar.PreClose = parseF(kv["preClose"]) / 100
        if !(bar.Amount > 0 || bar.Volume > 0) { continue }
        out = append(out, bar)
    }
    if startDate != "" {
        s := strings.ReplaceAll(startDate, "-", "")
        filtered := make([]IndexDailyBar, 0, len(out))
        for _, it := range out {
            d := strings.ReplaceAll(it.TradeDate, "-", "")
            if d >= s { filtered = append(filtered, it) }
        }
        out = filtered
    }
    return out, nil
}

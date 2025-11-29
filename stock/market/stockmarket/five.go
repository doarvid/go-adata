package stockmarket

import (
    "encoding/json"
    "fmt"
    "time"

    httpc "go-adata/pkg/adata/common/http"
)

type Five struct {
    StockCode string  `json:"stock_code"`
    ShortName string  `json:"short_name"`
    S5        float64 `json:"s5"`
    Sv5       int64   `json:"sv5"`
    S4        float64 `json:"s4"`
    Sv4       int64   `json:"sv4"`
    S3        float64 `json:"s3"`
    Sv3       int64   `json:"sv3"`
    S2        float64 `json:"s2"`
    Sv2       int64   `json:"sv2"`
    S1        float64 `json:"s1"`
    Sv1       int64   `json:"sv1"`
    B1        float64 `json:"b1"`
    Bv1       int64   `json:"bv1"`
    B2        float64 `json:"b2"`
    Bv2       int64   `json:"bv2"`
    B3        float64 `json:"b3"`
    Bv3       int64   `json:"bv3"`
    B4        float64 `json:"b4"`
    Bv4       int64   `json:"bv4"`
    B5        float64 `json:"b5"`
    Bv5       int64   `json:"bv5"`
}

func GetMarketFiveBaidu(stockCode string, wait time.Duration) (Five, error) {
    client := httpc.NewClient()
    url := fmt.Sprintf("https://finance.pae.baidu.com/vapi/v1/getquotation?srcid=5353&all=1&pointType=string&group=quotation_minute_ab&query=%s&code=%s&market_type=ab&newFormat=1&finClientType=pc", stockCode, stockCode)
    var res struct{
        Result struct{
            Askinfos []map[string]any `json:"askinfos"`
            Buyinfos []map[string]any `json:"buyinfos"`
            Basicinfos struct{ Name string `json:"name"` } `json:"basicinfos"`
        } `json:"Result"`
    }
    if wait>0{ time.Sleep(wait)}
    resp, err := client.R().Get(url)
    if err != nil { return Five{}, err }
    if err := json.Unmarshal(resp.Body(), &res); err != nil { return Five{}, err }
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

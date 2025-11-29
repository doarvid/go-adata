package sentiment

import (
    "encoding/json"
    "strconv"
    "time"

    httpc "go-adata/pkg/adata/common/http"
)

type SecuritiesMarginRow struct {
    TradeDate string  `json:"trade_date"`
    Rzye      float64 `json:"rzye"`
    Rqye      float64 `json:"rqye"`
    Rzrqye    float64 `json:"rzrqye"`
    Rzrqyecz  float64 `json:"rzrqyecz"`
}

func SecuritiesMargin(startDate string, wait time.Duration) ([]SecuritiesMarginRow, error) {
    client := httpc.NewClient()
    totalPages := 1
    currPage := 1
    pageSize := 250
    startDateStr := startDate
    var start time.Time
    hasStart := false
    if startDate != "" { t, err := time.Parse("2006-01-02", startDate); if err == nil { start = t; hasStart = true } }
    out := make([]SecuritiesMarginRow, 0, 512)
    for currPage <= totalPages {
        url := "https://datacenter-web.eastmoney.com/api/data/v1/get?reportName=RPTA_RZRQ_LSHJ&columns=ALL&source=WEB&sortColumns=dim_date&sortTypes=-1&pageNumber=" + strconv.Itoa(currPage) + "&pageSize=" + strconv.Itoa(pageSize)
        if wait > 0 { time.Sleep(wait) }
        resp, err := client.R().Get(url)
        if err != nil { break }
        var res struct{ Success bool `json:"success"`; Result struct{ Pages int `json:"pages"`; Data []map[string]any `json:"data"` } `json:"result"` }
        if err := json.Unmarshal(resp.Body(), &res); err != nil { break }
        if !res.Success { break }
        if currPage == 1 { totalPages = res.Result.Pages }
        data := res.Result.Data
        for _, it := range data {
            dt, _ := time.Parse("2006-01-02 15:04:05", toString(it["DIM_DATE"]))
            out = append(out, SecuritiesMarginRow{
                TradeDate: dt.Format("2006-01-02"),
                Rzye: parseF(toString(it["RZYE"])),
                Rqye: parseF(toString(it["RQYE"])),
                Rzrqye: parseF(toString(it["RZRQYE"])),
                Rzrqyecz: parseF(toString(it["RZRQYECZ"])),
            })
        }
        if !hasStart { break }
        if hasStart {
            last := data[len(data)-1]
            dmin, _ := time.Parse("2006-01-02 15:04:05", toString(last["DIM_DATE"]))
            if !dmin.Before(start) { break }
        }
        currPage++
    }
    if startDateStr != "" {
        out2 := make([]SecuritiesMarginRow, 0, len(out))
        for _, r := range out { if r.TradeDate > startDateStr { out2 = append(out2, r) } }
        out = out2
    }
    return out, nil
}

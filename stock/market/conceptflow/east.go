package conceptflow

import (
    "encoding/json"
    "fmt"
    "strconv"
    "strings"
    "time"

    httpc "github.com/doarvid/go-adata/common/http"
)

type ConceptFlow struct {
    IndexCode          string  `json:"index_code"`
    IndexName          string  `json:"index_name"`
    ChangePct          float64 `json:"change_pct"`
    MainNetInflow      float64 `json:"main_net_inflow"`
    MainNetInflowRate  float64 `json:"main_net_inflow_rate"`
    MaxNetInflow       float64 `json:"max_net_inflow"`
    MaxNetInflowRate   float64 `json:"max_net_inflow_rate"`
    LgNetInflow        float64 `json:"lg_net_inflow"`
    LgNetInflowRate    float64 `json:"lg_net_inflow_rate"`
    MidNetInflow       float64 `json:"mid_net_inflow"`
    MidNetInflowRate   float64 `json:"mid_net_inflow_rate"`
    SmNetInflow        float64 `json:"sm_net_inflow"`
    SmNetInflowRate    float64 `json:"sm_net_inflow_rate"`
    StockCode          string  `json:"stock_code"`
    StockName          string  `json:"stock_name"`
}

func ListConceptCapitalFlowEast(daysType int, wait time.Duration) ([]ConceptFlow, error) {
    if daysType == 0 { daysType = 1 }
    fid, fields := buildParams(daysType)
    client := httpc.NewClient()
    max := 6
    sem := make(chan struct{}, max)
    type pageRes struct{ page int; rows []ConceptFlow }
    ch := make(chan pageRes, 50)
    for p := 1; p <= 50; p++ {
        page := p
        sem <- struct{}{}
        go func() {
            defer func() { <-sem }()
            url := "https://push2.eastmoney.com/api/qt/clist/get?fid=" + fid + "&po=1&pz=50&pn=" + toStr(page) + "&np=1&fltt=2&invt=2&fs=m:90 t:3&fields=" + fields
            if wait > 0 { time.Sleep(wait) }
            resp, err := client.R().Get(url)
            if err != nil { ch <- pageRes{page: page, rows: []ConceptFlow{}}; return }
            text := resp.String()
            if len(text) < 2 || resp.StatusCode() != 200 { ch <- pageRes{page: page, rows: []ConceptFlow{}}; return }
            l := strings.Index(text, "(") ; r := strings.LastIndex(text, ")")
            if l < 0 || r <= l { ch <- pageRes{page: page, rows: []ConceptFlow{}}; return }
            payload := text[l+1 : r]
            var j struct{ Data struct{ Diff []map[string]any `json:"diff"` } `json:"data"` }
            if err := json.Unmarshal([]byte(payload), &j); err != nil { ch <- pageRes{page: page, rows: []ConceptFlow{}}; return }
            if j.Data.Diff == nil || len(j.Data.Diff) == 0 { ch <- pageRes{page: page, rows: []ConceptFlow{}}; return }
            cols := strings.Split(fields, ",")
            rows := make([]ConceptFlow, 0, len(j.Data.Diff))
            for _, item := range j.Data.Diff {
                cf := ConceptFlow{}
                v := func(k string) string { return toString(item[k]) }
                cf.IndexCode = v(cols[0])
                cf.IndexName = v(cols[1])
                cf.ChangePct = parseF(v(cols[2]))
                cf.MainNetInflow = parseF(v(cols[3]))
                cf.MainNetInflowRate = parseF(v(cols[4]))
                cf.MaxNetInflow = parseF(v(cols[5]))
                cf.MaxNetInflowRate = parseF(v(cols[6]))
                cf.LgNetInflow = parseF(v(cols[7]))
                cf.LgNetInflowRate = parseF(v(cols[8]))
                cf.MidNetInflow = parseF(v(cols[9]))
                cf.MidNetInflowRate = parseF(v(cols[10]))
                cf.SmNetInflow = parseF(v(cols[11]))
                cf.SmNetInflowRate = parseF(v(cols[12]))
                cf.StockCode = v(cols[13])
                cf.StockName = v(cols[14])
                rows = append(rows, cf)
            }
            ch <- pageRes{page: page, rows: rows}
        }()
    }
    pages := make(map[int][]ConceptFlow, 50)
    for i := 0; i < 50; i++ {
        res := <-ch
        if len(res.rows) > 0 { pages[res.page] = res.rows }
    }
    out := make([]ConceptFlow, 0, 2000)
    for p := 1; p <= 50; p++ {
        if rows, ok := pages[p]; ok { out = append(out, rows...) }
    }
    return NormalizeConceptFlows(out), nil
}

func buildParams(daysType int) (string, string) {
    switch daysType {
    case 5:
        return "f164", "f12,f14,f109,f164,f165,f166,f167,f168,f169,f170,f171,f172,f173,f257,f258"
    case 10:
        return "f174", "f12,f14,f160,f174,f175,f176,f177,f178,f179,f180,f181,f182,f183,f260,f261"
    default:
        return "f62", "f12,f14,f3,f62,f184,f66,f69,f72,f75,f78,f81,f84,f87,f204,f205"
    }
}

func parseF(s string) float64 {
    s = strings.TrimSpace(strings.ReplaceAll(s, "%", ""))
    if s == "" || s == "-" || s == "--" { return 0 }
    v, _ := strconv.ParseFloat(s, 64)
    return v
}

func toStr(i int) string { return strconv.Itoa(i) }
func toString(v any) string { return strings.TrimSpace(fmt.Sprintf("%v", v)) }

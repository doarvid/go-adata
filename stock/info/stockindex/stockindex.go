package stockindex

import (
    "encoding/json"
    "fmt"
    "strconv"
    "time"

    httpc "github.com/doarvid/go-adata/common/http"
)

type IndexCode struct {
    IndexCode   string `json:"index_code"`
    ConceptCode string `json:"concept_code"`
    Name        string `json:"name"`
    Source      string `json:"source"`
}

type Constituent struct {
    IndexCode string `json:"index_code"`
    StockCode string `json:"stock_code"`
    ShortName string `json:"short_name"`
}

func AllIndexCodeEast(wait time.Duration) ([]IndexCode, error) {
    client := httpc.NewClient()
    var out []IndexCode
    for i := 0; i < 2; i++ {
        page := 1
        for page < 88 {
            var base string
            if i == 0 {
                base = "https://39.push2.eastmoney.com/api/qt/clist/get"
            } else {
                base = "https://31.push2.eastmoney.com/api/qt/clist/get"
            }
            params := map[string]string{
                "pn":    strconv.Itoa(page),
                "pz":    "20",
                "po":    "1",
                "np":    "1",
                "ut":    "bd1d9ddb04089700cf9c27f6f7426281",
                "fltt":  "2",
                "invt":  "2",
                "dect":  "1",
                "wbp2u": "|0|0|0|web",
                "fid":   "f3",
                "fs":    func() string { if i == 0 { return "m:1+s:2" } ; return "m:0+t:5" }(),
                "fields": "f12,f13,f14",
            }
            if wait > 0 { time.Sleep(wait) }
            resp, err := client.R().SetQueryParams(params).Get(base)
            if err != nil { return out, err }
            var data struct{
                Data struct{ Diff []map[string]any `json:"diff"` } `json:"data"`
            }
            if err := json.Unmarshal(resp.Body(), &data); err != nil { return out, err }
            if len(data.Data.Diff) == 0 { break }
            for _, d := range data.Data.Diff {
                idx := fmt.Sprintf("%v", d["f12"]) ; name := fmt.Sprintf("%v", d["f14"]) 
                out = append(out, IndexCode{IndexCode: idx, ConceptCode: "", Name: name, Source: "东方财富"})
            }
            page++
        }
    }
    return out, nil
}

func IndexConstituentBaidu(indexCode string, wait time.Duration) ([]Constituent, error) {
    client := httpc.NewClient()
    var out []Constituent
    for page := 0; page < 100; page++ {
        url := fmt.Sprintf("https://gushitong.baidu.com/opendata?resource_id=5352&query=%s&code=%s&market=ab&group=asyn_ranking&pn=%d&rn=100&pc_web=1&finClientType=pc", indexCode, indexCode, page*50)
        if wait > 0 { time.Sleep(wait) }
        resp, err := client.R().Get(url)
        if err != nil { return out, err }
        var data struct{
            ResultCode string `json:"ResultCode"`
            Result     []struct{ DisplayData struct{ ResultData struct{ TplData struct{ Result struct{ List []struct{ Code string `json:"code"`; Name string `json:"name"` } `json:"list"` } `json:"result"` } `json:"tplData"` } `json:"resultData"` } `json:"DisplayData"` } `json:"Result"`
        }
        if err := json.Unmarshal(resp.Body(), &data); err != nil { return out, err }
        if data.ResultCode != "0" || len(data.Result) == 0 { break }
        lists := data.Result[len(data.Result)-1].DisplayData.ResultData.TplData.Result.List
        if len(lists) == 0 { break }
        for _, it := range lists {
            out = append(out, Constituent{IndexCode: indexCode, StockCode: it.Code, ShortName: it.Name})
        }
    }
    // 去重
    seen := map[string]bool{}
    uniq := make([]Constituent, 0, len(out))
    for _, c := range out {
        k := c.IndexCode + ":" + c.StockCode
        if !seen[k] { seen[k] = true; uniq = append(uniq, c) }
    }
    return uniq, nil
}

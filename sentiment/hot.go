package sentiment

import (
    "encoding/json"
    "fmt"
    "strings"
    "time"

    httpc "go-adata/pkg/adata/common/http"
)

type PopRankRow struct {
	Rank      int     `json:"rank"`
	StockCode string  `json:"stock_code"`
	ShortName string  `json:"short_name"`
	Price     float64 `json:"price"`
	Change    float64 `json:"change"`
	ChangePct float64 `json:"change_pct"`
}

type HotRankRow struct {
	Rank       int     `json:"rank"`
	StockCode  string  `json:"stock_code"`
	ShortName  string  `json:"short_name"`
	ChangePct  float64 `json:"change_pct"`
	HotValue   float64 `json:"hot_value"`
	PopTag     string  `json:"pop_tag"`
	ConceptTag string  `json:"concept_tag"`
}

type HotConceptRow struct {
	Rank        int     `json:"rank"`
	ConceptCode string  `json:"concept_code"`
	ConceptName string  `json:"concept_name"`
	ChangePct   float64 `json:"change_pct"`
	HotValue    float64 `json:"hot_value"`
	HotTag      string  `json:"hot_tag"`
}

func PopRank100East(wait time.Duration) ([]PopRankRow, error) {
	client := httpc.NewClient()
    if wait > 0 { time.Sleep(wait) }
    resp, err := client.R().Post("https://emappdata.eastmoney.com/stockrank/getAllCurrentList")
    if err != nil {
        return nil, err
    }
    var res struct {
        Data []map[string]any `json:"data"`
    }
    if err := json.Unmarshal(resp.Body(), &res); err != nil {
        return nil, err
    }
	sc := make([]string, 0, len(res.Data))
	for _, it := range res.Data {
		sc = append(sc, toString(it["sc"]))
	}
	marks := make([]string, 0, len(sc))
	for _, item := range sc {
		if strings.HasPrefix(item, "SZ") {
			marks = append(marks, "0."+item[2:])
		} else {
			marks = append(marks, "1."+item[2:])
		}
	}
	q := strings.Join(marks, ",")
	url := "https://push2.eastmoney.com/api/qt/ulist.np/get?ut=f057cbcbce2a86e2866ab8877db1d059&fltt=2&invt=2&fields=f14,f3,f12,f2&secids=" + q
    if wait > 0 { time.Sleep(wait) }
    resp2, err := client.R().Get(url)
    if err != nil {
        return nil, err
    }
    var res2 struct {
        Data struct {
            Diff []map[string]any `json:"diff"`
        } `json:"data"`
    }
    if err := json.Unmarshal(resp2.Body(), &res2); err != nil {
        return nil, err
    }
	out := make([]PopRankRow, 0, len(res2.Data.Diff))
	rank := 1
	for _, it := range res2.Data.Diff {
		price := parseF(toString(it["f2"]))
		pct := parseF(toString(it["f3"]))
		out = append(out, PopRankRow{Rank: rank, StockCode: toString(it["f12"]), ShortName: toString(it["f14"]), Price: price, ChangePct: pct, Change: price * pct / 100})
		rank++
	}
	return out, nil
}

func HotRank100Ths(wait time.Duration) ([]HotRankRow, error) {
	client := httpc.NewClient()
	url := "https://dq.10jqka.com.cn/fuyao/hot_list_data/out/hot_list/v1/stock?stock_type=a&type=hour&list_type=normal"
    if wait > 0 { time.Sleep(wait) }
    resp, err := client.R().Get(url)
    if err != nil {
        return nil, err
    }
    var res struct {
        Data struct {
            StockList []map[string]any `json:"stock_list"`
        } `json:"data"`
    }
    if err := json.Unmarshal(resp.Body(), &res); err != nil {
        return nil, err
    }
	out := make([]HotRankRow, 0, len(res.Data.StockList))
	for _, d := range res.Data.StockList {
		conceptTags := []string{}
		if v, ok := d["tag"].(map[string]any); ok {
			if ct, ok2 := v["concept_tag"].([]any); ok2 {
				for _, x := range ct {
					conceptTags = append(conceptTags, toString(x))
				}
			}
		}
		popTag := ""
		if v, ok := d["tag"].(map[string]any); ok {
			if pt, ok2 := v["popularity_tag"].(string); ok2 {
				popTag = strings.ReplaceAll(pt, "\n", "")
			}
		}
		out = append(out, HotRankRow{
			Rank:       int(parseF(toString(d["order"]))),
			StockCode:  toString(d["code"]),
			ShortName:  toString(d["name"]),
			ChangePct:  parseF(toString(d["rise_and_fall"])),
			HotValue:   parseF(toString(d["rate"])),
			PopTag:     popTag,
			ConceptTag: strings.Join(conceptTags, ";"),
		})
	}
	return out, nil
}

func HotConcept20Ths(plateType int, wait time.Duration) ([]HotConceptRow, error) {
	client := httpc.NewClient()
	t := "concept"
	if plateType != 1 {
		t = "industry"
	}
	url := fmt.Sprintf("https://dq.10jqka.com.cn/fuyao/hot_list_data/out/hot_list/v1/plate?type=%s", t)
    if wait > 0 { time.Sleep(wait) }
    resp, err := client.R().Get(url)
    if err != nil {
        return nil, err
    }
    var res struct {
        Data struct {
            PlateList []map[string]any `json:"plate_list"`
        } `json:"data"`
    }
    if err := json.Unmarshal(resp.Body(), &res); err != nil {
        return nil, err
    }
	out := make([]HotConceptRow, 0, len(res.Data.PlateList))
	for _, d := range res.Data.PlateList {
		out = append(out, HotConceptRow{
			Rank:        int(parseF(toString(d["order"]))),
			ConceptCode: toString(d["code"]),
			ConceptName: toString(d["name"]),
			ChangePct:   parseF(toString(d["rise_and_fall"])),
			HotValue:    parseF(toString(d["rate"])),
			HotTag:      toString(d["hot_tag"]),
		})
	}
	return out, nil
}

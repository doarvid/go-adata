package sentiment

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	httpc "github.com/doarvid/go-adata/common/http"
)

type PopRankRow struct {
	Rank      int     `json:"rank"`       // 排名，示例：100
	StockCode string  `json:"stock_code"` // 股票代码，示例：000799
	ShortName string  `json:"short_name"` // 股票简称，示例：酒鬼酒
	Price     float64 `json:"price"`      // 最新价格，示例：88
	Change    float64 `json:"change"`     // 涨跌额，示例：2
	ChangePct float64 `json:"change_pct"` // 涨跌幅（%），示例：10.0020
}

type HotRankRow struct {
	Rank       int     `json:"rank"`        // 排名，示例：100
	StockCode  string  `json:"stock_code"`  // 股票代码，示例：000799
	ShortName  string  `json:"short_name"`  // 股票简称，示例：酒鬼酒
	ChangePct  float64 `json:"change_pct"`  // 涨跌幅（%），示例：10.002
	HotValue   float64 `json:"hot_value"`   // 热度值，示例：432509.0
	PopTag     string  `json:"pop_tag"`     // 人气标签，示例：首板涨停
	ConceptTag string  `json:"concept_tag"` // 概念板块，示例：白酒概念;国企改革
}

type HotConceptRow struct {
	Rank        int     `json:"rank"`         // 排名，示例：1
	ConceptCode string  `json:"concept_code"` // 概念代码，示例：881157
	ConceptName string  `json:"concept_name"` // 概念名称，示例：证券
	ChangePct   float64 `json:"change_pct"`   // 涨跌幅（%），示例：0.2488
	HotValue    float64 `json:"hot_value"`    // 热度值，示例：1130204.5
	HotTag      string  `json:"hot_tag"`      // 热度标签，示例：连续351天上榜
}

// 东方财富人气榜100
// http://guba.eastmoney.com/rank/
func PopRank100East(wait time.Duration) ([]PopRankRow, error) {
	client := httpc.NewClient()
	if wait > 0 {
		time.Sleep(wait)
	}

	params := map[string]any{
		"appId":      "appId01",
		"globalId":   "786e4c21-70dc-435a-93bb-38",
		"marketType": "",
		"pageNo":     1,
		"pageSize":   100,
	}
	resp, err := client.R().SetBody(params).Post("https://emappdata.eastmoney.com/stockrank/getAllCurrentList")
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
	if wait > 0 {
		time.Sleep(wait)
	}
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
	if wait > 0 {
		time.Sleep(wait)
	}
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

type PlateType string

const (
	PlateTypeConcept  PlateType = "concept"
	PlateTypeIndustry PlateType = "industry"
)

func HotConcept20Ths(plateType PlateType, wait time.Duration) ([]HotConceptRow, error) {
	client := httpc.NewClient()
	if plateType != PlateTypeConcept && plateType != PlateTypeIndustry {
		return nil, fmt.Errorf("invalid plate type: %s", plateType)
	}
	t := string(plateType)
	url := fmt.Sprintf("https://dq.10jqka.com.cn/fuyao/hot_list_data/out/hot_list/v1/plate?type=%s", t)
	if wait > 0 {
		time.Sleep(wait)
	}
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

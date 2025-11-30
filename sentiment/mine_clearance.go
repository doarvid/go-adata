package sentiment

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	httpc "github.com/doarvid/go-adata/common/http"
)

type MineRow struct {
	StockCode string  `json:"stock_code"`
	ShortName string  `json:"short_name"`
	Score     float64 `json:"score"`
	FType     string  `json:"f_type"`
	SType     string  `json:"s_type"`
	TType     string  `json:"t_type"`
	Reason    string  `json:"reason"`
}

func MineClearanceTDX(stockCode string, wait time.Duration) ([]MineRow, error) {
	if stockCode == "" {
		return nil, fmt.Errorf("stock code is empty")
	}
	client := httpc.NewClient()
	url := "http://page3.tdx.com.cn:7615/site/pcwebcall_static/bxb/json/" + stockCode + ".json"
	if wait > 0 {
		time.Sleep(wait)
	}
	resp, err := client.R().Get(url)
	if err != nil {
		return []MineRow{{StockCode: "", ShortName: "", Score: 0, FType: "暂无数据"}}, nil
	}
	var res struct {
		Name string           `json:"name"`
		Data []map[string]any `json:"data"`
	}
	if err := json.Unmarshal(resp.Body(), &res); err != nil {
		return []MineRow{{StockCode: "", ShortName: "", Score: 0, FType: "暂无数据"}}, nil
	}
	name := res.Name
	data := res.Data
	out := make([]MineRow, 0, 128)
	score := 100.0
	sTypeDeduct := map[string]bool{}
	for _, i := range data {
		ftype := toString(i["name"])
		rows, _ := i["rows"].([]any)
		for _, k := range rows {
			kk, _ := k.(map[string]any)
			if toString(kk["trigyy"]) != "" {
				com, _ := kk["commonlxid"].([]any)
				if len(com) == 0 {
					out = append(out, MineRow{StockCode: stockCode, ShortName: name, FType: ftype, SType: toString(kk["lx"]), TType: "", Reason: toString(kk["trigyy"]), Score: parseF(toString(kk["fs"]))})
					if toInt(kk["trig"]) == 1 {
						score -= parseF(toString(kk["fs"]))
					}
				}
				for _, j := range com {
					jj, _ := j.(map[string]any)
					if toString(jj["trigyy"]) != "" {
						out = append(out, MineRow{StockCode: stockCode, ShortName: name, FType: ftype, SType: toString(kk["lx"]), TType: toString(jj["lx"]), Reason: toString(jj["trigyy"]), Score: parseF(toString(jj["fs"]))})
						if toInt(jj["trig"]) == 1 && !sTypeDeduct[toString(kk["lx"])] {
							score -= parseF(toString(jj["fs"]))
							sTypeDeduct[toString(kk["lx"])] = true
						}
					}
				}
			}
		}
	}
	if len(out) == 0 {
		if strings.HasSuffix(name, "退") {
			return []MineRow{{StockCode: stockCode, ShortName: name, Score: -1, FType: "已退市"}}, nil
		}
		if score < 1 {
			score = 1
		}
		return []MineRow{{StockCode: stockCode, ShortName: name, Score: score, FType: "暂无风险项"}}, nil
	}
	if score < 1 {
		score = 1
	}
	for i := range out {
		out[i].Score = score
	}
	return out, nil
}

func toInt(v any) int {
	s := strings.TrimSpace(fmt.Sprintf("%v", v))
	if s == "" {
		return 0
	}
	n, _ := strconv.Atoi(s)
	return n
}

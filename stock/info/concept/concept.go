package concept

import (
    "encoding/csv"
    "encoding/json"
    "fmt"
    "io"
    "os"
    "strconv"
    "strings"
    "time"

    "github.com/doarvid/go-adata/common/codeutils"
    httpc "github.com/doarvid/go-adata/common/http"
    "github.com/doarvid/go-adata/stock/cache"
)

type ConceptCode struct {
	ConceptCode string `json:"concept_code"`
	IndexCode   string `json:"index_code"`
	Name        string `json:"name"`
	Source      string `json:"source"`
}

type ConceptInfo struct {
	StockCode   string `json:"stock_code"`
	ConceptCode string `json:"concept_code"`
	Name        string `json:"name"`
	Source      string `json:"source"`
	Reason      string `json:"reason"`
}

type Constituent struct {
	StockCode string `json:"stock_code"`
	ShortName string `json:"short_name"`
}

func LoadAllConceptCodesFromCSV() ([]ConceptCode, error) {
	p := cache.GetAllConceptCodeEastCSVPath()
	f, err := os.Open(p)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	r := csv.NewReader(f)
	if _, err := r.Read(); err != nil {
		return nil, err
	}
	var out []ConceptCode
	for {
		rec, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		out = append(out, ConceptCode{ConceptCode: rec[0], IndexCode: rec[1], Name: rec[2], Source: rec[3]})
	}
	return out, nil
}

func AllConceptCodesEast(wait time.Duration) ([]ConceptCode, error) {
	client := httpc.NewClient()
	page := 1
	size := 100
	var out []ConceptCode
    for page < 50 {
        params := map[string]string{
            "pn":     strconv.Itoa(page),
            "pz":     strconv.Itoa(size),
            "po":     "1",
            "np":     "1",
            "fields": "f12,f13,f14,f62",
            "fid":    "f62",
            "fs":     "m:90+t:3",
        }
        if wait > 0 {
            time.Sleep(wait)
        }
        resp, err := client.R().SetQueryParams(params).Get("https://push2.eastmoney.com/api/qt/clist/get")
        if err != nil {
            return out, err
        }
        var data struct {
            Data struct {
                Diff []map[string]any `json:"diff"`
            } `json:"data"`
        }
        if err := json.Unmarshal(resp.Body(), &data); err != nil {
            return out, err
        }
		if len(data.Data.Diff) == 0 {
			break
		}
		for _, d := range data.Data.Diff {
			name := strings.TrimSpace(toString(d["f14"]))
			code := toString(d["f12"]) // BKxxxx
			out = append(out, ConceptCode{ConceptCode: code, IndexCode: code, Name: name, Source: "东方财富"})
		}
		if len(data.Data.Diff) < size {
			break
		}
		page++
	}
	// 合并缓存去重
	cached, _ := LoadAllConceptCodesFromCSV()
	seen := map[string]bool{}
	for _, c := range out {
		seen[c.ConceptCode] = true
	}
	for _, c := range cached {
		if !seen[c.ConceptCode] {
			out = append(out, c)
		}
	}
	return out, nil
}

func GetConceptEast(stockCode string, wait time.Duration) ([]ConceptInfo, error) {
	client := httpc.NewClient()
	sc := codeutils.CompileExchangeByStockCode(stockCode)
    params := map[string]string{
        "reportName":  "RPT_F10_CORETHEME_BOARDTYPE",
        "columns":     "SECUCODE,SECURITY_CODE,SECURITY_NAME_ABBR,NEW_BOARD_CODE,BOARD_NAME,SELECTED_BOARD_REASON,IS_PRECISE,BOARD_RANK,BOARD_YIELD,DERIVE_BOARD_CODE",
        "quoteColumns": "f3~05~NEW_BOARD_CODE~BOARD_YIELD",
        "filter":      "(SECUCODE=\"" + sc + "\")(IS_PRECISE=\"1\")",
        "pageNumber":  "1",
        "pageSize":    "50",
        "sortTypes":   "1",
        "sortColumns": "BOARD_RANK",
        "source":      "HSF10",
        "client":      "PC",
    }
    if wait > 0 {
        time.Sleep(wait)
    }
    resp, err := client.R().SetQueryParams(params).Get("https://datacenter.eastmoney.com/securities/api/data/v1/get")
    if err != nil {
        return nil, err
    }
    var data struct {
        Result struct {
            Data []map[string]any `json:"data"`
        } `json:"result"`
    }
    if err := json.Unmarshal(resp.Body(), &data); err != nil {
        return nil, err
    }
	var out []ConceptInfo
	for _, d := range data.Result.Data {
		out = append(out, ConceptInfo{
			StockCode:   stockCode,
			ConceptCode: toString(d["NEW_BOARD_CODE"]),
			Name:        toString(d["BOARD_NAME"]),
			Source:      "东方财富",
			Reason:      toString(d["SELECTED_BOARD_REASON"]),
		})
	}
	return out, nil
}

func ConceptConstituentEast(conceptCode string, wait time.Duration) ([]Constituent, error) {
	client := httpc.NewClient()
	var out []Constituent
	page := 1
    for page < 100 {
        params := map[string]string{
            "fid":    "f62",
            "po":     "1",
            "pz":     "200",
            "pn":     strconv.Itoa(page),
            "np":     "1",
            "fltt":   "2",
            "invt":   "2",
            "fs":     "b:" + conceptCode,
            "fields": "f12,f14",
        }
        if wait > 0 {
            time.Sleep(wait)
        }
        resp, err := client.R().SetQueryParams(params).Get("https://push2.eastmoney.com/api/qt/clist/get")
        if err != nil {
            return out, err
        }
        var data struct {
            Data struct {
                Diff []map[string]any `json:"diff"`
            } `json:"data"`
        }
        if err := json.Unmarshal(resp.Body(), &data); err != nil {
            return out, err
        }
		if len(data.Data.Diff) == 0 {
			break
		}
		for _, d := range data.Data.Diff {
			out = append(out, Constituent{StockCode: toString(d["f12"]), ShortName: toString(d["f14"])})
		}
		page++
	}
	return out, nil
}

func toString(v any) string { return strings.TrimSpace(fmt.Sprintf("%v", v)) }

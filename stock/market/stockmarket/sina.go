package stockmarket

import (
    "strings"
    "time"

    "github.com/doarvid/go-adata/common/codeutils"
    httpc "github.com/doarvid/go-adata/common/http"
)

type CurrentQuote struct {
	StockCode string  `json:"stock_code"`
	ShortName string  `json:"short_name"`
	Price     float64 `json:"price"`
	Change    float64 `json:"change"`
	ChangePct float64 `json:"change_pct"`
	Volume    float64 `json:"volume"`
	Amount    float64 `json:"amount"`
}

func ListMarketCurrentSina(codeList []string, wait time.Duration) ([]CurrentQuote, error) {
	if len(codeList) == 0 {
		return []CurrentQuote{}, nil
	}
	client := httpc.NewClient()
	api := "https://hq.sinajs.cn/list="
	for _, code := range codeList {
		ex := strings.ToLower(codeutils.GetExchangeByStockCode(code))
		api += "s_" + ex + code + ","
	}
    if wait > 0 {
        time.Sleep(wait)
    }
    headers := map[string]string{"Referer": "https://finance.sina.com.cn/", "User-Agent": "Mozilla/5.0"}
    resp, err := client.R().SetHeaders(headers).Get(api)
    if err != nil {
        return nil, err
    }
    text := resp.String()
    if len(text) < 1 || resp.StatusCode() != 200 {
        return []CurrentQuote{}, nil
    }
	parts := strings.Split(text, ";")
	out := make([]CurrentQuote, 0, len(parts))
	for _, p := range parts {
		if len(p) < 8 {
			continue
		}
		idx := strings.Index(p, "=")
		if idx <= 0 || idx-6 < 0 {
			continue
		}
		code := p[idx-6 : idx]
		vals := strings.Split(p[idx+2:len(p)-1], ",")
		if len(vals) != 6 {
			continue
		}
		out = append(out, CurrentQuote{
			StockCode: code,
			ShortName: vals[0],
			Price:     parseF(vals[1]),
			Change:    parseF(vals[2]),
			ChangePct: parseF(vals[3]),
			Volume:    parseF(vals[4]) * 100,   // 北京返回手，换算成股
			Amount:    parseF(vals[5]) * 10000, // 北京返回万元，换算为元
		})
	}
	return out, nil
}

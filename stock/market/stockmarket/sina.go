package stockmarket

import (
	"context"
	"strings"
	"time"

	"github.com/doarvid/go-adata/common/codeutils"
	"github.com/doarvid/go-adata/common/utils"
)

type CurrentQuote struct {
	StockCode   string   `json:"stock_code"`    // 股票代码
	ShortName   string   `json:"short_name"`    // 股票简称
	Price       float64  `json:"price"`         // 最新价/当前价
	Change      float64  `json:"change"`        // 涨跌额
	ChangePct   float64  `json:"change_pct"`    // 涨跌幅 (%)
	Volume      float64  `json:"volume"`        // 成交量 (股)
	Amount      float64  `json:"amount"`        // 成交额 (元)
	TotalCap    *float64 `json:"total_cap"`     // 总市值 (元)
	FloatCap    *float64 `json:"float_cap"`     // 流通市值 (元)
	TotalShares *float64 `json:"total_shares"`  // 总股本
	FloatShares *float64 `json:"float_shares"`  // 流通股本
	Change5d    *float64 `json:"change_5d"`     // 5 日涨跌幅 (%)
	Change10d   *float64 `json:"change_10d"`    // 10 日涨跌幅 (%)
	Change20d   *float64 `json:"change_20d"`    // 20 日涨跌幅 (%)
	Change60d   *float64 `json:"change_60d"`    // 60 日涨跌幅 (%)
	Change120d  *float64 `json:"change_120d"`   // 120 日涨跌幅 (%)
	Change250d  *float64 `json:"change_250d"`   // 250 日涨跌幅 (%)
	ChangeYtd   *float64 `json:"change_ytd"`    // 年初至今涨跌幅 (%)
}

func (m *Market) ListCurrentSina(ctx context.Context, codeList []string) ([]CurrentQuote, error) {
	if len(codeList) == 0 {
		return []CurrentQuote{}, nil
	}
	client := m.client
	api := "https://hq.sinajs.cn/list="
	for _, code := range codeList {
		ex := strings.ToLower(codeutils.GetExchangeByStockCode(code))
		api += "s_" + ex + code + ","
	}
	if m.MinWait > 0 {
		time.Sleep(m.MinWait)
	}
	headers := map[string]string{"Referer": "https://finance.sina.com.cn/", "User-Agent": "Mozilla/5.0"}
	resp, err := client.R().SetContext(ctx).SetHeaders(headers).Get(api)
	if err != nil {
		return nil, err
	}
	text, err := utils.GBKToUTF8([]byte(resp.String()))
	if err != nil {
		return nil, err
	}
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

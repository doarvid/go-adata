package stockmarket

import (
	"context"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/doarvid/go-adata/common/codeutils"
)

func (m *Market) ListCurrentQQ(ctx context.Context, codeList []string) ([]CurrentQuote, error) {
	if len(codeList) == 0 {
		return []CurrentQuote{}, nil
	}
	client := m.client
	api := "https://qt.gtimg.cn/q="
	for _, code := range codeList {
		ex := strings.ToLower(codeutils.GetExchangeByStockCode(code))
		api += "s_" + ex + code + ","
	}
	if m.MinWait > 0 {
		time.Sleep(m.MinWait)
	}
	resp, err := client.R().SetContext(ctx).Get(api)
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
		arr := strings.Split(p, "~")
		if len(arr) != 11 {
			continue
		}
		out = append(out, CurrentQuote{
			ShortName: arr[1],
			StockCode: arr[2],
			Price:     parseF(arr[3]),
			Change:    parseF(arr[4]),
			ChangePct: parseF(arr[5]),
			Volume:    parseF(arr[6]) * 100,
			Amount:    parseF(arr[7]) * 10000,
		})
	}
	return out, nil
}

func (m *Market) GetFiveQQ(ctx context.Context, stockCode string) (Five, error) {
	list, err := m.ListFiveQQ(ctx, []string{stockCode})
	if err != nil || len(list) == 0 {
		return Five{}, err
	}
	return list[0], nil
}

func (m *Market) ListFiveQQ(ctx context.Context, codeList []string) ([]Five, error) {
	if len(codeList) == 0 {
		return []Five{}, nil
	}
	client := m.client
	api := "https://web.sqt.gtimg.cn/q="
	for _, code := range codeList {
		ex := strings.ToLower(codeutils.GetExchangeByStockCode(code))
		api += ex + code + ","
	}
	if m.MinWait > 0 {
		time.Sleep(m.MinWait)
	}
	resp, err := client.R().SetContext(ctx).Get(api)
	if err != nil {
		return nil, err
	}
	text := resp.String()
	if len(text) < 1 || resp.StatusCode() != 200 {
		return []Five{}, nil
	}
	parts := strings.Split(text, ";")
	out := make([]Five, 0, len(parts))
	for _, p := range parts {
		if len(p) < 8 {
			continue
		}
		arr := strings.Split(p, "~")
		if len(arr) < 85 {
			continue
		}
		code := arr[2]
		name := arr[1]
		f := Five{StockCode: code, ShortName: name}
		f.B1 = parseF(arr[27])
		f.Bv1 = toInt64(arr[28]) * 100
		f.B2 = parseF(arr[25])
		f.Bv2 = toInt64(arr[26]) * 100
		f.B3 = parseF(arr[23])
		f.Bv3 = toInt64(arr[24]) * 100
		f.B4 = parseF(arr[21])
		f.Bv4 = toInt64(arr[22]) * 100
		f.B5 = parseF(arr[19])
		f.Bv5 = toInt64(arr[20]) * 100
		f.S5 = parseF(arr[9])
		f.Sv5 = toInt64(arr[10]) * 100
		f.S4 = parseF(arr[11])
		f.Sv4 = toInt64(arr[12]) * 100
		f.S3 = parseF(arr[13])
		f.Sv3 = toInt64(arr[14]) * 100
		f.S2 = parseF(arr[15])
		f.Sv2 = toInt64(arr[16]) * 100
		f.S1 = parseF(arr[17])
		f.Sv1 = toInt64(arr[18]) * 100
		out = append(out, f)
	}
	return out, nil
}

func (m *Market) GetBarQQ(ctx context.Context, stockCode string) ([]TickBar, error) {
	client := m.client
	ex := strings.ToLower(codeutils.GetExchangeByStockCode(stockCode))
	code := ex + stockCode
	out := make([]TickBar, 0, 512)
	re := regexp.MustCompile(`\[\s*\d+\s*,\s*"([^"]+)"`)
	for page := 0; page < 10000; page++ {
		params := map[string]string{
			"appn":   "detail",
			"action": "data",
			"c":      code,
			"p":      strconv.Itoa(page),
		}
		if m.MinWait > 0 {
			time.Sleep(m.MinWait)
		}
		resp, err := client.R().SetContext(ctx).SetQueryParams(params).Get("http://stock.gtimg.cn/data/index.php")
		if err != nil {
			break
		}
		text := resp.String()
		m := re.FindStringSubmatch(text)
		if len(m) < 2 {
			break
		}
		payload := m[1]
		rows := strings.Split(payload, "|")
		if len(rows) == 0 {
			break
		}
		for _, r := range rows {
			cols := strings.Split(r, "/")
			// expect: no/time/price/x/volume/x/bs_type
			if len(cols) < 7 {
				continue
			}
			out = append(out, TickBar{
				TradeTime: cols[1],
				Price:     parseF(cols[2]),
				Volume:    toInt64(cols[4]) * 100,
				BsType:    cols[6],
				StockCode: stockCode,
			})
		}
	}
	return out, nil
}

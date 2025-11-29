package sentiment

import (
    "regexp"
    "strconv"
    "strings"
    "time"

    httpc "go-adata/pkg/adata/common/http"
)

type StockLiftingRow struct {
	StockCode string  `json:"stock_code"`
	ShortName string  `json:"short_name"`
	LiftDate  string  `json:"lift_date"`
	Volume    int64   `json:"volume"`
	Amount    int64   `json:"amount"`
	Ratio     float64 `json:"ratio"`
	Price     float64 `json:"price"`
}

func StockLiftingLastMonth(wait time.Duration) ([]StockLiftingRow, error) {
	client := httpc.NewClient()
    url := "http://data.10jqka.com.cn/market/xsjj/field/enddate/order/desc/ajax/1/free/1/"
    if wait > 0 {
        time.Sleep(wait)
    }
    resp, err := client.R().Get(url)
    if err != nil {
        return []StockLiftingRow{}, nil
    }
    text := resp.String()
	if !(strings.Contains(text, "解禁日期") || strings.Contains(text, "解禁股")) {
		return []StockLiftingRow{}, nil
	}
	rows := parseTableRows(text)
	out := make([]StockLiftingRow, 0, len(rows))
	for _, r := range rows {
		if len(r) < 8 {
			continue
		}
		vol := toUnitInt(r[3])
		amt := toUnitInt(r[5])
		out = append(out, StockLiftingRow{
			StockCode: r[0],
			ShortName: r[1],
			LiftDate:  r[2],
			Volume:    vol,
			Ratio:     parseF(r[6]),
			Price:     parseF(r[4]),
			Amount:    amt,
		})
	}
	return out, nil
}

func parseTableRows(html string) [][]string {
	trs := strings.Split(html, "<tr")
	out := [][]string{}
	for i, seg := range trs {
		if i == 0 {
			continue
		}
		cols := []string{}
		re := regexp.MustCompile(`>([^<]+)</td>`)
		for _, m := range re.FindAllStringSubmatch(seg, -1) {
			cols = append(cols, strings.TrimSpace(m[1]))
		}
		if len(cols) >= 8 {
			out = append(out, []string{cols[1], cols[2], cols[3], cols[4], cols[5], cols[6], cols[7], cols[8]})
		}
	}
	return out
}

func toUnitInt(s string) int64 {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0
	}
	if strings.HasSuffix(s, "万") {
		v, _ := strconv.ParseFloat(strings.TrimSuffix(s, "万"), 64)
		return int64(v * 10000)
	}
	if strings.HasSuffix(s, "亿") {
		v, _ := strconv.ParseFloat(strings.TrimSuffix(s, "亿"), 64)
		return int64(v * 100000000)
	}
	v, _ := strconv.ParseFloat(s, 64)
	return int64(v)
}

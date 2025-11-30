package sentiment

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/doarvid/go-adata/common/header"
	httpc "github.com/doarvid/go-adata/common/http"
	"github.com/doarvid/go-adata/common/utils"
)

// StockLiftingRow 股票解禁数据行
type StockLiftingRow struct {
	StockCode string  `json:"stock_code"` // 股票代码，示例：300539
	ShortName string  `json:"short_name"` // 股票简称，示例：横河精密
	LiftDate  string  `json:"lift_date"`  // 解禁日期，示例：2023-06-05
	Volume    int64   `json:"volume"`     // 解禁股数(股)，示例：123400
	Amount    int64   `json:"amount"`     // 当前解禁市值(元)，根据当前价格计算，示例：1621596
	Ratio     float64 `json:"ratio"`      // 占总股本比例(%)，示例：0.36
	Price     float64 `json:"price"`      // 当前价格(元)，示例：13.14
}

// 获取最近一个月的股票解禁数据
// 帮助提前规避解禁股票，对于大额解禁个股参考意义巨大。
func StockLiftingLastMonth(wait time.Duration) ([]StockLiftingRow, error) {
	client := httpc.NewClient()
	url := "http://data.10jqka.com.cn/market/xsjj/field/enddate/order/desc/ajax/1/free/1/"
	if wait > 0 {
		time.Sleep(wait)
	}
	headers := header.DefaultHeaders()
	headers["Host"] = "data.10jqka.com.cn"
	headers["Cookie"] = utils.ThsCookie()
	resp, err := client.R().SetHeaders(headers).Get(url)
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

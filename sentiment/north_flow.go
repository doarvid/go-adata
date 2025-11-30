package sentiment

import (
	"encoding/json"
	"io"
	"strings"
	"time"

	httpc "github.com/doarvid/go-adata/common/http"
	"github.com/doarvid/go-adata/stock/info/tradecalendar"
)

type NorthFlowDaily struct {
	// 交易时间，格式如 2023-06-01
	TradeDate string `json:"trade_date"`
	// 沪港通净买入金额（元），买入和卖出合计，示例：405050400
	NetHgt float64 `json:"net_hgt"`
	// 沪港通买入金额（元）
	BuyHgt float64 `json:"buy_hgt"`
	// 沪港通卖出金额（元）
	SellHgt float64 `json:"sell_hgt"`
	// 深港通净买入金额（元），买入和卖出合计，示例：151704400
	NetSgt float64 `json:"net_sgt"`
	// 深港通买入金额（元）
	BuySgt float64 `json:"buy_sgt"`
	// 深港通卖出金额（元）
	SellSgt float64 `json:"sell_sgt"`
	// 北向净买入金额（元），沪港通和深港通合计，示例：556754800
	NetTgt float64 `json:"net_tgt"`
	// 北向买入金额（元）
	BuyTgt float64 `json:"buy_tgt"`
	// 北向卖出金额（元）
	SellTgt float64 `json:"sell_tgt"`
}

type NorthFlowMinute struct {
	// 交易时间，格式如 2023-06-01 09:30:00
	TradeTime string `json:"trade_time"`
	// 沪港通净买入金额（元），示例：405050400
	NetHgt float64 `json:"net_hgt"`
	// 深港通净买入金额（元），示例：151704400
	NetSgt float64 `json:"net_sgt"`
	// 北向净买入金额（元），沪港通和深港通合计，示例：556754800
	NetTgt float64 `json:"net_tgt"`
}

// 获取北向的历史流入行情
// startDate 开始日期，格式为"2006-01-02"
// wait 等待时间，单位为秒
func NorthFlow(startDate string, wait time.Duration) ([]NorthFlowDaily, error) {
	return northFlowEast(startDate, wait)
}

func northFlowEast(startDate string, wait time.Duration) ([]NorthFlowDaily, error) {
	client := httpc.NewClient()
	currPage := 1
	out := make([]NorthFlowDaily, 0, 1024)
	var start time.Time
	var hasStart bool
	if startDate != "" {
		t, err := time.Parse("2006-01-02", startDate)
		if err == nil {
			start = t
			hasStart = true
		}
		dateMin, _ := time.Parse("2006-01-02", "2017-01-01")
		if start.Before(dateMin) {
			start = dateMin
		}
	}
	for currPage < 18 {
		base := "https://datacenter-web.eastmoney.com/api/data/v1/get?sortColumns=TRADE_DATE&sortTypes=-1&pageSize=1000&pageNumber=" + toString(currPage) + "&reportName=RPT_MUTUAL_DEAL_HISTORY&columns=ALL&source=WEB&client=WEB&"
		sgtURL := base + "filter=(MUTUAL_TYPE=%27001%27)"
		hgtURL := base + "filter=(MUTUAL_TYPE=%27003%27)"
		if wait > 0 {
			time.Sleep(wait)
		}
		resp1, err1 := client.R().Get(sgtURL)
		if err1 != nil {
			break
		}
		resp2, err2 := client.R().Get(hgtURL)
		if err2 != nil {
			break
		}
		buf1 := new(strings.Builder)
		buf2 := new(strings.Builder)
		if _, err := ioCopy(buf1, strings.NewReader(resp1.String())); err != nil {
			break
		}
		if _, err := ioCopy(buf2, strings.NewReader(resp2.String())); err != nil {
			break
		}
		sgtText := strings.ReplaceAll(buf1.String(), "null", "0")
		hgtText := strings.ReplaceAll(buf2.String(), "null", "0")
		l1 := strings.Index(sgtText, "{")
		l2 := strings.Index(hgtText, "{")
		if l1 < 0 || l2 < 0 {
			break
		}
		var sgtRes struct {
			Result struct {
				Data []map[string]any `json:"data"`
			} `json:"result"`
		}
		var hgtRes struct {
			Result struct {
				Data []map[string]any `json:"data"`
			} `json:"result"`
		}
		if err := json.Unmarshal([]byte(sgtText), &sgtRes); err != nil {
			break
		}
		if err := json.Unmarshal([]byte(hgtText), &hgtRes); err != nil {
			break
		}
		sgtData := sgtRes.Result.Data
		hgtData := hgtRes.Result.Data
		if len(sgtData) == 0 {
			break
		}
		isEnd := false
		for i := range hgtData {
			if !hasStart && i >= 30 {
				isEnd = true
				break
			}
			dt, _ := time.Parse("2006-01-02 15:04:05", toString(hgtData[i]["TRADE_DATE"]))
			if hasStart && start.After(dt) {
				isEnd = true
				break
			}
			nh := parseF(toString(hgtData[i]["NET_DEAL_AMT"])) * 1000000
			bh := parseF(toString(hgtData[i]["BUY_AMT"])) * 1000000
			sh := parseF(toString(hgtData[i]["SELL_AMT"])) * 1000000
			ns := parseF(toString(sgtData[i]["NET_DEAL_AMT"])) * 1000000
			bs := parseF(toString(sgtData[i]["BUY_AMT"])) * 1000000
			ss := parseF(toString(sgtData[i]["SELL_AMT"])) * 1000000
			out = append(out, NorthFlowDaily{
				TradeDate: dt.Format("2006-01-02"),
				NetHgt:    nh, BuyHgt: bh, SellHgt: sh,
				NetSgt: ns, BuySgt: bs, SellSgt: ss,
				NetTgt: nh + ns, BuyTgt: bh + bs, SellTgt: sh + ss,
			})
		}
		if isEnd {
			break
		}
		currPage++
	}
	return out, nil
}

func NorthFlowMin(wait time.Duration) ([]NorthFlowMinute, error) {
	r, _ := northFlowMinEast(wait)
	if len(r) == 0 {
		r, _ = northFlowMinThs(wait)
	}
	return r, nil
}

func NorthFlowCurrent(wait time.Duration) (NorthFlowMinute, error) {
	mins, _ := NorthFlowMin(wait)
	if len(mins) == 0 {
		return NorthFlowMinute{}, nil
	}
	return mins[len(mins)-1], nil
}

func northFlowMinThs(wait time.Duration) ([]NorthFlowMinute, error) {
	client := httpc.NewClient()
	if wait > 0 {
		time.Sleep(wait)
	}
	resp, err := client.R().Get("https://data.hexin.cn/market/hsgtApi/method/dayChart/")
	if err != nil {
		return []NorthFlowMinute{}, nil
	}
	var res struct {
		Time []string  `json:"time"`
		Hgt  []float64 `json:"hgt"`
		Sgt  []float64 `json:"sgt"`
	}
	if err := json.Unmarshal(resp.Body(), &res); err != nil {
		return []NorthFlowMinute{}, nil
	}
	now := time.Now()
	yrs := tradecalendar.CalendarYears()
	var days []tradecalendar.Day
	for _, y := range yrs {
		if y == now.Year() {
			d, _ := tradecalendar.TradeCalendar(y)
			days = d
			break
		}
	}
	var latest string
	if len(days) > 0 {
		for i := len(days) - 1; i >= 0; i-- {
			dt, _ := time.Parse("2006-01-02", days[i].TradeDate)
			if !dt.After(now) && days[i].TradeStatus == 1 {
				latest = days[i].TradeDate
				break
			}
		}
	}
	out := make([]NorthFlowMinute, 0, len(res.Time))
	for i := range res.Time {
		tr := latest + " " + res.Time[i]
		out = append(out, NorthFlowMinute{TradeTime: tr, NetHgt: float64(int64(res.Hgt[i] * 100000000)), NetSgt: float64(int64(res.Sgt[i] * 100000000)), NetTgt: float64(int64((res.Hgt[i] + res.Sgt[i]) * 100000000))})
	}
	return out, nil
}

func northFlowMinEast(wait time.Duration) ([]NorthFlowMinute, error) {
	client := httpc.NewClient()
	url := "https://push2.eastmoney.com/api/qt/kamt.rtmin/get?fields1=f1,f3&fields2=f51,f52,f54,f56&ut=b2884a393a59ad64002292a3e90d46a5"
	if wait > 0 {
		time.Sleep(wait)
	}
	resp, err := client.R().Get(url)
	if err != nil {
		return []NorthFlowMinute{}, nil
	}
	buf := new(strings.Builder)
	if _, err := io.Copy(buf, strings.NewReader(resp.String())); err != nil {
		return []NorthFlowMinute{}, nil
	}
	text := buf.String()
	l := strings.Index(text, "{")
	if l < 0 {
		return []NorthFlowMinute{}, nil
	}
	var res struct {
		Data struct {
			S2nDate string   `json:"s2nDate"`
			S2n     []string `json:"s2n"`
		} `json:"data"`
	}
	if err := json.Unmarshal([]byte(text[l:len(text)-2]), &res); err != nil {
		return []NorthFlowMinute{}, nil
	}
	y := time.Now().Format("2006")
	out := make([]NorthFlowMinute, 0, len(res.Data.S2n))
	for _, row := range res.Data.S2n {
		cols := strings.Split(row, ",")
		if len(cols) < 4 {
			continue
		}
		if cols[1] == "-" {
			continue
		}
		out = append(out, NorthFlowMinute{TradeTime: y + "-" + res.Data.S2nDate + " " + cols[0], NetHgt: float64(int64(parseF(cols[1]) * 10000)), NetSgt: float64(int64(parseF(cols[2]) * 10000)), NetTgt: float64(int64(parseF(cols[3]) * 10000))})
	}
	return out, nil
}

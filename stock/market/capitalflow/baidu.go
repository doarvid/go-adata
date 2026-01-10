package capitalflow

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func GetStockCapitalFlowMinBaidu(stockCode string, wait time.Duration) ([]FlowMin, error) {
	if stockCode == "" {
		return []FlowMin{}, nil
	}
	client := getHTTPClient()
	url := fmt.Sprintf("https://finance.pae.baidu.com/vapi/v1/fundflow?finance_type=stock&fund_flow_type=&type=stock&market=ab&code=%s&belongs=stocklevelone&finClientType=pc", stockCode)
	if wait > 0 {
		time.Sleep(wait)
	}
	resp, err := client.R().Get(url)
	if err != nil {
		return nil, err
	}

	var res struct {
		Result struct {
			Content struct {
				FundFlowMinute struct {
					Data string `json:"data"`
				} `json:"fundFlowMinute"`
			} `json:"content"`
		} `json:"Result"`
	}
	if err := json.Unmarshal(resp.Body(), &res); err != nil {
		return nil, err
	}
	data := res.Result.Content.FundFlowMinute.Data
	if data == "" {
		return []FlowMin{}, nil
	}
	rows := strings.Split(data, ";")
	out := make([]FlowMin, 0, len(rows))
	for _, r := range rows {
		cols := strings.Split(r, ",")
		if len(cols) < 8 {
			continue
		}
		fm := FlowMin{StockCode: stockCode}
		fm.TradeTime = strings.TrimSpace(cols[0])
		fm.MainNetInflow = parseMoney(cols[2])
		fm.SmNetInflow = parseMoney(cols[7])
		fm.MidNetInflow = parseMoney(cols[6])
		fm.LgNetInflow = parseMoney(cols[5])
		fm.MaxNetInflow = parseMoney(cols[4])
		out = append(out, fm)
	}
	return out, nil
}

func GetStockCapitalFlowBaidu(stockCode string, startDate string, endDate string, wait time.Duration) ([]FlowDaily, error) {
	if stockCode == "" {
		return []FlowDaily{}, nil
	}
	client := getHTTPClient()
	if endDate == "" {
		endDate = time.Now().Format("20060102")
	} else {
		endDate = strings.ReplaceAll(endDate, "-", "")
	}
	if startDate != "" {
		startDate = strings.ReplaceAll(startDate, "-", "")
	}
	out := make([]FlowDaily, 0, 256)
	isEnd := false
	for i := 0; i < 500; i++ {
		url := fmt.Sprintf("https://finance.pae.baidu.com/vapi/v1/fundsortlist?code=%s&market=ab&finance_type=stock&tab=day&from=history&date=%s&pn=0&rn=20&finClientType=pc", stockCode, endDate)
		if wait > 0 {
			time.Sleep(wait)
		}
		resp, err := client.R().Get(url)
		if err != nil {
			return nil, err
		}
		var res struct {
			Result struct {
				Content []map[string]string `json:"content"`
			} `json:"Result"`
		}
		if err := json.Unmarshal(resp.Body(), &res); err != nil {
			return nil, err
		}
		list := res.Result.Content
		if len(list) == 0 {
			break
		}
		for _, row := range list {
			date := strings.ReplaceAll(row["date"], "/", "-")
			yyyymmdd := strings.ReplaceAll(row["date"], "/", "")
			if startDate != "" && yyyymmdd < startDate {
				isEnd = true
				break
			}
			fd := FlowDaily{StockCode: stockCode}
			fd.TradeDate = date
			fd.MainNetInflow = parseMoney(row["extMainIn"])
			fd.SmNetInflow = parseMoney(row["littleNetIn"])
			fd.MidNetInflow = parseMoney(row["mediumNetIn"])
			fd.LgNetInflow = parseMoney(row["largeNetIn"])
			fd.MaxNetInflow = parseMoney(row["superNetIn"])
			out = append(out, fd)
		}
		if isEnd {
			break
		}
		if len(out) > 0 {
			endDate = strings.ReplaceAll(out[len(out)-1].TradeDate, "-", "")
		}
	}
	return out, nil
}

func parseMoney(s string) float64 {
	s = strings.TrimSpace(strings.ReplaceAll(s, "元", ""))
	if s == "" {
		return 0
	}
	re := regexp.MustCompile(`([-+]?\d*\.\d+|\d+)([亿万]?)`)
	m := re.FindStringSubmatch(s)
	if len(m) < 3 {
		v, _ := strconvParseFloat(s)
		return v
	}
	num, _ := strconvParseFloat(m[1])
	unit := m[2]
	if unit == "亿" {
		return num * 100000000
	}
	if unit == "万" {
		return num * 10000
	}
	return num
}

func strconvParseFloat(s string) (float64, error) {
	return strconv.ParseFloat(strings.TrimSpace(strings.ReplaceAll(s, ",", "")), 64)
}

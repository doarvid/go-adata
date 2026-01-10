package indexmarket

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/doarvid/go-adata/stock/info/stockindex"
)

func mapIndexToTHS(indexCode string) string {
	m, err := stockindex.LoadIndexCodeRelTHS()
	if err != nil {
		return indexCode
	}
	if v, ok := m[indexCode]; ok && v != "" {
		return v
	}
	return indexCode
}

func yearsFrom(startDate string) []string {
	if startDate == "" {
		return []string{time.Now().Format("2006")}
	}
	s := strings.Split(startDate, "-")
	y := s[0]
	sy, _ := strconv.Atoi(y)
	cy := time.Now().Year()
	out := make([]string, 0, cy-int(sy)+1)
	for yy := int(sy); yy <= cy; yy++ {
		out = append(out, fmt.Sprintf("%04d", yy))
	}
	return out
}

func (im *IndexMarket) GetDailyThs(ctx context.Context, indexCode, startDate string, kType int, wait time.Duration) ([]IndexDailyBar, error) {
	if indexCode == "" {
		return []IndexDailyBar{}, nil
	}
	concept := mapIndexToTHS(indexCode)
	yrs := yearsFrom(startDate)
	out := make([]IndexDailyBar, 0, 1024)
	for _, year := range yrs {
		url := fmt.Sprintf("http://d.10jqka.com.cn/v4/line/zs_%s/%d1/%s.js", concept, kType-1, year)
		if wait > 0 {
			time.Sleep(wait)
		}
		resp, err := im.client.R().SetContext(ctx).Get(url)
		if err != nil {
			continue
		}
		text := resp.String()
		l := strings.Index(text, "{")
		if l < 0 {
			continue
		}
		payload := text[l : len(text)-1]
		var j struct {
			Data string `json:"data"`
		}
		if err := json.Unmarshal([]byte(payload), &j); err != nil {
			continue
		}
		rows := strings.Split(j.Data, ";")
		for _, r := range rows {
			cols := strings.Split(r, ",")
			if len(cols) < 7 {
				continue
			}
			bar := IndexDailyBar{IndexCode: indexCode}
			bar.TradeDate = parseDate8(cols[0])
			bar.TradeTime = bar.TradeDate
			bar.Open = parseF(cols[1])
			bar.High = parseF(cols[2])
			bar.Low = parseF(cols[3])
			bar.Close = parseF(cols[4])
			bar.Volume = parseF(cols[5])
			bar.Amount = parseF(cols[6])
			out = append(out, bar)
		}
	}
	if len(out) == 0 {
		return []IndexDailyBar{}, nil
	}
	sort.Slice(out, func(i, j int) bool { return out[i].TradeDate < out[j].TradeDate })
	for i := range out {
		if i > 0 {
			out[i].Change = out[i].Close - out[i-1].Close
			if out[i-1].Close != 0 {
				out[i].ChangePct = out[i].Change / out[i-1].Close * 100
			}
		}
	}
	return out, nil
}

func (im *IndexMarket) GetMinuteThs(ctx context.Context, indexCode string, wait time.Duration) ([]IndexMinuteBar, error) {
	if indexCode == "" {
		return []IndexMinuteBar{}, nil
	}
	concept := mapIndexToTHS(indexCode)
	url := fmt.Sprintf("http://d.10jqka.com.cn/v4/time/zs_%s/last.js", concept)
	if wait > 0 {
		time.Sleep(wait)
	}
	resp, err := im.client.R().SetContext(ctx).Get(url)
	if err != nil {
		return nil, err
	}
	text := resp.String()
	l := strings.Index(text, "{")
	if l < 0 {
		return []IndexMinuteBar{}, nil
	}
	payload := text[l : len(text)-1]
	var obj map[string]map[string]any
	if err := json.Unmarshal([]byte(payload), &obj); err != nil {
		return nil, err
	}
	key := fmt.Sprintf("zs_%s", concept)
	data := obj[key]
	if data == nil {
		return []IndexMinuteBar{}, nil
	}
	pre := parseF(toString(data["pre"]))
	tradeDate := toString(data["date"]) // yyyymmdd
	lines := strings.Split(toString(data["data"]), ";")
	out := make([]IndexMinuteBar, 0, len(lines))
	for _, r := range lines {
		if r == "" {
			continue
		}
		cols := strings.Split(r, ",")
		if len(cols) < 5 {
			continue
		}
		t := tradeDate + cols[0]
		tm, _ := time.Parse("200601021504", t)
		bar := IndexMinuteBar{IndexCode: indexCode}
		bar.TradeTime = tm.Format("2006-01-02 15:04:05")
		bar.TradeDate = tm.Format("2006-01-02")
		bar.Price = parseF(cols[1])
		bar.Amount = parseF(cols[2])
		bar.AvgPrice = parseF(cols[3])
		bar.Volume = int64(parseF(cols[4]))
		bar.Change = bar.Price - pre
		if pre != 0 {
			bar.ChangePct = bar.Change / pre * 100
		}
		out = append(out, bar)
	}
	return out, nil
}

func (im *IndexMarket) GetCurrentThs(ctx context.Context, indexCode string, wait time.Duration) (IndexCurrent, error) {
	if indexCode == "" {
		return IndexCurrent{}, nil
	}
	concept := mapIndexToTHS(indexCode)
	url := fmt.Sprintf("http://d.10jqka.com.cn/v4/line/zs_%s/01/today.js", concept)
	if wait > 0 {
		time.Sleep(wait)
	}
	resp, err := im.client.R().SetContext(ctx).Get(url)
	if err != nil {
		return IndexCurrent{}, err
	}
	text := resp.String()
	l := strings.Index(text, "{")
	if l < 0 {
		return IndexCurrent{}, nil
	}
	payload := text[l : len(text)-1]
	var obj map[string]map[string]any
	if err := json.Unmarshal([]byte(payload), &obj); err != nil {
		return IndexCurrent{}, err
	}
	key := fmt.Sprintf("zs_%s", concept)
	data := obj[key]
	if data == nil {
		return IndexCurrent{}, nil
	}
	cur := IndexCurrent{IndexCode: indexCode}
	cur.TradeDate = toString(data["1"]) // yyyymmdd
	dt := toString(data["dt"])          // HHMM
	tm, _ := time.Parse("200601021504", cur.TradeDate+dt)
	cur.TradeDate = tm.Format("2006-01-02")
	cur.TradeTime = tm.Format("2006-01-02 15:04:05")
	cur.Open = parseF(toString(data["7"]))
	cur.High = parseF(toString(data["8"]))
	cur.Low = parseF(toString(data["9"]))
	cur.Price = parseF(toString(data["11"]))
	cur.Volume = parseF(toString(data["13"]))
	cur.Amount = parseF(toString(data["19"]))
	return cur, nil
}

func parseDate8(s string) string {
	if len(s) != 8 {
		return s
	}
	t, err := time.Parse("20060102", s)
	if err != nil {
		return s
	}
	return t.Format("2006-01-02")
}

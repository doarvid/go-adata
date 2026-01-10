package conceptmarket

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

func GetConceptDailyThs(indexCode string, kType int, adjustType int) ([]ConceptDailyBar, error) {
	if indexCode == "" {
		return []ConceptDailyBar{}, nil
	}
	if !strings.HasPrefix(indexCode, "8") {
		return []ConceptDailyBar{}, nil
	}
	client := getHTTPClient()
	url := fmt.Sprintf("http://d.10jqka.com.cn/v6/line/48_%s/%02d%d/last1800.js", indexCode, kType-1, adjustType)
	// wait handled by caller
	resp, err := client.R().Get(url)
	if err != nil {
		return nil, err
	}
	text := resp.String()
	l := strings.Index(text, "{")
	if l < 0 {
		return []ConceptDailyBar{}, nil
	}
	payload := text[l:]
	var j struct {
		Data string `json:"data"`
	}
	if err := json.Unmarshal([]byte(payload[:len(payload)-1]), &j); err != nil {
		return nil, err
	}
	rows := strings.Split(j.Data, ";")
	out := make([]ConceptDailyBar, 0, len(rows))
	for _, r := range rows {
		if r == "" {
			continue
		}
		cols := strings.Split(r, ",")
		if len(cols) < 7 {
			continue
		}
		dt := cols[0]
		bar := ConceptDailyBar{IndexCode: indexCode}
		bar.TradeDate = parseDate8(dt)
		bar.TradeTime = bar.TradeDate
		bar.Open = parseF(cols[1])
		bar.High = parseF(cols[2])
		bar.Low = parseF(cols[3])
		bar.Close = parseF(cols[4])
		bar.Volume = parseF(cols[5])
		bar.Amount = parseF(cols[6])
		out = append(out, bar)
	}
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

func GetConceptMinuteThs(indexCode string) ([]ConceptMinuteBar, error) {
	if indexCode == "" {
		return []ConceptMinuteBar{}, nil
	}
	if !strings.HasPrefix(indexCode, "8") {
		return []ConceptMinuteBar{}, nil
	}
	client := getHTTPClient()
	url := fmt.Sprintf("http://d.10jqka.com.cn/v6/time/48_%s/last.js", indexCode)
	// wait handled by caller
	resp, err := client.R().Get(url)
	if err != nil {
		return nil, err
	}
	text := resp.String()
	l := strings.Index(text, "{")
	if l < 0 {
		return []ConceptMinuteBar{}, nil
	}
	payload := text[l : len(text)-1]
	var obj map[string]map[string]any
	if err := json.Unmarshal([]byte(payload), &obj); err != nil {
		return nil, err
	}
	key := fmt.Sprintf("48_%s", indexCode)
	data := obj[key]
	if data == nil {
		return []ConceptMinuteBar{}, nil
	}
	pre := parseF(toString(data["pre"]))
	tradeDate := toString(data["date"]) // yyyymmdd
	dl := toString(data["data"])        // semicolon lines
	lines := strings.Split(dl, ";")
	out := make([]ConceptMinuteBar, 0, len(lines))
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
		bar := ConceptMinuteBar{IndexCode: indexCode}
		bar.TradeTime = tm.Format("2006-01-02 15:04:05")
		bar.TradeDate = tm.Format("2006-01-02")
		bar.Price = parseF(cols[1])
		bar.Amount = parseF(cols[2])
		bar.AvgPrice = 0
		bar.Volume = int64(parseF(cols[4]))
		bar.Change = bar.Price - pre
		if pre != 0 {
			bar.ChangePct = bar.Change / pre * 100
		}
		out = append(out, bar)
	}
	return out, nil
}

func GetConceptCurrentThs(indexCode string, kType int) (ConceptCurrent, error) {
	if indexCode == "" {
		return ConceptCurrent{}, nil
	}
	if !strings.HasPrefix(indexCode, "8") {
		return ConceptCurrent{}, nil
	}
	client := getHTTPClient()
	url := fmt.Sprintf("http://d.10jqka.com.cn/v6/line/48_%s/%02d1/today.js", indexCode, kType-1)
	// wait handled by caller
	resp, err := client.R().Get(url)
	if err != nil {
		return ConceptCurrent{}, err
	}
	text := resp.String()
	l := strings.Index(text, "{")
	if l < 0 {
		return ConceptCurrent{}, nil
	}
	payload := text[l : len(text)-1]
	var obj map[string]map[string]any
	if err := json.Unmarshal([]byte(payload), &obj); err != nil {
		return ConceptCurrent{}, err
	}
	key := fmt.Sprintf("48_%s", indexCode)
	data := obj[key]
	if data == nil {
		return ConceptCurrent{}, nil
	}
	cur := ConceptCurrent{IndexCode: indexCode}
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

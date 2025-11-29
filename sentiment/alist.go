package sentiment

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	httpc "github.com/doarvid/go-adata/common/http"
)

type AListDaily struct {
	TradeDate     string  `json:"trade_date"`
	ShortName     string  `json:"short_name"`
	StockCode     string  `json:"stock_code"`
	Close         float64 `json:"close"`
	ChangeCpt     float64 `json:"change_cpt"`
	TurnoverRatio float64 `json:"turnover_ratio"`
	ANetAmount    float64 `json:"a_net_amount"`
	ABuyAmount    float64 `json:"a_buy_amount"`
	ASellAmount   float64 `json:"a_sell_amount"`
	AAmount       float64 `json:"a_amount"`
	Amount        float64 `json:"amount"`
	NetAmountRate float64 `json:"net_amount_rate"`
	AAmountRate   float64 `json:"a_amount_rate"`
	Reason        string  `json:"reason"`
}

type AListInfo struct {
	TradeDate       string  `json:"trade_date"`
	StockCode       string  `json:"stock_code"`
	OperateCode     string  `json:"operate_code"`
	OperateName     string  `json:"operate_name"`
	ABuyAmount      float64 `json:"a_buy_amount"`
	ASellAmount     float64 `json:"a_sell_amount"`
	ANetAmount      float64 `json:"a_net_amount"`
	ABuyAmountRate  float64 `json:"a_buy_amount_rate"`
	ASellAmountRate float64 `json:"a_sell_amount_rate"`
	Reason          string  `json:"reason"`
}

func ListAListDaily(reportDate string, wait time.Duration) ([]AListDaily, error) {
	if reportDate == "" {
		reportDate = time.Now().Format("2006-01-02")
	}
	client := httpc.NewClient()
	url := "https://datacenter-web.eastmoney.com/api/data/v1/get?sortColumns=SECURITY_CODE,TRADE_DATE&sortTypes=1,-1&pageSize=500&pageNumber=1&reportName=RPT_DAILYBILLBOARD_DETAILSNEW&columns=SECURITY_CODE,SECUCODE,SECURITY_NAME_ABBR,TRADE_DATE,EXPLAIN,CLOSE_PRICE,CHANGE_RATE,BILLBOARD_NET_AMT,BILLBOARD_BUY_AMT,BILLBOARD_SELL_AMT,BILLBOARD_DEAL_AMT,ACCUM_AMOUNT,DEAL_NET_RATIO,DEAL_AMOUNT_RATIO,TURNOVERRATE,FREE_MARKET_CAP,EXPLANATION,D1_CLOSE_ADJCHRATE,D2_CLOSE_ADJCHRATE,D5_CLOSE_ADJCHRATE,D10_CLOSE_ADJCHRATE,SECURITY_TYPE_CODE&source=WEB&client=WEB&filter=(TRADE_DATE='" + reportDate + "')(TRADE_DATE='" + reportDate + "')"
	if wait > 0 {
		time.Sleep(wait)
	}
	resp, err := client.R().Post(url)
	if err != nil {
		return nil, err
	}
	buf := new(strings.Builder)
	if _, err := ioCopy(buf, strings.NewReader(resp.String())); err != nil {
		return nil, err
	}
	text := buf.String()
	l := strings.Index(text, "{")
	if l < 0 {
		return []AListDaily{}, nil
	}
	var res struct {
		Result struct {
			Data []map[string]any `json:"data"`
		} `json:"result"`
	}
	if err := json.Unmarshal([]byte(text[l:len(text)-2]), &res); err != nil {
		return nil, err
	}
	data := res.Result.Data
	out := make([]AListDaily, 0, len(data))
	for _, it := range data {
		out = append(out, AListDaily{
			StockCode:     toString(it["SECURITY_CODE"]),
			ShortName:     strings.ReplaceAll(toString(it["SECURITY_NAME_ABBR"]), " ", ""),
			TradeDate:     toDate(toString(it["TRADE_DATE"])),
			Close:         parseF(toString(it["CLOSE_PRICE"])),
			ChangeCpt:     parseF(toString(it["CHANGE_RATE"])),
			TurnoverRatio: parseF(toString(it["TURNOVERRATE"])),
			ANetAmount:    parseF(toString(it["BILLBOARD_NET_AMT"])),
			ABuyAmount:    parseF(toString(it["BILLBOARD_BUY_AMT"])),
			ASellAmount:   parseF(toString(it["BILLBOARD_SELL_AMT"])),
			AAmount:       parseF(toString(it["BILLBOARD_DEAL_AMT"])),
			Amount:        parseF(toString(it["ACCUM_AMOUNT"])),
			NetAmountRate: parseF(toString(it["DEAL_NET_RATIO"])),
			AAmountRate:   parseF(toString(it["DEAL_AMOUNT_RATIO"])),
			Reason:        toString(it["EXPLANATION"]),
		})
	}
	return out, nil
}

func GetAListInfo(stockCode string, reportDate string, wait time.Duration) ([]AListInfo, error) {
	if stockCode == "" {
		return []AListInfo{}, nil
	}
	if reportDate == "" {
		reportDate = time.Now().Format("2006-01-02")
	}
	urls := []string{
		"https://datacenter-web.eastmoney.com/api/data/v1/get?reportName=RPT_BILLBOARD_DAILYDETAILSBUY&columns=ALL&filter=(TRADE_DATE=%27" + reportDate + "%27)(SECURITY_CODE=%22" + stockCode + "%22)&pageNumber=1&pageSize=50&sortTypes=-1&sortColumns=BUY&source=WEB&client=WEB",
		"https://datacenter-web.eastmoney.com/api/data/v1/get?reportName=RPT_BILLBOARD_DAILYDETAILSSELL&columns=ALL&filter=(TRADE_DATE=%27" + reportDate + "%27)(SECURITY_CODE=%22" + stockCode + "%22)&pageNumber=1&pageSize=50&sortTypes=-1&sortColumns=BUY&source=WEB&client=WEB",
	}
	client := httpc.NewClient()
	out := make([]AListInfo, 0, 0)
	for _, url := range urls {
		if wait > 0 {
			time.Sleep(wait)
		}
		resp, err := client.R().Post(url)
		if err != nil {
			continue
		}
		var res struct {
			Result struct {
				Data []map[string]any `json:"data"`
			} `json:"result"`
		}
		if err := json.Unmarshal(resp.Body(), &res); err != nil {
			continue
		}
		for _, it := range res.Result.Data {
			out = append(out, AListInfo{
				StockCode:       toString(it["SECURITY_CODE"]),
				TradeDate:       toDate(toString(it["TRADE_DATE"])),
				OperateCode:     toString(it["OPERATEDEPT_CODE"]),
				OperateName:     toString(it["OPERATEDEPT_NAME"]),
				ABuyAmount:      parseF(toString(it["BUY"])),
				ASellAmount:     parseF(toString(it["SELL"])),
				ANetAmount:      parseF(toString(it["NET"])),
				ABuyAmountRate:  parseF(toString(it["TOTAL_BUYRIO"])),
				ASellAmountRate: parseF(toString(it["TOTAL_SELLRIO"])),
				Reason:          toString(it["EXPLANATION"]),
			})
		}
	}
	return out, nil
}

func parseF(s string) float64 {
	s = strings.TrimSpace(strings.ReplaceAll(s, ",", ""))
	if s == "" || s == "--" {
		return 0
	}
	v, _ := strconv.ParseFloat(strings.TrimSuffix(s, "%"), 64)
	return v
}
func toString(v any) string { return strings.TrimSpace(fmt.Sprintf("%v", v)) }
func toDate(s string) string {
	t, err := time.Parse("2006-01-02 15:04:05", s)
	if err == nil {
		return t.Format("2006-01-02")
	}
	tt, err2 := time.Parse("2006-01-02", s)
	if err2 == nil {
		return tt.Format("2006-01-02")
	}
	return s
}
func ioCopy(w *strings.Builder, r io.Reader) (int64, error) { return io.Copy(w, r) }

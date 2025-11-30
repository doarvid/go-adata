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
	TradeDate     string  `json:"trade_date"`      // 交易日期 2024-07-01
	ShortName     string  `json:"short_name"`      // 股票简称 全新好
	StockCode     string  `json:"stock_code"`      // 股票代码 000007
	Close         float64 `json:"close"`           // 收盘价(元) 5.16000
	ChangeCpt     float64 `json:"change_cpt"`      // 涨跌幅(%) -9.94760
	TurnoverRatio float64 `json:"turnover_ratio"`  // 换手率(%) 8.37860
	ANetAmount    float64 `json:"a_net_amount"`    // 龙虎榜净买入额(元) 4939641.98000
	ABuyAmount    float64 `json:"a_buy_amount"`    // 龙虎榜买入额(元) 23347567.29000
	ASellAmount   float64 `json:"a_sell_amount"`   // 龙虎榜卖出额(元) 18407925.31000
	AAmount       float64 `json:"a_amount"`        // 龙虎榜成交额(元) 41755492.60000
	Amount        float64 `json:"amount"`          // 总成交额(元) 137998593
	NetAmountRate float64 `json:"net_amount_rate"` // 龙虎榜净买额占总成交额比例(%) 3.57949
	AAmountRate   float64 `json:"a_amount_rate"`   // 龙虎榜成交额占总成交额比例(%) 30.25791
	Reason        string  `json:"reason"`          // 上榜原因 日跌幅偏离值达到7%的前5只证券
}

type AListInfo struct {
	TradeDate       string  `json:"trade_date"`         // 交易日期 2024-07-01
	StockCode       string  `json:"stock_code"`         // 股票代码 000007
	OperateCode     string  `json:"operate_code"`       // 营业部代码 10678762
	OperateName     string  `json:"operate_name"`       // 营业部名称 东方财富证券股份有限公司拉萨金融城南环路证券营业部
	ABuyAmount      float64 `json:"a_buy_amount"`       // 龙虎榜买入额(元) 23347567.29000
	ASellAmount     float64 `json:"a_sell_amount"`      // 龙虎榜卖出额(元) 18407925.31000
	ANetAmount      float64 `json:"a_net_amount"`       // 龙虎榜净买入额(元) 4939641.98000
	ABuyAmountRate  float64 `json:"a_buy_amount_rate"`  // 龙虎榜买入额占总成交额比例(%) 3.57949
	ASellAmountRate float64 `json:"a_sell_amount_rate"` // 龙虎榜卖出额占总成交额比例(%) 30.25791
	Reason          string  `json:"reason"`             // 上榜原因 有价格涨跌幅限制的日收盘价格涨幅偏离值达到7%的前五只证券
}

// 每日龙虎榜，默认为当天
// http://guba.eastmoney.com/rank/
func ListAListDaily(reportDate string, wait time.Duration) ([]AListDaily, error) {
	if reportDate == "" {
		reportDate = time.Now().Format("2006-01-02")
	}
	client := httpc.NewClient()
	url := "https://datacenter-web.eastmoney.com/api/data/v1/get?sortColumns=SECURITY_CODE,TRADE_DATE&sortTypes=1,-1&pageSize=500&pageNumber=1&reportName=RPT_DAILYBILLBOARD_DETAILSNEW&columns=SECURITY_CODE,SECUCODE,SECURITY_NAME_ABBR,TRADE_DATE,EXPLAIN,CLOSE_PRICE,CHANGE_RATE,BILLBOARD_NET_AMT,BILLBOARD_BUY_AMT,BILLBOARD_SELL_AMT,BILLBOARD_DEAL_AMT,ACCUM_AMOUNT,DEAL_NET_RATIO,DEAL_AMOUNT_RATIO,TURNOVERRATE,FREE_MARKET_CAP,EXPLANATION,D1_CLOSE_ADJCHRATE,D2_CLOSE_ADJCHRATE,D5_CLOSE_ADJCHRATE,D10_CLOSE_ADJCHRATE,SECURITY_TYPE_CODE&source=WEB&client=WEB&filter=(TRADE_DATE=%27" + reportDate + "%27)(TRADE_DATE=%27" + reportDate + "%27)"
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
	if err := json.Unmarshal([]byte(text), &res); err != nil {
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

// 获取单个龙虎榜的数据，买5和卖5
// https://datacenter-web.eastmoney.com/api/data/v1/get?callback=jQuery1123015874658470862357_1721014447038&reportName=RPT_BILLBOARD_DAILYDETAILSBUY&columns=ALL&filter=(TRADE_DATE='2024-07-12')(SECURITY_CODE="600297")&pageNumber=1&pageSize=50&sortTypes=-1&sortColumns=BUY&source=WEB&client=WEB&_=1721014447040
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

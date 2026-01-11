package stockmarket

import (
	"context"
	"time"

	browser "github.com/EDDYCJY/fake-useragent"
	"github.com/go-resty/resty/v2"
)

type KType int

const (
	KTypeDay      KType = 1
	KTypeMinute   KType = 2
	KTypeMonth    KType = 3
	KTypeQuarter  KType = 4
	KTypeMinite5  KType = 5
	KTypeMinite15 KType = 15
	KTypeMinite30 KType = 30
	KTypeMinite60 KType = 60
)

type AdjustType int

const (
	AdjustTypeNone AdjustType = 0
	AdjustTypePre  AdjustType = 1
	AdjustTypePost AdjustType = 2
)

type Market struct {
	MinWait time.Duration
	Retries int

	proxy  string
	debug  bool
	client *resty.Client
}

type DailyBar struct {
	// 股票代码，示例：600001
	StockCode string `json:"stock_code"`
	// 交易时间，示例：1990-01-01 00:00:00；分时图使用具体的时间
	TradeTime string `json:"trade_time"`
	// 交易日期，示例：1990-01-01
	TradeDate string `json:"trade_date"`
	// 开盘价(元)，示例：9.98
	Open float64 `json:"open"`
	// 收盘价(元)，示例：9.98
	Close float64 `json:"close"`
	// 最高价(元)，示例：9.98
	High float64 `json:"high"`
	// 最低价(元)，示例：9.98
	Low float64 `json:"low"`
	// 成交量(股)，示例：64745722
	Volume float64 `json:"volume"`
	// 成交额(元)，示例：934285179.00
	Amount float64 `json:"amount"`
	// 涨跌额(元)，示例：-0.02
	Change float64 `json:"change"`
	// 涨跌幅(%)，示例：-0.16
	ChangePct float64 `json:"change_pct"`
	// 换手率(%)，示例：0.38
	TurnoverRatio string `json:"turnover_ratio"`
	// 昨收(元)，示例：10.00
	PreClose float64 `json:"pre_close"`
}

// 分时 K 线结构体
type MinuteBar struct {
	Time      int64   `json:"time"`       // 时间戳（秒）  例如：1710000000
	Price     float64 `json:"price"`      // 价格(元)     例如：9.98
	ChangePct float64 `json:"change_pct"` // 涨跌幅(%)    例如：-0.16
	Change    float64 `json:"change"`     // 涨跌额(元)   例如：-0.02
	AvgPrice  float64 `json:"avg_price"`  // 平均价(元)   例如：9.98
	Volume    int64   `json:"volume"`     // 成交量(股)   例如：64745722
	Amount    float64 `json:"amount"`     // 成交额(元)   例如：934285179.00
	Open      float64 `json:"open"`       // 开盘价(元)   例如：10.00
	Close     float64 `json:"close"`      // 收盘价(元)   例如：9.98
	High      float64 `json:"high"`       // 最高价(元)   例如：10.05
	Low       float64 `json:"low"`        // 最低价(元)   例如：9.95
	TradeTime string  `json:"trade_time"` // 交易时间     例如：2024-01-01 14:55:00
	TradeDate string  `json:"trade_date"` // 交易日期     例如：2024-01-01
	StockCode string  `json:"stock_code"` // 股票代码     例如：600001
}

// 逐笔成交结构体
type TickBar struct {
	TradeTime string  `json:"trade_time"` // 成交时间	例如：2023-09-13 09:31:45
	Volume    int64   `json:"volume"`     // 成交量(股)	例如：34452500
	Price     float64 `json:"price"`      // 当前价格(元)	例如：12.36
	Type      string  `json:"type"`       // 类型	例如：--
	BsType    string  `json:"bs_type"`    // 买卖类型	B：买入，S：卖出
	StockCode string  `json:"stock_code"` // 代码	例如：600001
}

type Five struct {
	StockCode string  `json:"stock_code"` // 代码，如：600001
	ShortName string  `json:"short_name"` // 简称，如：平安银行
	S5        float64 `json:"s5"`         // 卖5价(元)，如：11.29
	Sv5       int64   `json:"sv5"`        // 卖5量(股)，如：2263
	S4        float64 `json:"s4"`         // 卖4价(元)，如：11.28
	Sv4       int64   `json:"sv4"`        // 卖4量(股)，如：2263
	S3        float64 `json:"s3"`         // 卖3价(元)，如：11.27
	Sv3       int64   `json:"sv3"`        // 卖3量(股)，如：2263
	S2        float64 `json:"s2"`         // 卖2价(元)，如：11.26
	Sv2       int64   `json:"sv2"`        // 卖2量(股)，如：2263
	S1        float64 `json:"s1"`         // 卖1价(元)，如：11.25
	Sv1       int64   `json:"sv1"`        // 卖1量(股)，如：2263
	B1        float64 `json:"b1"`         // 买1价(元)，如：11.24
	Bv1       int64   `json:"bv1"`        // 买1量(股)，如：2263
	B2        float64 `json:"b2"`         // 买2价(元)，如：11.23
	Bv2       int64   `json:"bv2"`        // 买2量(股)，如：2263
	B3        float64 `json:"b3"`         // 买3价(元)，如：11.22
	Bv3       int64   `json:"bv3"`        // 买3量(股)，如：2263
	B4        float64 `json:"b4"`         // 买4价(元)，如：11.21
	Bv4       int64   `json:"bv4"`        // 买4量(股)，如：2263
	B5        float64 `json:"b5"`         // 买5价(元)，如：11.20
	Bv5       int64   `json:"bv5"`        // 买5量(股)，如：2263
}

type MarketOpt func(*Market)

func WithProxy(proxy string) MarketOpt {
	return func(m *Market) {
		m.proxy = proxy
	}
}
func WithDebug(enable bool) MarketOpt {
	return func(m *Market) {
		m.debug = enable
	}
}

func NewMarket(opts ...MarketOpt) *Market {
	m := &Market{MinWait: 50 * time.Millisecond, Retries: 2}
	for _, opt := range opts {
		opt(m)
	}
	c := resty.New()
	c.SetTimeout(15 * time.Second)
	c.SetHeader("User-Agent", browser.Random())
	if m.proxy != "" {
		c.SetProxy(m.proxy)
	}
	if m.debug {
		c.SetDebug(true)
	}
	m.client = c
	return m
}

func (m *Market) GetDaily(stockCode, startDate, endDate string, kType KType, adjustType AdjustType) ([]DailyBar, error) {
	if stockCode == "" {
		return []DailyBar{}, nil
	}
	var east []DailyBar
	var err error
	for i := 0; i <= m.Retries; i++ {
		east, err = m.GetDailyEast(context.Background(), stockCode, startDate, endDate, kType, adjustType)
		if err == nil && len(east) > 0 {
			return NormalizeDaily(east), nil
		}
		time.Sleep(m.MinWait)
	}
	var bd []DailyBar
	var err2 error
	for i := 0; i <= m.Retries; i++ {
		bd, err2 = m.GetDailyBaidu(context.Background(), stockCode, startDate, kType)
		if err2 == nil && len(bd) > 0 {
			return NormalizeDaily(bd), nil
		}
		time.Sleep(m.MinWait)
	}
	if len(bd) > 0 {
		return NormalizeDaily(bd), err
	}
	return NormalizeDaily(east), err
}

func (m *Market) GetMinute(stockCode string) ([]MinuteBar, error) {
	if stockCode == "" {
		return []MinuteBar{}, nil
	}
	var east []MinuteBar
	var err error
	for i := 0; i <= m.Retries; i++ {
		east, err = m.GetMinuteEast(context.Background(), stockCode)
		if err == nil && len(east) > 0 {
			return NormalizeMinute(east), nil
		}
		time.Sleep(m.MinWait)
	}
	var bd []MinuteBar
	var err2 error
	for i := 0; i <= m.Retries; i++ {
		bd, err2 = m.GetMinuteBaidu(context.Background(), stockCode)
		if err2 == nil && len(bd) > 0 {
			return NormalizeMinute(bd), nil
		}
		time.Sleep(m.MinWait)
	}
	if len(bd) > 0 {
		return NormalizeMinute(bd), err
	}
	return NormalizeMinute(east), err
}

func (m *Market) GetBar(stockCode string) ([]TickBar, error) {
	if stockCode == "" {
		return []TickBar{}, nil
	}
	var bd []TickBar
	var err error
	for i := 0; i <= m.Retries; i++ {
		bd, err := m.GetBarBaidu(context.Background(), stockCode)
		if err == nil && len(bd) > 0 {
			return NormalizeTick(bd), nil
		}
		time.Sleep(m.MinWait)
	}
	var qq []TickBar
	var err2 error
	for i := 0; i <= m.Retries; i++ {
		qq, err2 = m.GetBarQQ(context.Background(), stockCode)
		if err2 == nil && len(qq) > 0 {
			return NormalizeTick(qq), nil
		}
		time.Sleep(m.MinWait)
	}
	if err == nil {
		return NormalizeTick(bd), err2
	}
	return NormalizeTick(qq), err
}

func (m *Market) GetFive(stockCode string) (Five, error) {
	if stockCode == "" {
		return Five{}, nil
	}
	qq, err := m.GetFiveQQ(context.Background(), stockCode)
	if err == nil && qq.ShortName != "" {
		return NormalizeFive(qq), nil
	}
	bd, err2 := m.GetFiveBaidu(context.Background(), stockCode)
	if err2 == nil && bd.ShortName != "" {
		return NormalizeFive(bd), nil
	}
	if err == nil {
		return NormalizeFive(qq), err2
	}
	return NormalizeFive(bd), err
}
func (m *Market) ListCurrent(codeList []string) ([]CurrentQuote, error) {
	// 优先新浪，失败或空则腾讯
	sina, err := m.ListCurrentSina(context.Background(), codeList)
	if err == nil && len(sina) > 0 {
		return NormalizeCurrent(sina), nil
	}
	qq, err2 := m.ListCurrentQQ(context.Background(), codeList)
	if err2 == nil && len(qq) > 0 {
		return NormalizeCurrent(qq), nil
	}
	if err == nil {
		return NormalizeCurrent(sina), err2
	}
	return NormalizeCurrent(qq), err
}

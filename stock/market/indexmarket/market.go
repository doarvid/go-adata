package indexmarket

import "time"

type IndexDailyBar struct {
	TradeTime     string  `json:"trade_time"`
	TradeDate     string  `json:"trade_date"`
	Open          float64 `json:"open"`
	Close         float64 `json:"close"`
	High          float64 `json:"high"`
	Low           float64 `json:"low"`
	Volume        float64 `json:"volume"`
	Amount        float64 `json:"amount"`
	Change        float64 `json:"change"`
	ChangePct     float64 `json:"change_pct"`
	IndexCode     string  `json:"index_code"`
	TurnoverRatio string  `json:"turnover_ratio"`
	PreClose      float64 `json:"pre_close"`
}

type IndexMinuteBar struct {
	TradeTime string  `json:"trade_time"`
	TradeDate string  `json:"trade_date"`
	Price     float64 `json:"price"`
	Change    float64 `json:"change"`
	ChangePct float64 `json:"change_pct"`
	Volume    int64   `json:"volume"`
	AvgPrice  float64 `json:"avg_price"`
	Amount    float64 `json:"amount"`
	Open      float64 `json:"open"`
	Close     float64 `json:"close"`
	High      float64 `json:"high"`
	Low       float64 `json:"low"`
	IndexCode string  `json:"index_code"`
}

type IndexCurrent struct {
	TradeTime string  `json:"trade_time"`
	TradeDate string  `json:"trade_date"`
	Open      float64 `json:"open"`
	High      float64 `json:"high"`
	Low       float64 `json:"low"`
	Price     float64 `json:"price"`
	Change    float64 `json:"change"`
	ChangePct float64 `json:"change_pct"`
	Volume    float64 `json:"volume"`
	Amount    float64 `json:"amount"`
	IndexCode string  `json:"index_code"`
}

type Market struct {
	MinWait time.Duration
	Retries int
}

func New() *Market { return &Market{MinWait: 50 * time.Millisecond, Retries: 2} }

func (m *Market) GetDaily(indexCode, startDate string, kType int, wait time.Duration) ([]IndexDailyBar, error) {
	if indexCode == "" {
		return []IndexDailyBar{}, nil
	}
	if wait < m.MinWait {
		wait = m.MinWait
	}
	var east []IndexDailyBar
	var err error
	for i := 0; i <= m.Retries; i++ {
		east, err = GetIndexDailyEast(indexCode, startDate, kType, wait)
		if err == nil && len(east) > 0 {
			return NormalizeIndexDaily(east), nil
		}
		time.Sleep(wait)
	}
	var bd []IndexDailyBar
	var err2 error
	for i := 0; i <= m.Retries; i++ {
		bd, err2 = GetIndexDailyBaidu(indexCode, startDate, kType, wait)
		if err2 == nil && len(bd) > 0 {
			return NormalizeIndexDaily(bd), nil
		}
		time.Sleep(wait)
	}
	var ths []IndexDailyBar
	var err3 error
	for i := 0; i <= m.Retries; i++ {
		ths, err3 = GetIndexDailyThs(indexCode, startDate, kType, wait)
		if err3 == nil && len(ths) > 0 {
			return NormalizeIndexDaily(ths), nil
		}
		time.Sleep(wait)
	}
	if len(bd) > 0 {
		return NormalizeIndexDaily(bd), err
	}
	if len(ths) > 0 {
		return NormalizeIndexDaily(ths), err
	}
	return NormalizeIndexDaily(east), err
}

func (m *Market) GetMinute(indexCode string, wait time.Duration) ([]IndexMinuteBar, error) {
	if indexCode == "" {
		return []IndexMinuteBar{}, nil
	}
	if wait < m.MinWait {
		wait = m.MinWait
	}
	var east []IndexMinuteBar
	var err error
	for i := 0; i <= m.Retries; i++ {
		east, err = GetIndexMinuteEast(indexCode, wait)
		if err == nil && len(east) > 0 {
			return NormalizeIndexMinute(east), nil
		}
		time.Sleep(wait)
	}
	var ths []IndexMinuteBar
	var err2 error
	for i := 0; i <= m.Retries; i++ {
		ths, err2 = GetIndexMinuteThs(indexCode, wait)
		if err2 == nil && len(ths) > 0 {
			return NormalizeIndexMinute(ths), nil
		}
		time.Sleep(wait)
	}
	if len(ths) > 0 {
		return NormalizeIndexMinute(ths), err
	}
	return NormalizeIndexMinute(east), err
}

func (m *Market) GetCurrent(indexCode string, wait time.Duration) (IndexCurrent, error) {
	if indexCode == "" {
		return IndexCurrent{}, nil
	}
	if wait < m.MinWait {
		wait = m.MinWait
	}
	var cur IndexCurrent
	var err error
	for i := 0; i <= m.Retries; i++ {
		cur, err = GetIndexCurrentEast(indexCode, wait)
		if cur.IndexCode != "" {
			return NormalizeIndexCurrent(cur), nil
		}
		time.Sleep(wait)
	}
	var ths IndexCurrent
	var err2 error
	for i := 0; i <= m.Retries; i++ {
		ths, err2 = GetIndexCurrentThs(indexCode, wait)
		if ths.IndexCode != "" {
			return NormalizeIndexCurrent(ths), nil
		}
		time.Sleep(wait)
	}
	var mins []IndexMinuteBar
	var merr error
	for i := 0; i <= m.Retries; i++ {
		mins, merr = GetIndexMinuteEast(indexCode, wait)
		if merr == nil && len(mins) > 0 {
			break
		}
		time.Sleep(wait)
	}
	if len(mins) > 0 {
		last := mins[len(mins)-1]
		cur2 := IndexCurrent{IndexCode: indexCode}
		cur2.Open = last.Open
		cur2.High = last.High
		cur2.Low = last.Low
		cur2.Price = last.Price
		cur2.Change = last.Change
		cur2.ChangePct = last.ChangePct
		cur2.Volume = float64(last.Volume)
		cur2.Amount = last.Amount
		cur2.TradeTime = last.TradeTime
		cur2.TradeDate = last.TradeDate
		return NormalizeIndexCurrent(cur2), nil
	}
	if err == nil {
		return NormalizeIndexCurrent(cur), err2
	}
	return NormalizeIndexCurrent(ths), err
}

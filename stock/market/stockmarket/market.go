package stockmarket

import "time"

type Market struct {
	MinWait time.Duration
	Retries int
}

func NewMarket() *Market { return &Market{MinWait: 50 * time.Millisecond, Retries: 2} }

func (m *Market) GetDaily(stockCode, startDate, endDate string, kType, adjustType int, wait time.Duration) ([]DailyBar, error) {
	if stockCode == "" {
		return []DailyBar{}, nil
	}
	if wait < m.MinWait {
		wait = m.MinWait
	}
	var east []DailyBar
	var err error
	for i := 0; i <= m.Retries; i++ {
		east, err = GetMarketDailyEast(stockCode, startDate, endDate, kType, adjustType, wait)
		if err == nil && len(east) > 0 {
			return NormalizeDaily(east), nil
		}
		time.Sleep(wait)
	}
	var bd []DailyBar
	var err2 error
	for i := 0; i <= m.Retries; i++ {
		bd, err2 = GetMarketDailyBaidu(stockCode, startDate, kType, wait)
		if err2 == nil && len(bd) > 0 {
			return NormalizeDaily(bd), nil
		}
		time.Sleep(wait)
	}
	if len(bd) > 0 {
		return NormalizeDaily(bd), err
	}
	return NormalizeDaily(east), err
}

func (m *Market) GetMinute(stockCode string, wait time.Duration) ([]MinuteBar, error) {
	if stockCode == "" {
		return []MinuteBar{}, nil
	}
	if wait < m.MinWait {
		wait = m.MinWait
	}
	var east []MinuteBar
	var err error
	for i := 0; i <= m.Retries; i++ {
		east, err = GetMarketMinuteEast(stockCode, wait)
		if err == nil && len(east) > 0 {
			return NormalizeMinute(east), nil
		}
		time.Sleep(wait)
	}
	var bd []MinuteBar
	var err2 error
	for i := 0; i <= m.Retries; i++ {
		bd, err2 = GetMarketMinuteBaidu(stockCode, wait)
		if err2 == nil && len(bd) > 0 {
			return NormalizeMinute(bd), nil
		}
		time.Sleep(wait)
	}
	if len(bd) > 0 {
		return NormalizeMinute(bd), err
	}
	return NormalizeMinute(east), err
}

func (m *Market) GetBar(stockCode string, wait time.Duration) ([]TickBar, error) {
	if stockCode == "" {
		return []TickBar{}, nil
	}
	if wait < m.MinWait {
		wait = m.MinWait
	}
	var bd []TickBar
	var err error
	for i := 0; i <= m.Retries; i++ {
		bd, err = GetMarketBarBaidu(stockCode, wait)
		if err == nil && len(bd) > 0 {
			return NormalizeTick(bd), nil
		}
		time.Sleep(wait)
	}
	var qq []TickBar
	var err2 error
	for i := 0; i <= m.Retries; i++ {
		qq, err2 = GetMarketBarQQ(stockCode, wait)
		if err2 == nil && len(qq) > 0 {
			return NormalizeTick(qq), nil
		}
		time.Sleep(wait)
	}
	if err == nil {
		return NormalizeTick(bd), err2
	}
	return NormalizeTick(qq), err
}

func (m *Market) GetFive(stockCode string, wait time.Duration) (Five, error) {
	if stockCode == "" {
		return Five{}, nil
	}
	if wait < m.MinWait {
		wait = m.MinWait
	}
	qq, err := GetMarketFiveQQ(stockCode, wait)
	if err == nil && qq.ShortName != "" {
		return NormalizeFive(qq), nil
	}
	bd, err2 := GetMarketFiveBaidu(stockCode, wait)
	if err2 == nil && bd.ShortName != "" {
		return NormalizeFive(bd), nil
	}
	if err == nil {
		return NormalizeFive(qq), err2
	}
	return NormalizeFive(bd), err
}
func (m *Market) ListCurrent(codeList []string, wait time.Duration) ([]CurrentQuote, error) {
	// 优先新浪，失败或空则腾讯
	if wait < m.MinWait {
		wait = m.MinWait
	}
	sina, err := ListMarketCurrentSina(codeList, wait)
	if err == nil && len(sina) > 0 {
		return NormalizeCurrent(sina), nil
	}
	qq, err2 := ListMarketCurrentQQ(codeList, wait)
	if err2 == nil && len(qq) > 0 {
		return NormalizeCurrent(qq), nil
	}
	if err == nil {
		return NormalizeCurrent(sina), err2
	}
	return NormalizeCurrent(qq), err
}

// Python 风格方法名包装
func (m *Market) GetMarket(stockCode, startDate, endDate string, kType, adjustType int, wait time.Duration) ([]DailyBar, error) {
	return m.GetDaily(stockCode, startDate, endDate, kType, adjustType, wait)
}
func (m *Market) GetMarketMin(stockCode string, wait time.Duration) ([]MinuteBar, error) {
	return m.GetMinute(stockCode, wait)
}
func (m *Market) GetMarketFive(stockCode string, wait time.Duration) (Five, error) {
	return m.GetFive(stockCode, wait)
}
func (m *Market) GetMarketBar(stockCode string, wait time.Duration) ([]TickBar, error) {
	return m.GetBar(stockCode, wait)
}

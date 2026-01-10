package conceptmarket

import "time"

type Market struct { MinWait time.Duration; Retries int }

func New() *Market { return &Market{MinWait: 50 * time.Millisecond, Retries: 2} }

func (m *Market) GetDaily(indexCode string, kType int) ([]ConceptDailyBar, error) {
    if indexCode == "" { return []ConceptDailyBar{}, nil }
    var ths []ConceptDailyBar; var err error
    for i := 0; i <= m.Retries; i++ {
        ths, err = GetConceptDailyThs(indexCode, kType, 1)
        if err == nil && len(ths) > 0 { return NormalizeConceptDaily(ths), nil }
        time.Sleep(m.MinWait)
    }
    var east []ConceptDailyBar; var err2 error
    for i := 0; i <= m.Retries; i++ {
        east, err2 = GetConceptDailyEast(indexCode, kType)
        if err2 == nil && len(east) > 0 { return NormalizeConceptDaily(east), nil }
        time.Sleep(m.MinWait)
    }
    if len(east) > 0 { return NormalizeConceptDaily(east), err }
    return NormalizeConceptDaily(ths), err
}

func (m *Market) GetMinute(indexCode string) ([]ConceptMinuteBar, error) {
    if indexCode == "" { return []ConceptMinuteBar{}, nil }
    var ths []ConceptMinuteBar; var err error
    for i := 0; i <= m.Retries; i++ {
        ths, err = GetConceptMinuteThs(indexCode)
        if err == nil && len(ths) > 0 { return NormalizeConceptMinute(ths), nil }
        time.Sleep(m.MinWait)
    }
    var east []ConceptMinuteBar; var err2 error
    for i := 0; i <= m.Retries; i++ {
        east, err2 = GetConceptMinuteEast(indexCode)
        if err2 == nil && len(east) > 0 { return NormalizeConceptMinute(east), nil }
        time.Sleep(m.MinWait)
    }
    if len(east) > 0 { return NormalizeConceptMinute(east), err }
    return NormalizeConceptMinute(ths), err
}

func (m *Market) GetCurrent(indexCode string, kType int) (ConceptCurrent, error) {
    if indexCode == "" { return ConceptCurrent{}, nil }
    var cur ConceptCurrent
    var err error
    for i := 0; i <= m.Retries; i++ {
        cur, err = GetConceptCurrentThs(indexCode, kType)
        if cur.IndexCode != "" { return NormalizeConceptCurrent(cur), nil }
        time.Sleep(m.MinWait)
    }
    if cur.IndexCode == "" {
        var mins []ConceptMinuteBar
        var merr error
        for i := 0; i <= m.Retries; i++ {
            mins, merr = GetConceptMinuteThs(indexCode)
            if merr == nil && len(mins) > 0 { break }
            time.Sleep(m.MinWait)
        }
        if len(mins) > 0 {
            last := mins[len(mins)-1]
            cur = ConceptCurrent{IndexCode: indexCode}
            cur.Price = last.Price
            cur.Change = last.Change
            cur.ChangePct = last.ChangePct
            cur.Volume = float64(last.Volume)
            cur.Amount = last.Amount
            cur.TradeTime = last.TradeTime
            cur.TradeDate = last.TradeDate
            return NormalizeConceptCurrent(cur), nil
        }
    }
    var cur2 ConceptCurrent
    var err2 error
    for i := 0; i <= m.Retries; i++ {
        cur2, err2 = GetConceptCurrentEast(indexCode)
        if cur2.IndexCode != "" { return NormalizeConceptCurrent(cur2), nil }
        time.Sleep(m.MinWait)
    }
    if cur2.IndexCode == "" {
        var mins2 []ConceptMinuteBar
        var merr2 error
        for i := 0; i <= m.Retries; i++ {
            mins2, merr2 = GetConceptMinuteEast(indexCode)
            if merr2 == nil && len(mins2) > 0 { break }
            time.Sleep(m.MinWait)
        }
        if len(mins2) > 0 {
            last := mins2[len(mins2)-1]
            cur2 = ConceptCurrent{IndexCode: indexCode}
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
            return NormalizeConceptCurrent(cur2), nil
        }
    }
    if err == nil { return NormalizeConceptCurrent(cur), err2 }
    return NormalizeConceptCurrent(cur2), err
}

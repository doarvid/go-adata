package capitalflow

import "time"

type Market struct { MinWait time.Duration; Retries int }

func New() *Market { return &Market{MinWait: 50 * time.Millisecond, Retries: 2} }

func (m *Market) GetMin(stockCode string) ([]FlowMin, error) {
    if stockCode == "" { return []FlowMin{}, nil }
    var east []FlowMin; var err error
    for i := 0; i <= m.Retries; i++ {
        east, err = GetStockCapitalFlowMinEast(stockCode)
        if err == nil && len(east) > 0 { return east, nil }
        time.Sleep(m.MinWait)
    }
    var bd []FlowMin; var err2 error
    for i := 0; i <= m.Retries; i++ {
        bd, err2 = GetStockCapitalFlowMinBaidu(stockCode)
        if err2 == nil && len(bd) > 0 { return bd, nil }
        time.Sleep(m.MinWait)
    }
    if len(bd) > 0 { return bd, err }
    return east, err
}

func (m *Market) GetDaily(stockCode, startDate, endDate string) ([]FlowDaily, error) {
    if stockCode == "" { return []FlowDaily{}, nil }
    var east []FlowDaily; var err error
    for i := 0; i <= m.Retries; i++ {
        east, err = GetStockCapitalFlowEast(stockCode, startDate, endDate)
        if err == nil && len(east) > 0 { return east, nil }
        time.Sleep(m.MinWait)
    }
    var bd []FlowDaily; var err2 error
    for i := 0; i <= m.Retries; i++ {
        bd, err2 = GetStockCapitalFlowBaidu(stockCode, startDate, endDate)
        if err2 == nil && len(bd) > 0 { return bd, nil }
        time.Sleep(m.MinWait)
    }
    if len(bd) > 0 { return bd, err }
    return east, err
}

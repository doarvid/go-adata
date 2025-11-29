package capitalflow

import "time"

type Market struct { MinWait time.Duration; Retries int }

func New() *Market { return &Market{MinWait: 50 * time.Millisecond, Retries: 2} }

func (m *Market) GetMin(stockCode string, wait time.Duration) ([]FlowMin, error) {
    if stockCode == "" { return []FlowMin{}, nil }
    if wait < m.MinWait { wait = m.MinWait }
    var east []FlowMin; var err error
    for i := 0; i <= m.Retries; i++ {
        east, err = GetStockCapitalFlowMinEast(stockCode, wait)
        if err == nil && len(east) > 0 { return east, nil }
        time.Sleep(wait)
    }
    var bd []FlowMin; var err2 error
    for i := 0; i <= m.Retries; i++ {
        bd, err2 = GetStockCapitalFlowMinBaidu(stockCode, wait)
        if err2 == nil && len(bd) > 0 { return bd, nil }
        time.Sleep(wait)
    }
    if len(bd) > 0 { return bd, err }
    return east, err
}

func (m *Market) GetDaily(stockCode, startDate, endDate string, wait time.Duration) ([]FlowDaily, error) {
    if stockCode == "" { return []FlowDaily{}, nil }
    if wait < m.MinWait { wait = m.MinWait }
    var east []FlowDaily; var err error
    for i := 0; i <= m.Retries; i++ {
        east, err = GetStockCapitalFlowEast(stockCode, startDate, endDate, wait)
        if err == nil && len(east) > 0 { return east, nil }
        time.Sleep(wait)
    }
    var bd []FlowDaily; var err2 error
    for i := 0; i <= m.Retries; i++ {
        bd, err2 = GetStockCapitalFlowBaidu(stockCode, startDate, endDate, wait)
        if err2 == nil && len(bd) > 0 { return bd, nil }
        time.Sleep(wait)
    }
    if len(bd) > 0 { return bd, err }
    return east, err
}

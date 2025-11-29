package stock

import (
	"time"

	"go-adata/pkg/adata/stock/market/capitalflow"
	"go-adata/pkg/adata/stock/market/conceptflow"
	"go-adata/pkg/adata/stock/market/conceptmarket"
	"go-adata/pkg/adata/stock/market/indexmarket"
	"go-adata/pkg/adata/stock/market/stockmarket"
)

type Stock struct {
	Market *stockmarket.Market
}

func New() *Stock { return &Stock{Market: stockmarket.NewMarket()} }

func (s *Stock) Daily(stockCode, startDate, endDate string, kType, adjustType int, wait time.Duration) ([]stockmarket.DailyBar, error) {
	return s.Market.GetDaily(stockCode, startDate, endDate, kType, adjustType, wait)
}

func (s *Stock) Minute(stockCode string, wait time.Duration) ([]stockmarket.MinuteBar, error) {
	return s.Market.GetMinute(stockCode, wait)
}

func (s *Stock) Bar(stockCode string, wait time.Duration) ([]stockmarket.TickBar, error) {
	return s.Market.GetBar(stockCode, wait)
}

func (s *Stock) Five(stockCode string, wait time.Duration) (stockmarket.Five, error) {
	return s.Market.GetFive(stockCode, wait)
}

func (s *Stock) Current(codeList []string, wait time.Duration) ([]stockmarket.CurrentQuote, error) {
	return s.Market.ListCurrent(codeList, wait)
}

// Python 风格方法名包装
func (s *Stock) GetMarket(stockCode, startDate, endDate string, kType, adjustType int, wait time.Duration) ([]stockmarket.DailyBar, error) {
	return s.Market.GetMarket(stockCode, startDate, endDate, kType, adjustType, wait)
}
func (s *Stock) GetMarketMin(stockCode string, wait time.Duration) ([]stockmarket.MinuteBar, error) {
	return s.Market.GetMarketMin(stockCode, wait)
}
func (s *Stock) GetMarketFive(stockCode string, wait time.Duration) (stockmarket.Five, error) {
	return s.Market.GetMarketFive(stockCode, wait)
}
func (s *Stock) GetMarketBar(stockCode string, wait time.Duration) ([]stockmarket.TickBar, error) {
	return s.Market.GetMarketBar(stockCode, wait)
}

func (s *Stock) IndexDaily(indexCode, startDate string, kType int, wait time.Duration) ([]indexmarket.IndexDailyBar, error) {
	im := indexmarket.New()
	return im.GetDaily(indexCode, startDate, kType, wait)
}

func (s *Stock) IndexMinute(indexCode string, wait time.Duration) ([]indexmarket.IndexMinuteBar, error) {
	im := indexmarket.New()
	return im.GetMinute(indexCode, wait)
}

func (s *Stock) IndexCurrent(indexCode string, wait time.Duration) (indexmarket.IndexCurrent, error) {
	im := indexmarket.New()
	return im.GetCurrent(indexCode, wait)
}

func (s *Stock) ConceptDaily(indexCode string, kType int, wait time.Duration) ([]conceptmarket.ConceptDailyBar, error) {
	cm := conceptmarket.New()
	return cm.GetDaily(indexCode, kType, wait)
}

func (s *Stock) ConceptMinute(indexCode string, wait time.Duration) ([]conceptmarket.ConceptMinuteBar, error) {
	cm := conceptmarket.New()
	return cm.GetMinute(indexCode, wait)
}

func (s *Stock) ConceptCurrent(indexCode string, wait time.Duration) (conceptmarket.ConceptCurrent, error) {
	cm := conceptmarket.New()
	return cm.GetCurrent(indexCode, 1, wait)
}

func (s *Stock) CapitalFlowMin(stockCode string, wait time.Duration) ([]capitalflow.FlowMin, error) {
	cf := capitalflow.New()
	return cf.GetMin(stockCode, wait)
}

func (s *Stock) CapitalFlowDaily(stockCode, startDate, endDate string, wait time.Duration) ([]capitalflow.FlowDaily, error) {
	cf := capitalflow.New()
	return cf.GetDaily(stockCode, startDate, endDate, wait)
}

func (s *Stock) ConceptFlowAll(daysType int, wait time.Duration) ([]conceptflow.ConceptFlow, error) {
	return conceptflow.ListConceptCapitalFlowEast(daysType, wait)
}

func (s *Stock) IndexCurrentBatch(indexCodes []string, wait time.Duration) ([]indexmarket.IndexCurrent, error) {
	if len(indexCodes) == 0 {
		return []indexmarket.IndexCurrent{}, nil
	}
	im := indexmarket.New()
	type pair struct {
		code string
		cur  indexmarket.IndexCurrent
	}
	max := 8
	sem := make(chan struct{}, max)
	ch := make(chan pair, len(indexCodes))
	for _, code := range indexCodes {
		c := code
		sem <- struct{}{}
		go func() {
			defer func() { <-sem }()
			var cur indexmarket.IndexCurrent
			var err error
			for i := 0; i <= im.Retries; i++ {
				cur, err = im.GetCurrent(c, wait)
				if cur.IndexCode != "" || err == nil {
					break
				}
				time.Sleep(wait)
			}
			ch <- pair{code: c, cur: cur}
		}()
	}
	got := make(map[string]indexmarket.IndexCurrent, len(indexCodes))
	for i := 0; i < len(indexCodes); i++ {
		p := <-ch
		got[p.code] = p.cur
	}
	out := make([]indexmarket.IndexCurrent, 0, len(indexCodes))
	for _, code := range indexCodes {
		if cur, ok := got[code]; ok {
			out = append(out, cur)
		}
	}
	return out, nil
}

func (s *Stock) ConceptCurrentBatch(indexCodes []string, wait time.Duration) ([]conceptmarket.ConceptCurrent, error) {
	if len(indexCodes) == 0 {
		return []conceptmarket.ConceptCurrent{}, nil
	}
	cm := conceptmarket.New()
	type pair struct {
		code string
		cur  conceptmarket.ConceptCurrent
	}
	max := 8
	sem := make(chan struct{}, max)
	ch := make(chan pair, len(indexCodes))
	for _, code := range indexCodes {
		c := code
		sem <- struct{}{}
		go func() {
			defer func() { <-sem }()
			var cur conceptmarket.ConceptCurrent
			var err error
			for i := 0; i <= cm.Retries; i++ {
				cur, err = cm.GetCurrent(c, 1, wait)
				if cur.IndexCode != "" || err == nil {
					break
				}
				time.Sleep(wait)
			}
			ch <- pair{code: c, cur: cur}
		}()
	}
	got := make(map[string]conceptmarket.ConceptCurrent, len(indexCodes))
	for i := 0; i < len(indexCodes); i++ {
		p := <-ch
		got[p.code] = p.cur
	}
	out := make([]conceptmarket.ConceptCurrent, 0, len(indexCodes))
	for _, code := range indexCodes {
		if cur, ok := got[code]; ok {
			out = append(out, cur)
		}
	}
	return out, nil
}

func (s *Stock) CapitalFlowMinBatch(stockCodes []string, wait time.Duration) ([][]capitalflow.FlowMin, error) {
	if len(stockCodes) == 0 {
		return [][]capitalflow.FlowMin{}, nil
	}
	max := 8
	sem := make(chan struct{}, max)
	type pair struct {
		idx  int
		rows []capitalflow.FlowMin
	}
	ch := make(chan pair, len(stockCodes))
	for i, code := range stockCodes {
		idx := i
		c := code
		sem <- struct{}{}
		go func() {
			defer func() { <-sem }()
			var rows []capitalflow.FlowMin
			var err error
			for r := 0; r < 3; r++ {
				rows, err = capitalflow.GetStockCapitalFlowMinEast(c, wait)
				if err == nil {
					break
				}
				time.Sleep(wait)
			}
			_ = err
			ch <- pair{idx: idx, rows: rows}
		}()
	}
	out := make([][]capitalflow.FlowMin, len(stockCodes))
	for i := 0; i < len(stockCodes); i++ {
		p := <-ch
		out[p.idx] = p.rows
	}
	return out, nil
}

func (s *Stock) CapitalFlowDailyBatch(stockCodes []string, startDate, endDate string, wait time.Duration) ([][]capitalflow.FlowDaily, error) {
	if len(stockCodes) == 0 {
		return [][]capitalflow.FlowDaily{}, nil
	}
	max := 8
	sem := make(chan struct{}, max)
	type pair struct {
		idx  int
		rows []capitalflow.FlowDaily
	}
	ch := make(chan pair, len(stockCodes))
	for i, code := range stockCodes {
		idx := i
		c := code
		sem <- struct{}{}
		go func() {
			defer func() { <-sem }()
			var rows []capitalflow.FlowDaily
			var err error
			for r := 0; r < 3; r++ {
				rows, err = capitalflow.GetStockCapitalFlowEast(c, startDate, endDate, wait)
				if err == nil {
					break
				}
				time.Sleep(wait)
			}
			_ = err
			ch <- pair{idx: idx, rows: rows}
		}()
	}
	out := make([][]capitalflow.FlowDaily, len(stockCodes))
	for i := 0; i < len(stockCodes); i++ {
		p := <-ch
		out[p.idx] = p.rows
	}
	return out, nil
}

// Python 风格别名
func (s *Stock) GetIndexCurrentBatch(indexCodes []string, wait time.Duration) ([]indexmarket.IndexCurrent, error) {
	return s.IndexCurrentBatch(indexCodes, wait)
}
func (s *Stock) GetConceptCurrentBatch(indexCodes []string, wait time.Duration) ([]conceptmarket.ConceptCurrent, error) {
	return s.ConceptCurrentBatch(indexCodes, wait)
}
func (s *Stock) GetCapitalFlowMinBatch(stockCodes []string, wait time.Duration) ([][]capitalflow.FlowMin, error) {
	return s.CapitalFlowMinBatch(stockCodes, wait)
}
func (s *Stock) GetCapitalFlowDailyBatch(stockCodes []string, startDate, endDate string, wait time.Duration) ([][]capitalflow.FlowDaily, error) {
	return s.CapitalFlowDailyBatch(stockCodes, startDate, endDate, wait)
}

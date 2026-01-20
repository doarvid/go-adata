package stockmarket

import (
	"context"
	"testing"

	"github.com/doarvid/go-adata/stock/info/tradecalendar"
)

func TestGetMarketDaily(t *testing.T) {
	stockCode := "002926"
	startDate := "2025-11-13"
	m := NewMarket()
	dailyBars, err := m.GetDaily(stockCode, startDate, tradecalendar.TradeDateNow(), KTypeDay, AdjustTypePre)
	if err != nil {
		t.Errorf("GetMarketDaily failed: %v", err)
	}
	if len(dailyBars) == 0 {
		t.Errorf("GetMarketDaily failed: empty daily bars")
	}
	for _, bar := range dailyBars {
		t.Logf("%v", bar)
	}
	t.Logf("total %d bars", len(dailyBars))
}

func TestGetMarketDailyBaidu(t *testing.T) {
	stockCode := "002926"
	startDate := "2025-11-13"
	dailyBars, err := NewMarket(WithDebug(true), WithProxy("http://192.168.31.100:20172")).GetDailyBaidu(context.Background(), stockCode, startDate, KTypeDay)
	if err != nil {
		t.Skipf("GetMarketDailyBaidu error: %v, skipping", err)
		return
	}
	if len(dailyBars) == 0 {
		t.Skipf("GetMarketDailyBaidu empty daily bars, skipping")
		return
	}
	for _, bar := range dailyBars {
		t.Logf("%v", bar)
	}
	t.Logf("total %d bars", len(dailyBars))
}

func TestGetMarketDailyEast(t *testing.T) {
	stockCode := "002926"
	startDate := "2025-11-13"
	endDate := "2025-11-18"
	dailyBars, err := NewMarket().GetDailyEast(context.Background(), stockCode, startDate, endDate, KTypeDay, 1)
	if err != nil {
		t.Errorf("GetMarketDailyEast failed: %v", err)
	}
	if len(dailyBars) == 0 {
		t.Errorf("GetMarketDailyEast failed: empty daily bars")
	}
	for _, bar := range dailyBars {
		t.Logf("%v", bar)
	}
	t.Logf("total %d bars", len(dailyBars))
}

func TestGetMarketMinuteEast(t *testing.T) {
	stockCode := "002926"
	minuteBars, err := NewMarket().GetMinuteEast(context.Background(), stockCode)
	if err != nil {
		t.Errorf("GetMarketMinuteEast failed: %v", err)
	}
	if len(minuteBars) == 0 {
		t.Errorf("GetMarketMinuteEast failed: empty minute bars")
	}
	for _, bar := range minuteBars {
		t.Logf("%v", bar)
	}
	t.Logf("total %d bars", len(minuteBars))
}

func TestGetMarketMinuteBaidu(t *testing.T) {
	stockCode := "002926"
	minuteBars, err := NewMarket().GetMinuteBaidu(context.Background(), stockCode)
	if err != nil {
		t.Skipf("GetMarketMinuteBaidu error: %v, skipping", err)
		return
	}
	if len(minuteBars) == 0 {
		t.Skipf("GetMarketMinuteBaidu empty minute bars, skipping")
		return
	}
	for _, bar := range minuteBars {
		t.Logf("%v", bar)
	}
	t.Logf("total %d bars", len(minuteBars))
}

func TestGetMarketCurrentSina(t *testing.T) {
	stockCode := "002926"
	currents, err := NewMarket().ListCurrentSina(context.Background(), []string{stockCode})
	if err != nil {
		t.Errorf("ListMarketCurrentSina failed: %v", err)
	}
	if len(currents) == 0 {
		t.Errorf("ListMarketCurrentSina failed: empty current")
	}
	for _, current := range currents {
		t.Logf("%v", current)
	}
	t.Logf("total %d current", len(currents))
}

func TestGetMarketCurrentQQ(t *testing.T) {
	stockCode := "002926"
	currents, err := NewMarket().ListCurrentQQ(context.Background(), []string{stockCode})
	if err != nil {
		t.Errorf("ListMarketCurrentQQ failed: %v", err)
	}
	if len(currents) == 0 {
		t.Errorf("ListMarketCurrentQQ failed: empty current")
	}
	for _, current := range currents {
		t.Logf("%v", current)
	}
	t.Logf("total %d current", len(currents))
}

func TestGetMarketFive(t *testing.T) {
	stockCode := "002926"
	five, err := NewMarket(WithDebug(true)).GetFiveBaidu(context.Background(), stockCode)
	if err != nil {
		t.Skipf("GetMarketFiveBaidu error: %v, skipping", err)
		return
	}
	if len(five.ShortName) == 0 {
		t.Skipf("empty five from baidu, skipping")
		return
	}
	t.Logf("five %+v", five)
}

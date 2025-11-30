package stockmarket

import (
	"testing"
	"time"
)

func TestGetMarketDailyBaidu(t *testing.T) {
	stockCode := "002926"
	startDate := "2025-11-13"
	kType := 1
	wait := 100 * time.Millisecond
	dailyBars, err := GetMarketDailyBaidu(stockCode, startDate, kType, wait)
	if err != nil {
		t.Errorf("GetMarketDailyBaidu failed: %v", err)
	}
	if len(dailyBars) == 0 {
		t.Errorf("GetMarketDailyBaidu failed: empty daily bars")
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
	kType := 1
	wait := 100 * time.Millisecond
	dailyBars, err := GetMarketDailyEast(stockCode, startDate, endDate, kType, 1, wait)
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
	wait := 100 * time.Millisecond
	minuteBars, err := GetMarketMinuteEast(stockCode, wait)
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
	wait := 100 * time.Millisecond
	minuteBars, err := GetMarketMinuteBaidu(stockCode, wait)
	if err != nil {
		t.Errorf("GetMarketMinuteBaidu failed: %v", err)
	}
	if len(minuteBars) == 0 {
		t.Errorf("GetMarketMinuteBaidu failed: empty minute bars")
	}
	for _, bar := range minuteBars {
		t.Logf("%v", bar)
	}
	t.Logf("total %d bars", len(minuteBars))
}

func TestGetMarketCurrentSina(t *testing.T) {
	stockCode := "002926"
	wait := 100 * time.Millisecond
	currents, err := ListMarketCurrentSina([]string{stockCode}, wait)
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
	wait := 100 * time.Millisecond
	currents, err := ListMarketCurrentQQ([]string{stockCode}, wait)
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

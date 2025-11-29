package indexmarket

import "testing"

func TestIndexThsDailyEmpty(t *testing.T) {
	if bars, _ := GetIndexDailyThs("", "2020-01-01", 1, 0); len(bars) != 0 {
		t.Fatalf("daily not empty")
	}
}

func TestIndexThsMinuteEmpty(t *testing.T) {
	if mins, _ := GetIndexMinuteThs("", 0); len(mins) != 0 {
		t.Fatalf("minute not empty")
	}
}

func TestIndexThsCurrentEmpty(t *testing.T) {
	if cur, _ := GetIndexCurrentThs("", 0); cur.Price != 0 || cur.TradeTime != "" {
		t.Fatalf("current not empty")
	}
}

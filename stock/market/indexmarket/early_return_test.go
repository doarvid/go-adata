package indexmarket

import "testing"

func TestIndexEarlyReturn(t *testing.T) {
    if bars, _ := GetIndexDailyEast("", "2020-01-01", 1, 0); len(bars) != 0 { t.Fatalf("daily not empty") }
    if mins, _ := GetIndexMinuteEast("", 0); len(mins) != 0 { t.Fatalf("minute not empty") }
    if cur, _ := GetIndexCurrentEast("", 0); (cur.Price != 0 || cur.TradeTime != "") { t.Fatalf("current not empty") }
}


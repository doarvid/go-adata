package indexmarket

import "testing"

func TestNormalizeIndexDailyClampAndDateFill(t *testing.T) {
    out := NormalizeIndexDaily([]IndexDailyBar{{Volume: -1, Amount: -2, TradeTime: "2024-01-01"}})
    if out[0].Volume != 0 || out[0].Amount != 0 { t.Fatalf("clamp failed") }
    if out[0].TradeDate != "2024-01-01" { t.Fatalf("trade date fill failed") }
}

func TestNormalizeIndexMinuteDateFill(t *testing.T) {
    out := NormalizeIndexMinute([]IndexMinuteBar{{TradeTime: "2024-02-03 10:00:00", Volume: 10}})
    if out[0].TradeDate != "2024-02-03" { t.Fatalf("minute date fill failed") }
}


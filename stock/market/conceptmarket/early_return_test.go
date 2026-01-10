package conceptmarket

import "testing"

func TestConceptEarlyReturn(t *testing.T) {
    if bars, _ := GetConceptDailyEast("", 1); len(bars) != 0 { t.Fatalf("daily not empty") }
    if mins, _ := GetConceptMinuteEast(""); len(mins) != 0 { t.Fatalf("minute not empty") }
    if cur, _ := GetConceptCurrentEast(""); (cur.Price != 0 || cur.TradeTime != "") { t.Fatalf("current not empty") }
}

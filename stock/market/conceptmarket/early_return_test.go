package conceptmarket

import "testing"

func TestConceptEarlyReturn(t *testing.T) {
    if bars, _ := GetConceptDailyEast("", 1, 0); len(bars) != 0 { t.Fatalf("daily not empty") }
    if mins, _ := GetConceptMinuteEast("", 0); len(mins) != 0 { t.Fatalf("minute not empty") }
    if cur, _ := GetConceptCurrentEast("", 0); (cur.Price != 0 || cur.TradeTime != "") { t.Fatalf("current not empty") }
}


package conceptmarket

import "testing"

func TestConceptThsDailyEmpty(t *testing.T) {
    if bars, _ := GetConceptDailyThs("", 1, 1, 0); len(bars) != 0 { t.Fatalf("daily not empty") }
}

func TestConceptThsMinuteEmpty(t *testing.T) {
    if mins, _ := GetConceptMinuteThs("", 0); len(mins) != 0 { t.Fatalf("minute not empty") }
}

func TestConceptThsCurrentEmpty(t *testing.T) {
    if cur, _ := GetConceptCurrentThs("", 1, 0); (cur.IndexCode != "" || cur.Price != 0) { t.Fatalf("current not empty") }
}

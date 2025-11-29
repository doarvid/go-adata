package stockmarket

import "testing"

func TestMarketEarlyReturn(t *testing.T) {
    m := NewMarket()
    if bars, _ := m.GetDaily("", "2020-01-01", "", 1, 1, 0); len(bars) != 0 { t.Fatalf("daily not empty") }
    if mins, _ := m.GetMinute("", 0); len(mins) != 0 { t.Fatalf("minute not empty") }
    if ticks, _ := m.GetBar("", 0); len(ticks) != 0 { t.Fatalf("bar not empty") }
    if five, _ := m.GetFive("", 0); (five.ShortName != "") { t.Fatalf("five not empty") }
}


package stockmarket

import "testing"

func TestNormalizeDailyClamp(t *testing.T) {
    out := NormalizeDaily([]DailyBar{{Volume: -1, Amount: -2, TradeTime: "2024-01-01"}})
    if out[0].Volume != 0 || out[0].Amount != 0 { t.Fatalf("clamp failed") }
    if out[0].TradeDate != "2024-01-01" { t.Fatalf("trade date fill failed") }
}

func TestNormalizeMinuteDateFill(t *testing.T) {
    out := NormalizeMinute([]MinuteBar{{TradeTime: "2024-02-03 10:00:00", Volume: 10}})
    if out[0].TradeDate != "2024-02-03" { t.Fatalf("minute date fill failed") }
}

func TestNormalizeTickClamp(t *testing.T) {
    out := NormalizeTick([]TickBar{{Price: -1, Volume: -5}})
    if out[0].Price != 0 || out[0].Volume != 0 { t.Fatalf("tick clamp failed") }
}

func TestNormalizeFiveClamp(t *testing.T) {
    f := NormalizeFive(Five{Sv1: -1, Sv2: -1, Sv3: -1, Sv4: -1, Sv5: -1, Bv1: -1, Bv2: -1, Bv3: -1, Bv4: -1, Bv5: -1})
    if f.Sv1 != 0 || f.Bv5 != 0 { t.Fatalf("five clamp failed") }
}

func TestNormalizeCurrentClamp(t *testing.T) {
    out := NormalizeCurrent([]CurrentQuote{{Volume: -1, Amount: -2}})
    if out[0].Volume != 0 || out[0].Amount != 0 { t.Fatalf("current clamp failed") }
}


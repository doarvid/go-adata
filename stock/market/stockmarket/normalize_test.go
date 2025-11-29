package stockmarket

import "testing"

func TestNormalizePassThrough(t *testing.T) {
    if len(NormalizeDaily([]DailyBar{})) != 0 { t.Fatalf("daily") }
    if len(NormalizeMinute([]MinuteBar{})) != 0 { t.Fatalf("minute") }
    if len(NormalizeTick([]TickBar{})) != 0 { t.Fatalf("tick") }
    if (NormalizeFive(Five{}).ShortName) != "" { t.Fatalf("five") }
    if len(NormalizeCurrent([]CurrentQuote{})) != 0 { t.Fatalf("current") }
}


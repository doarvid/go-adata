package stock

import (
    "os"
    "testing"
    "time"
)

func TestIndexDailyThreshold(t *testing.T) {
    if os.Getenv("ADATA_E2E") != "1" { t.Skip("external network") }
    s := New()
    bars, err := s.IndexDaily("000001", "2021-01-01", 1, 50*time.Millisecond)
    if err != nil { t.Fatalf("err: %v", err) }
    if len(bars) <= 30 { t.Fatalf("len=%d", len(bars)) }
}

func TestIndexMinuteThreshold(t *testing.T) {
    if os.Getenv("ADATA_E2E") != "1" { t.Skip("external network") }
    s := New()
    mins, err := s.IndexMinute("000001", 50*time.Millisecond)
    if err != nil { t.Fatalf("err: %v", err) }
    if len(mins) <= 2 { t.Fatalf("len=%d", len(mins)) }
}

func TestIndexCurrentNotEmpty(t *testing.T) {
    if os.Getenv("ADATA_E2E") != "1" { t.Skip("external network") }
    s := New()
    cur, err := s.IndexCurrent("000001", 50*time.Millisecond)
    if err != nil { t.Fatalf("err: %v", err) }
    if cur.TradeTime == "" && cur.Price == 0 { t.Fatalf("empty current") }
}

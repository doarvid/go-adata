package stock

import (
    "os"
    "testing"
    "time"
)

func TestConceptDailyThreshold(t *testing.T) {
    if os.Getenv("ADATA_E2E") != "1" { t.Skip("external network") }
    s := New()
    bars, err := s.ConceptDaily("886041", 1, 50*time.Millisecond)
    if err != nil { t.Fatalf("err: %v", err) }
    if len(bars) <= 30 { t.Fatalf("len=%d", len(bars)) }
}

func TestConceptMinuteThreshold(t *testing.T) {
    if os.Getenv("ADATA_E2E") != "1" { t.Skip("external network") }
    s := New()
    mins, err := s.ConceptMinute("886041", 50*time.Millisecond)
    if err != nil { t.Fatalf("err: %v", err) }
    if len(mins) <= 2 { t.Fatalf("len=%d", len(mins)) }
}

func TestConceptCurrentNotEmpty(t *testing.T) {
    if os.Getenv("ADATA_E2E") != "1" { t.Skip("external network") }
    s := New()
    cur, err := s.ConceptCurrent("886041", 50*time.Millisecond)
    if err != nil { t.Fatalf("err: %v", err) }
    if cur.TradeTime == "" && cur.Price == 0 { t.Fatalf("empty current") }
}

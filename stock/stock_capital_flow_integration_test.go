package stock

import (
    "os"
    "testing"
    "time"
)

func TestCapitalFlowMinThreshold(t *testing.T) {
    if os.Getenv("ADATA_E2E") != "1" { t.Skip("external network") }
    s := New()
    rows, err := s.CapitalFlowMin("000001", 50*time.Millisecond)
    if err != nil { t.Fatalf("err: %v", err) }
    if len(rows) == 0 { t.Fatalf("len=0") }
}

func TestCapitalFlowDailyThreshold(t *testing.T) {
    if os.Getenv("ADATA_E2E") != "1" { t.Skip("external network") }
    s := New()
    rows, err := s.CapitalFlowDaily("688403", "2020-01-01", "2024-12-31", 50*time.Millisecond)
    if err != nil { t.Fatalf("err: %v", err) }
    if len(rows) <= 200 { t.Fatalf("len=%d", len(rows)) }
}

func TestConceptFlowAllThreshold(t *testing.T) {
    if os.Getenv("ADATA_E2E") != "1" { t.Skip("external network") }
    s := New()
    rows, err := s.ConceptFlowAll(5, 50*time.Millisecond)
    if err != nil { t.Fatalf("err: %v", err) }
    if len(rows) <= 200 { t.Fatalf("len=%d", len(rows)) }
}

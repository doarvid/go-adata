package sentiment

import (
    "os"
    "testing"
    "time"
)

func e2eNorth(t *testing.T) bool { if os.Getenv("ADATA_E2E") == "1" { return true }; t.Skip("skip e2e; set ADATA_E2E=1 to enable"); return false }

func TestNorthFlowCurrentCount(t *testing.T) {
    if !e2eNorth(t) { return }
    r, err := NorthFlowCurrent(50 * time.Millisecond)
    if err != nil { t.Fatalf("err: %v", err) }
    if r.TradeTime == "" { t.Fatalf("empty current") }
}

func TestNorthFlowMinCount(t *testing.T) {
    if !e2eNorth(t) { return }
    rows, err := NorthFlowMin(50 * time.Millisecond)
    if err != nil { t.Fatalf("err: %v", err) }
    if len(rows) == 0 { t.Fatalf("len=0") }
}

func TestNorthFlowHistoryMinCount(t *testing.T) {
    if !e2eNorth(t) { return }
    rows, err := NorthFlow("2023-01-01", 50 * time.Millisecond)
    if err != nil { t.Fatalf("err: %v", err) }
    if len(rows) < 100 { t.Fatalf("len=%d", len(rows)) }
}


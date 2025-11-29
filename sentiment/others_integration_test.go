package sentiment

import (
    "os"
    "testing"
    "time"
)

func e2eOther(t *testing.T) bool { if os.Getenv("ADATA_E2E") == "1" { return true }; t.Skip("skip e2e; set ADATA_E2E=1 to enable"); return false }

func TestSecuritiesMarginMinCount(t *testing.T) {
    if !e2eOther(t) { return }
    rows, err := SecuritiesMargin("2020-01-01", 50 * time.Millisecond)
    if err != nil { t.Fatalf("err: %v", err) }
    if len(rows) < 250 { t.Fatalf("len=%d", len(rows)) }
}

func TestStockLiftingLastMonthCount(t *testing.T) {
    if !e2eOther(t) { return }
    rows, err := StockLiftingLastMonth(50 * time.Millisecond)
    if err != nil { t.Fatalf("err: %v", err) }
    if len(rows) <= 1 { t.Fatalf("len=%d", len(rows)) }
}


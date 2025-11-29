package sentiment

import (
    "os"
    "testing"
    "time"
)

func e2e(t *testing.T) bool { if os.Getenv("ADATA_E2E") == "1" { return true }; t.Skip("skip e2e; set ADATA_E2E=1 to enable"); return false }

func TestPopRank100EastCount(t *testing.T) {
    if !e2e(t) { return }
    rows, err := PopRank100East(50 * time.Millisecond)
    if err != nil { t.Fatalf("err: %v", err) }
    if len(rows) != 100 { t.Fatalf("len=%d", len(rows)) }
}

func TestHotRank100ThsCount(t *testing.T) {
    if !e2e(t) { return }
    rows, err := HotRank100Ths(50 * time.Millisecond)
    if err != nil { t.Fatalf("err: %v", err) }
    if len(rows) != 100 { t.Fatalf("len=%d", len(rows)) }
}

func TestHotConcept20ThsCount(t *testing.T) {
    if !e2e(t) { return }
    rows, err := HotConcept20Ths(1, 50 * time.Millisecond)
    if err != nil { t.Fatalf("err: %v", err) }
    if len(rows) != 20 { t.Fatalf("len=%d", len(rows)) }
}

func TestListAListDailyMinCount(t *testing.T) {
    if !e2e(t) { return }
    rows, err := ListAListDaily("2024-07-04", 50 * time.Millisecond)
    if err != nil { t.Fatalf("err: %v", err) }
    if len(rows) < 20 { t.Fatalf("len=%d", len(rows)) }
}

func TestGetAListInfoMinCount(t *testing.T) {
    if !e2e(t) { return }
    rows, err := GetAListInfo("600297", "2024-07-12", 50 * time.Millisecond)
    if err != nil { t.Fatalf("err: %v", err) }
    if len(rows) < 10 { t.Fatalf("len=%d", len(rows)) }
}


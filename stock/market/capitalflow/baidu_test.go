package capitalflow

import "testing"

func TestCapitalFlowBaiduMinEmpty(t *testing.T) {
    if rows, _ := GetStockCapitalFlowMinBaidu("", 0); len(rows) != 0 { t.Fatalf("min not empty") }
}

func TestCapitalFlowBaiduDailyEmpty(t *testing.T) {
    if rows, _ := GetStockCapitalFlowBaidu("", "2024-01-01", "2024-12-31", 0); len(rows) != 0 { t.Fatalf("daily not empty") }
}

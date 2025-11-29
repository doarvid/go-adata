package capitalflow

import "testing"

func TestCapitalFlowEarlyReturn(t *testing.T) {
    if mins, _ := GetStockCapitalFlowMinEast("", 0); len(mins) != 0 { t.Fatalf("min not empty") }
    if days, _ := GetStockCapitalFlowEast("", "2024-01-01", "2024-12-31", 0); len(days) != 0 { t.Fatalf("daily not empty") }
}


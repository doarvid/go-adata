package indexmarket

import (
	"context"
	"testing"
)

func TestGetIndexFive(t *testing.T) {
	indexCode := "000001"
	rows, err := NewIndexMarket().GetDailyEast(context.Background(), indexCode, "2025-10-21", 1)
	if err != nil {
		t.Errorf("GetIndexDailyEast failed: %v", err)
	}
	if len(rows) == 0 {
		t.Errorf("GetIndexDailyEast failed: empty five")
	}
	t.Logf("total %+v five", rows)
}

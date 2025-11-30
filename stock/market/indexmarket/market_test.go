package indexmarket

import (
	"testing"
	"time"
)

func TestGetIndexFive(t *testing.T) {
	indexCode := "000001"
	wait := 100 * time.Millisecond
	rows, err := GetIndexDailyEast(indexCode, "2025-10-21", 1, wait)
	if err != nil {
		t.Errorf("GetIndexDailyEast failed: %v", err)
	}
	if len(rows) == 0 {
		t.Errorf("GetIndexDailyEast failed: empty five")
	}
	t.Logf("total %+v five", rows)
}

package lifting

import (
	"context"
	"testing"
)

func TestStockLiftingLastMonth(t *testing.T) {
	sl := New()
	rows, err := sl.LastMonth(context.Background(), 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) == 0 {
		t.Fatalf("stock lifting last month not empty")
	}
	for _, row := range rows {
		t.Logf("%+v\n", row)
	}
	t.Logf("total %d rows", len(rows))
}


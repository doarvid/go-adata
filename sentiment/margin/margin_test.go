package margin

import (
	"context"
	"testing"
)

func TestSecuritiesMargin(t *testing.T) {
	m := New()
	rows, err := m.History(context.Background(), "2023-07-21", 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) == 0 {
		t.Fatal("no rows")
	}
	for _, row := range rows {
		t.Logf("%+v\n", row)
	}
	t.Logf("total %d rows", len(rows))
}


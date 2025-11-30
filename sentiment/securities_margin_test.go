package sentiment

import "testing"

func TestSecuritiesMargin(t *testing.T) {
	rows, err := SecuritiesMargin("2023-07-21", 0)
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

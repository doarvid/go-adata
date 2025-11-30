package sentiment

import "testing"

func TestStockLiftingLastMonth(t *testing.T) {
	rows, err := StockLiftingLastMonth(0)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) == 0 {
		t.Fatalf("stock lifting last month not empty")
	}
	for _, row := range rows {
		t.Logf("%+v\n", row)
	}
}

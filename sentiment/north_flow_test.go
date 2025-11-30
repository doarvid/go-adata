package sentiment

import "testing"

func TestNorthFlow(t *testing.T) {
	rows, err := NorthFlow("2025-11-12", 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) == 0 {
		t.Fatalf("north flow not empty")
	}
	for _, row := range rows {
		t.Logf("%+v\n", row)
	}
}

func TestNorthFlowMinute(t *testing.T) {
	rows, err := NorthFlowMin(10)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) == 0 {
		t.Fatalf("north flow minute not empty")
	}
	for _, row := range rows {
		t.Logf("%+v\n", row)
	}
}

func TestNorthFlowCurrent(t *testing.T) {
	row, err := NorthFlowCurrent(10)
	if err != nil {
		t.Fatal(err)
	}
	if row.TradeTime == "" {
		t.Fatalf("trade time not empty")
	}

	t.Logf("%+v\n", row)
}

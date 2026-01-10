package northflow

import (
	"context"
	"testing"
)

func TestNorthFlow(t *testing.T) {
	c := New()
	rows, err := c.History(context.Background(), "2025-11-12")
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
	c := New()
	rows, err := c.Minute(context.Background())
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
	c := New()
	row, err := c.Current(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if row.TradeTime == "" {
		t.Fatalf("trade time not empty")
	}
	t.Logf("%+v\n", row)
}

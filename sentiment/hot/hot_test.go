package hot

import (
	"context"
	"testing"
)

func TestPopRank100East(t *testing.T) {
	h := New()
	rows, err := h.Popular(context.Background(), 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) == 0 {
		t.Fatalf("pop rank 100 east not empty")
	}
	for _, row := range rows {
		t.Logf("%+v\n", row)
	}
}

func TestHotRank100Ths(t *testing.T) {
	h := New()
	rows, err := h.Stocks(context.Background(), 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) == 0 {
		t.Fatalf("hot rank 100 ths not empty")
	}
	for _, row := range rows {
		t.Logf("%+v\n", row)
	}
}

func TestHotConcept20Ths(t *testing.T) {
	h := New()
	rows, err := h.Concepts(context.Background(), PlateTypeConcept, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) == 0 {
		t.Fatalf("hot concept 20 ths not empty")
	}
	for _, row := range rows {
		t.Logf("%+v\n", row)
	}
}


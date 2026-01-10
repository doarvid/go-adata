package mine

import (
	"context"
	"testing"
)

func TestMineClearanceTDX(t *testing.T) {
	mc := New()
	rows, err := mc.EvaluateTDX(context.Background(), "600811", 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) == 0 {
		t.Fatalf("mine clearance tdx not empty")
	}
	for _, row := range rows {
		t.Logf("%+v\n", row)
	}
}


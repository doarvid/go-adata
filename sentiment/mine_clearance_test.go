package sentiment

import "testing"

func TestMineClearanceTDX(t *testing.T) {
	rows, err := MineClearanceTDX("600811", 0)
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

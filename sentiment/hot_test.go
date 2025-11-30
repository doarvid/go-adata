package sentiment

import (
	"testing"
)

func TestPopRank100East(t *testing.T) {
	rows, err := PopRank100East(0)
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
	rows, err := HotRank100Ths(0)
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
	rows, err := HotConcept20Ths(PlateTypeConcept, 0)
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

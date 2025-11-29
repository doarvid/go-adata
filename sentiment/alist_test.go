package sentiment

import (
	"fmt"
	"testing"
)

func TestGetAListInfo(t *testing.T) {
	if rows, _ := GetAListInfo("", "2024-01-01", 0); len(rows) != 0 {
		t.Fatalf("alist info not empty")
	}
}

func TestGetAListInfo2(t *testing.T) {
	rows, _ := GetAListInfo("600297", "2024-07-12", 10)
	if len(rows) == 0 {
		t.Fatalf("alist info not empty")
	}
	for _, row := range rows {
		fmt.Printf("row: %+v\n", row)
	}
}

package stockcode

import (
	"path/filepath"
	"runtime"
	"testing"
)

func TestMergeListDateFromCSV(t *testing.T) {
	codes := []StockCode{{StockCode: "000001"}}
	_, file, _, _ := runtime.Caller(0)
	base := filepath.Dir(file)
	root := filepath.Join(base, "..", "..", "..", "..", "..")
	csv := filepath.Join(root, "python", "adata", "adata", "stock", "cache", "code.csv")
	res := mergeListDateFromCSV(codes, csv)
	if res[0].ListDate == nil {
		t.Fatalf("list date not merged for 000001")
	}
}

func TestMarketRankEast(t *testing.T) {
	t.Skip("skip external network call in unit tests")
}

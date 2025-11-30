package stockcode

import (
	"testing"
	"time"
)

func TestAllStockCodes(t *testing.T) {
	codes, err := AllCode(100 * time.Millisecond)
	if err != nil {
		t.Fatalf("AllCode failed: %v", err)
	}
	if len(codes) == 0 {
		t.Fatalf("no stock codes")
	}
	t.Logf("total %d stock codes", len(codes))
}

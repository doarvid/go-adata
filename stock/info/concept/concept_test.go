package concept

import (
	"context"
	"testing"
	"time"
)

func TestLoadAllConceptCodesFromCSV(t *testing.T) {
	rows, err := LoadAllConceptCodesFromCSV()
	if err != nil {
		t.Fatalf("csv load failed: %v", err)
	}
	if len(rows) == 0 {
		t.Fatalf("csv empty")
	}
}

func TestNewConceptOptions(t *testing.T) {
	c := NewConcept(WithWait(50*time.Millisecond), WithRetries(1), WithUserAgent("test/ua"))
	if c == nil {
		t.Fatalf("NewConcept nil")
	}
}

func TestConceptAPI_Smoke(t *testing.T) {
	c := NewConcept(WithWait(50 * time.Millisecond), WithRetries(1))
	ctx := context.Background()
	codes, err := c.AllConceptCodesEast(ctx)
	if err != nil || len(codes) == 0 {
		t.Skipf("AllConceptCodesEast skip: err=%v len=%d", err, len(codes))
		return
	}
	infos, err2 := c.GetConceptEast(ctx, "600000")
	if err2 != nil || len(infos) == 0 {
		t.Skipf("GetConceptEast skip: err=%v len=%d", err2, len(infos))
	} else {
		if infos[0].StockCode == "" {
			t.Fatalf("GetConceptEast invalid")
		}
	}
	cs, err3 := c.ConstituentEast(ctx, codes[0].ConceptCode)
	if err3 != nil || len(cs) == 0 {
		t.Skipf("ConstituentEast skip: err=%v len=%d", err3, len(cs))
	} else {
		if cs[0].StockCode == "" {
			t.Fatalf("ConstituentEast invalid")
		}
	}
}

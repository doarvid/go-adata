package stock

import (
	"os"
	"testing"
	"time"
)

func TestGetMarketDailyThresholds(t *testing.T) {
	if os.Getenv("ADATA_E2E") != "1" {
		t.Skip("external network")
	}
	s := New()
	bars, err := s.GetMarket("000001", "2021-01-01", "", 1, 1, 50*time.Millisecond)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(bars) <= 300 {
		t.Fatalf("len=%d", len(bars))
	}
}

func TestGetMarketDailyOtherStock(t *testing.T) {
	if os.Getenv("ADATA_E2E") != "1" {
		t.Skip("external network")
	}
	s := New()
	bars, err := s.GetMarket("002824", "2007-01-01", "", 1, 1, 50*time.Millisecond)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(bars) <= 300 {
		t.Fatalf("len=%d", len(bars))
	}
}

func TestGetMarketWeeklyThreshold(t *testing.T) {
	if os.Getenv("ADATA_E2E") != "1" {
		t.Skip("external network")
	}
	s := New()
	bars, err := s.GetMarket("000001", "2021-01-01", "", 2, 1, 50*time.Millisecond)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(bars) <= 80 {
		t.Fatalf("len=%d", len(bars))
	}
}

func TestGetMarketMonthlyThreshold(t *testing.T) {
	if os.Getenv("ADATA_E2E") != "1" {
		t.Skip("external network")
	}
	s := New()
	bars, err := s.GetMarket("000001", "2021-01-01", "", 3, 1, 50*time.Millisecond)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(bars) <= 20 {
		t.Fatalf("len=%d", len(bars))
	}
}

func TestGetMarketMinThreshold(t *testing.T) {
	if os.Getenv("ADATA_E2E") != "1" {
		t.Skip("external network")
	}
	s := New()
	mins, err := s.GetMarketMin("000001", 50*time.Millisecond)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(mins) <= 1 {
		t.Fatalf("len=%d", len(mins))
	}
}

func TestListMarketCurrentThreshold(t *testing.T) {
	if os.Getenv("ADATA_E2E") != "1" {
		t.Skip("external network")
	}
	s := New()
	cur, err := s.Current([]string{"000001", "600001", "000795", "872925"}, 50*time.Millisecond)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(cur) <= 2 {
		t.Fatalf("len=%d", len(cur))
	}
}

func TestGetMarketFiveNotEmpty(t *testing.T) {
	if os.Getenv("ADATA_E2E") != "1" {
		t.Skip("external network")
	}
	s := New()
	five, err := s.GetMarketFive("000001", 50*time.Millisecond)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if five.ShortName == "" {
		t.Fatalf("empty five")
	}
}

func TestGetMarketBarThreshold(t *testing.T) {
	if os.Getenv("ADATA_E2E") != "1" {
		t.Skip("external network")
	}
	s := New()
	ticks, err := s.GetMarketBar("000001", 50*time.Millisecond)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(ticks) < 100 {
		t.Fatalf("len=%d", len(ticks))
	}
}

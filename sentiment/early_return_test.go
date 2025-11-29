package sentiment

import "testing"

func TestAListInfoEmpty(t *testing.T) {
    if rows, _ := GetAListInfo("", "2024-01-01", 0); len(rows) != 0 { t.Fatalf("alist info not empty") }
}

func TestNorthFlowCurrentSkip(t *testing.T) { t.Skip("external network") }
func TestNorthFlowMinSkip(t *testing.T) { t.Skip("external network") }
func TestNorthFlowDailySkip(t *testing.T) { t.Skip("external network") }
func TestSecuritiesMarginSkip(t *testing.T) { t.Skip("external network") }
func TestHotRankSkip(t *testing.T) { t.Skip("external network") }
func TestHotConceptSkip(t *testing.T) { t.Skip("external network") }
func TestPopRankSkip(t *testing.T) { t.Skip("external network") }
func TestStockLiftingSkip(t *testing.T) { t.Skip("external network") }

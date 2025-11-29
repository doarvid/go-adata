package stockmarket

import "testing"

func TestListMarketCurrentSinaEmpty(t *testing.T) {
    res, err := ListMarketCurrentSina([]string{}, 0)
    if err != nil { t.Fatalf("unexpected error: %v", err) }
    if len(res) != 0 { t.Fatalf("expected empty result, got %d", len(res)) }
}

func TestListMarketCurrentQQEmpty(t *testing.T) {
    res, err := ListMarketCurrentQQ([]string{}, 0)
    if err != nil { t.Fatalf("unexpected error: %v", err) }
    if len(res) != 0 { t.Fatalf("expected empty result, got %d", len(res)) }
}

func TestMarketListCurrentEmpty(t *testing.T) {
    m := NewMarket()
    res, err := m.ListCurrent([]string{}, 0)
    if err != nil { t.Fatalf("unexpected error: %v", err) }
    if len(res) != 0 { t.Fatalf("expected empty result, got %d", len(res)) }
}


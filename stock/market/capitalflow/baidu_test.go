package capitalflow

import (
	"context"
	"testing"
)

func TestCapitalFlowBaiduMinEmpty(t *testing.T) {
    c := NewClient()
    if rows, _ := c.MinutesBaidu(context.Background(), ""); len(rows) != 0 { t.Fatalf("min not empty") }
}

func TestCapitalFlowBaiduDailyEmpty(t *testing.T) {
    c := NewClient()
    if rows, _ := c.DailyBaidu(context.Background(), "", "2024-01-01", "2024-12-31"); len(rows) != 0 { t.Fatalf("daily not empty") }
}

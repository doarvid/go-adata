package indexmarket

import "testing"

func TestIndexBaiduDailyEmpty(t *testing.T) {
    if bars, _ := GetIndexDailyBaidu("", "2020-01-01", 1, 0); len(bars) != 0 { t.Fatalf("baidu daily not empty") }
}

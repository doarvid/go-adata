package codeutils

import "testing"

func TestCompileExchangeByStockCode(t *testing.T) {
    if got := CompileExchangeByStockCode("200039"); got != "200039.SZ" {
        t.Fatalf("want 200039.SZ, got %s", got)
    }
}

func TestGetExchangeByStockCode(t *testing.T) {
    if got := GetExchangeByStockCode("600000"); got != "SH" {
        t.Fatalf("want SH, got %s", got)
    }
}


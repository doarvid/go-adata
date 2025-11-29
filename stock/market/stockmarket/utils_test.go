package stockmarket

import "testing"

func TestParseF(t *testing.T) {
    cases := map[string]float64{
        "": 0,
        "--": 0,
        "1.23": 1.23,
        "-4.56%": -4.56,
        "+7.89": 7.89,
    }
    for in, exp := range cases {
        got := parseF(in)
        if got != exp {
            t.Fatalf("parseF(%q)=%v, exp=%v", in, got, exp)
        }
    }
}

func TestToInt64(t *testing.T) {
    if toInt64("123") != 123 { t.Fatalf("toInt64 123") }
    if toInt64("") != 0 { t.Fatalf("toInt64 empty") }
    if toInt64("-7") != -7 { t.Fatalf("toInt64 -7") }
}

func TestToString(t *testing.T) {
    if toString(" a ") != "a" { t.Fatalf("toString trim") }
}


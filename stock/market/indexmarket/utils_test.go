package indexmarket

import "testing"

func TestParseF(t *testing.T) {
    cases := map[string]float64{"":0, "--":0, "1.23":1.23, "-4.56%":-4.56, "+7.89":7.89}
    for in, exp := range cases {
        got := parseF(in)
        if got != exp { t.Fatalf("parseF(%q)=%v, exp=%v", in, got, exp) }
    }
}


package cache

import (
    "os"
    "testing"
)

func TestGetCodeCSVPathExists(t *testing.T) {
    p := GetCodeCSVPath()
    if _, err := os.Stat(p); err != nil {
        t.Fatalf("code.csv not found at %s", p)
    }
}

func TestGetCalendarCSVPathExists(t *testing.T) {
    p := GetCalendarCSVPath(2021)
    if _, err := os.Stat(p); err != nil {
        t.Fatalf("calendar 2021 not found at %s", p)
    }
}

func TestCalendarYearsContains(t *testing.T) {
    yrs := CalendarYears()
    if len(yrs) == 0 { t.Fatalf("empty years") }
    found := false
    for _, y := range yrs { if y == 2025 { found = true; break } }
    if !found { t.Fatalf("2025 not in years") }
}

func TestLoadIndexCodeRelTHS(t *testing.T) {
    m, err := LoadIndexCodeRelTHS()
    if err != nil { t.Fatalf("load mapping error: %v", err) }
    // sample key from file
    if m["000819"] == "" { t.Fatalf("expected mapping for 000819") }
}


package tradecalendar

import (
	"strconv"
	"strings"
	"testing"
)

func TestTradeCalendarStructureAndRange(t *testing.T) {
	for year := 2000; year <= 2025; year++ {
		t.Run(strconv.Itoa(year), func(t *testing.T) {
			days, err := TradeCalendar(year)
			if err != nil {
				t.Fatalf("err: %v", err)
			}
			if len(days) == 0 {
				t.Skipf("no data for %d", year)
				return
			}
			for i, d := range days {
				if d.TradeDate == "" {
					t.Fatalf("empty trade_date at %d", i)
				}
				parts := strings.Split(d.TradeDate, "-")
				if len(parts) != 3 {
					t.Fatalf("bad date format: %s", d.TradeDate)
				}
				if parts[0] != strconv.Itoa(year) {
					t.Fatalf("year mismatch: %s != %d", parts[0], year)
				}
				if d.TradeStatus != 0 && d.TradeStatus != 1 {
					t.Fatalf("trade_status out of range: %d", d.TradeStatus)
				}
				if d.DayWeek < 1 || d.DayWeek > 7 {
					t.Fatalf("day_week out of range: %d", d.DayWeek)
				}
			}
		})
	}
}

func TestAvailTradeDayN(t *testing.T) {
	days, err := AvailTradeDayN(365)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(days) == 0 {
		t.Skipf("no data for latest year")
		return
	}
}

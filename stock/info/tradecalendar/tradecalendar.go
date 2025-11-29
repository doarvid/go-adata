package tradecalendar

import (
    "encoding/csv"
    "github.com/doarvid/go-adata/stock/cache"
    "io"
    "os"
    "strconv"
)

type Day struct {
	TradeDate   string `json:"trade_date"`
	TradeStatus int    `json:"trade_status"`
	DayWeek     int    `json:"day_week"`
}

func TradeCalendar(year int) ([]Day, error) {
	p := cache.GetCalendarCSVPath(year)
	f, err := os.Open(p)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	r := csv.NewReader(f)
	// skip header
	if _, err := r.Read(); err != nil {
		return nil, err
	}
	var out []Day
	for {
		rec, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		status, _ := strconv.Atoi(rec[1])
		week, _ := strconv.Atoi(rec[2])
		out = append(out, Day{TradeDate: rec[0], TradeStatus: status, DayWeek: week})
	}
	return out, nil
}

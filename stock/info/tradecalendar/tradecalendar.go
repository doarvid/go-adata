package tradecalendar

import (
	"embed"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"slices"
	"strconv"
	"strings"
	"time"
)

//go:embed calendar/*
var staticFiles embed.FS

type Day struct {
	TradeDate   string `json:"trade_date"`
	TradeStatus int    `json:"trade_status"`
	DayWeek     int    `json:"day_week"`
}

func TradeCalendar(year int) ([]Day, error) {
	file, err := staticFiles.Open(fmt.Sprintf("calendar/%d.csv", year))
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	r := csv.NewReader(file)
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

func TradeDate(t time.Time) string {
	return t.Format("2006-01-02")
}
func TradeDateNow() string {
	return TradeDate(time.Now())
}
func tradeDayN(days int) ([]Day, error) {
	years := CalendarYears()
	var tradeDates []Day
	for i := 0; i < len(years); i++ {
		tdays, err := TradeCalendar(years[len(years)-1-i])
		if err != nil {
			return nil, err
		}
		tradeDates = append(tradeDates, tdays...)
		if len(tradeDates) >= days+365 {
			break
		}
	}
	slices.SortFunc(tradeDates, func(a, b Day) int {
		return strings.Compare(a.TradeDate, b.TradeDate)
	})
	tradeDate := TradeDate(time.Now())
	var ret []Day
	for i := len(tradeDates) - 1; i > 0; i-- {
		if tradeDates[i].TradeDate == tradeDate {
			for j := i; j > 0 && len(ret) < days; j-- {
				if tradeDates[j].TradeStatus == 1 {
					ret = append(ret, tradeDates[j])
				}
			}
			break
		}
	}
	slices.SortFunc(ret, func(a, b Day) int {
		return strings.Compare(a.TradeDate, b.TradeDate)
	})
	return ret, nil
}

func TradeDayN(days int) ([]string, error) {
	tradeDates, err := tradeDayN(days)
	if err != nil {
		return nil, err
	}
	var ret []string
	for _, tradeDate := range tradeDates {
		ret = append(ret, tradeDate.TradeDate)
	}
	return ret, nil
}

func CalendarYears() []int {
	return []int{2004, 2005, 2006, 2007, 2008, 2009, 2010, 2011, 2012, 2013, 2014, 2015, 2016, 2017, 2018, 2019, 2020, 2021, 2022, 2023, 2024, 2025, 2026}
}

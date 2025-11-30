package tradecalendar

import (
	"embed"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"strconv"
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

func CalendarYears() []int {
	return []int{2004, 2005, 2006, 2007, 2008, 2009, 2010, 2011, 2012, 2013, 2014, 2015, 2016, 2017, 2018, 2019, 2020, 2021, 2022, 2023, 2024, 2025}
}

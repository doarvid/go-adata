package utils

import "time"

func GetNDaysDate(days int, fmt string) string {
    if fmt == "" {
        fmt = "2006-01-02"
    }
    current := time.Now()
    target := current.AddDate(0, 0, days)
    return target.Format(fmt)
}

func GetCurTime(fmt string) string {
    if fmt == "" {
        fmt = "2006-01-02 15:04:05"
    }
    return time.Now().Format(fmt)
}


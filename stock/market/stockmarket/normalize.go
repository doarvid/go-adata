package stockmarket

import "strings"

func NormalizeDaily(in []DailyBar) []DailyBar {
    for i := range in {
        if in[i].Volume < 0 { in[i].Volume = 0 }
        if in[i].Amount < 0 { in[i].Amount = 0 }
        if in[i].TradeDate == "" { in[i].TradeDate = in[i].TradeTime }
    }
    return in
}

func NormalizeMinute(in []MinuteBar) []MinuteBar {
    for i := range in {
        if in[i].Volume < 0 { in[i].Volume = 0 }
        if in[i].Amount < 0 { in[i].Amount = 0 }
        if in[i].TradeDate == "" {
            tt := in[i].TradeTime
            if idx := strings.Index(tt, " "); idx > 0 { in[i].TradeDate = tt[:idx] } else { in[i].TradeDate = tt }
        }
    }
    return in
}

func NormalizeTick(in []TickBar) []TickBar {
    for i := range in {
        if in[i].Volume < 0 { in[i].Volume = 0 }
        if in[i].Price < 0 { in[i].Price = 0 }
    }
    return in
}

func NormalizeFive(in Five) Five {
    if in.Sv1 < 0 { in.Sv1 = 0 }
    if in.Sv2 < 0 { in.Sv2 = 0 }
    if in.Sv3 < 0 { in.Sv3 = 0 }
    if in.Sv4 < 0 { in.Sv4 = 0 }
    if in.Sv5 < 0 { in.Sv5 = 0 }
    if in.Bv1 < 0 { in.Bv1 = 0 }
    if in.Bv2 < 0 { in.Bv2 = 0 }
    if in.Bv3 < 0 { in.Bv3 = 0 }
    if in.Bv4 < 0 { in.Bv4 = 0 }
    if in.Bv5 < 0 { in.Bv5 = 0 }
    return in
}

func NormalizeCurrent(in []CurrentQuote) []CurrentQuote {
    for i := range in {
        if in[i].Volume < 0 { in[i].Volume = 0 }
        if in[i].Amount < 0 { in[i].Amount = 0 }
    }
    return in
}

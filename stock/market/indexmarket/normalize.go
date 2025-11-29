package indexmarket

func NormalizeIndexDaily(in []IndexDailyBar) []IndexDailyBar {
    for i := range in {
        if in[i].Volume < 0 { in[i].Volume = 0 }
        if in[i].Amount < 0 { in[i].Amount = 0 }
        if in[i].TradeDate == "" { in[i].TradeDate = in[i].TradeTime }
    }
    return in
}

func NormalizeIndexMinute(in []IndexMinuteBar) []IndexMinuteBar {
    for i := range in {
        if in[i].Volume < 0 { in[i].Volume = 0 }
        if in[i].Amount < 0 { in[i].Amount = 0 }
        if in[i].TradeDate == "" { tt := in[i].TradeTime; if len(tt) >= 10 { in[i].TradeDate = tt[:10] } else { in[i].TradeDate = tt } }
    }
    return in
}

func NormalizeIndexCurrent(in IndexCurrent) IndexCurrent { return in }


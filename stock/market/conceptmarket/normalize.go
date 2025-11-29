package conceptmarket

func NormalizeConceptDaily(in []ConceptDailyBar) []ConceptDailyBar {
    for i := range in {
        if in[i].Volume < 0 { in[i].Volume = 0 }
        if in[i].Amount < 0 { in[i].Amount = 0 }
        if in[i].TradeDate == "" { in[i].TradeDate = in[i].TradeTime }
    }
    return in
}

func NormalizeConceptMinute(in []ConceptMinuteBar) []ConceptMinuteBar {
    for i := range in {
        if in[i].Volume < 0 { in[i].Volume = 0 }
        if in[i].Amount < 0 { in[i].Amount = 0 }
        if in[i].TradeDate == "" { tt := in[i].TradeTime; if len(tt) >= 10 { in[i].TradeDate = tt[:10] } else { in[i].TradeDate = tt } }
    }
    return in
}

func NormalizeConceptCurrent(in ConceptCurrent) ConceptCurrent { return in }


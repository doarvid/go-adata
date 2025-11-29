package codeutils

var exchangeSuffix = map[string]string{
    "00": ".SZ",
    "20": ".SZ",
    "30": ".SZ",
    "43": ".BJ",
    "60": ".SH",
    "68": ".SH",
    "83": ".BJ",
    "87": ".BJ",
    "90": ".SH",
    "92": ".BJ",
}

func CompileExchangeByStockCode(stockCode string) string {
    if len(stockCode) < 2 {
        return stockCode
    }
    prefix := stockCode[:2]
    if suffix, ok := exchangeSuffix[prefix]; ok {
        return stockCode + suffix
    }
    return stockCode
}

func GetExchangeByStockCode(stockCode string) string {
    if len(stockCode) < 2 {
        return ""
    }
    prefix := stockCode[:2]
    if suffix, ok := exchangeSuffix[prefix]; ok && len(suffix) > 1 {
        return suffix[1:]
    }
    return ""
}


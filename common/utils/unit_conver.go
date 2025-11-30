package utils

import (
    "regexp"
    "strconv"
    "strings"
)

func ConvertToYuan(input map[string]string) map[string]float64 {
    unit := map[string]float64{"亿": 100000000, "万": 10000}
    re := regexp.MustCompile(`([-+]?\d*\.\d+|\d+)([亿万]?)`)
    out := make(map[string]float64, len(input))
    for k, v := range input {
        vv := strings.ReplaceAll(v, "元", "")
        m := re.FindStringSubmatch(vv)
        if len(m) >= 3 {
            num, _ := strconv.ParseFloat(m[1], 64)
            mult := unit[m[2]]
            if mult == 0 { mult = 1 }
            out[k] = num * mult
        } else {
            n, _ := strconv.ParseFloat(vv, 64)
            out[k] = n
        }
    }
    return out
}


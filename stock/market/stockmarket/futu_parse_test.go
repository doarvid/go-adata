package stockmarket

import (
	"fmt"
	"math"
	"strings"
	"testing"
)

// almostEqual 判断两个浮点数是否近似相等
func almostEqual(a, b float64) bool {
	return math.Abs(a-b) < 0.01
}

// almostEqualPtr 判断指针 float64 值是否近似相等
func almostEqualPtr(ptr *float64, expected float64) bool {
	if ptr == nil {
		return false
	}
	return math.Abs(*ptr-expected) < 0.01
}

func TestParseFutuHTML(t *testing.T) {
	// 实际的 HTML 片段 - 包含所有字段
	htmlContent := `<div class="content-main" data-v-47a1769a="">
<a href="/stock/920069-BJ" target="_blank" class="list-item" data-v-92fc86ce="">
  <div class="fix-left" data-v-92fc86ce="">
    <span class="order" data-v-92fc86ce="">1</span>
    <span title="920069" class="code ellipsis" data-v-92fc86ce="">920069</span>
    <span title="N 普昂" class="name ellipsis" data-v-92fc86ce="">N 普昂</span>
  </div>
  <div class="middle" data-v-92fc86ce="">
    <span title="43.48" class="value ellipsis direct-up" data-v-92fc86ce="">43.48</span>
    <span title="+25.10" class="value ellipsis direct-up" data-v-92fc86ce="">+25.10</span>
    <span title="+136.56%" class="value ellipsis direct-up" data-v-92fc86ce="">+136.56%</span>
    <span title="968.34 万" class="value ellipsis" data-v-92fc86ce="">968.34 万</span>
    <span title="4.10 亿" class="value ellipsis" data-v-92fc86ce="">4.10 亿</span>
    <span title="23.01 亿" class="value ellipsis" data-v-92fc86ce="">23.01 亿</span>
    <span title="6.43 亿" class="value ellipsis" data-v-92fc86ce="">6.43 亿</span>
    <span title="5291.40 万" class="value ellipsis" data-v-92fc86ce="">5291.40 万</span>
    <span title="1477.77 万" class="value ellipsis" data-v-92fc86ce="">1477.77 万</span>
    <span title="+136.56%" class="value ellipsis direct-up" data-v-92fc86ce="">+136.56%</span>
    <span title="+136.56%" class="value ellipsis direct-up" data-v-92fc86ce="">+136.56%</span>
    <span title="+136.56%" class="value ellipsis direct-up" data-v-92fc86ce="">+136.56%</span>
  </div>
</a>
<a href="/stock/001313-SZ" target="_blank" class="list-item" data-v-92fc86ce="">
  <div class="fix-left" data-v-92fc86ce="">
    <span class="order" data-v-92fc86ce="">2</span>
    <span title="001313" class="code ellipsis" data-v-92fc86ce="">001313</span>
    <span title="粤海饲料" class="name ellipsis" data-v-92fc86ce="">粤海饲料</span>
  </div>
  <div class="middle" data-v-92fc86ce="">
    <span title="7.56" class="value ellipsis direct-up" data-v-92fc86ce="">7.56</span>
    <span title="+0.69" class="value ellipsis direct-up" data-v-92fc86ce="">+0.69</span>
    <span title="+10.04%" class="value ellipsis direct-up" data-v-92fc86ce="">+10.04%</span>
    <span title="1825.99 万" class="value ellipsis" data-v-92fc86ce="">1825.99 万</span>
    <span title="1.35 亿" class="value ellipsis" data-v-92fc86ce="">1.35 亿</span>
    <span title="52.92 亿" class="value ellipsis" data-v-92fc86ce="">52.92 亿</span>
    <span title="52.85 亿" class="value ellipsis" data-v-92fc86ce="">52.85 亿</span>
  </div>
</a>
</div>`

	quotes, err := parseFutuHTML(htmlContent)
	if err != nil {
		t.Fatalf("parseFutuHTML 失败：%v", err)
	}

	if len(quotes) != 2 {
		t.Fatalf("期望解析 2 条数据，实际：%d", len(quotes))
	}

	// 验证第一条数据
	q1 := quotes[0]
	if q1.StockCode != "920069" {
		t.Errorf("股票代码期望 920069，实际：%s", q1.StockCode)
	}
	if q1.ShortName != "N 普昂" {
		t.Errorf("股票名称期望 N 普昂，实际：%s", q1.ShortName)
	}
	if !almostEqual(q1.Price, 43.48) {
		t.Errorf("价格期望 43.48，实际：%f", q1.Price)
	}
	if !almostEqual(q1.Change, 25.10) {
		t.Errorf("涨跌额期望 25.10，实际：%f", q1.Change)
	}
	if !almostEqual(q1.ChangePct, 136.56) {
		t.Errorf("涨跌幅期望 136.56，实际：%f", q1.ChangePct)
	}
	if !almostEqual(q1.Volume, 9683400) {
		t.Errorf("成交量期望 9683400，实际：%f", q1.Volume)
	}
	if !almostEqual(q1.Amount, 410000000) {
		t.Errorf("成交额期望 410000000，实际：%f", q1.Amount)
	}
	if !almostEqualPtr(q1.TotalCap, 2301000000) {
		t.Errorf("总市值期望 2301000000，实际：%v", q1.TotalCap)
	}
	if !almostEqualPtr(q1.FloatCap, 643000000) {
		t.Errorf("流通市值期望 643000000，实际：%v", q1.FloatCap)
	}
	// 验证总股本和流通股本
	if !almostEqualPtr(q1.TotalShares, 52914000) {
		t.Errorf("总股本期望 52914000，实际：%v", q1.TotalShares)
	}
	if !almostEqualPtr(q1.FloatShares, 14777700) {
		t.Errorf("流通股本期望 14777700，实际：%v", q1.FloatShares)
	}
	// 验证 5 日/10 日/20 日涨跌幅
	if !almostEqualPtr(q1.Change5d, 136.56) {
		t.Errorf("5 日涨跌幅期望 136.56，实际：%v", q1.Change5d)
	}
	if !almostEqualPtr(q1.Change10d, 136.56) {
		t.Errorf("10 日涨跌幅期望 136.56，实际：%v", q1.Change10d)
	}
	if !almostEqualPtr(q1.Change20d, 136.56) {
		t.Errorf("20 日涨跌幅期望 136.56，实际：%v", q1.Change20d)
	}

	// 验证第二条数据
	q2 := quotes[1]
	if q2.StockCode != "001313" {
		t.Errorf("股票代码期望 001313，实际：%s", q2.StockCode)
	}
	if q2.ShortName != "粤海饲料" {
		t.Errorf("股票名称期望 粤海饲料，实际：%s", q2.ShortName)
	}
	if !almostEqual(q2.Price, 7.56) {
		t.Errorf("价格期望 7.56，实际：%f", q2.Price)
	}
	if !almostEqualPtr(q2.TotalCap, 5292000000) {
		t.Errorf("总市值期望 5292000000，实际：%v", q2.TotalCap)
	}
	if !almostEqualPtr(q2.FloatCap, 5285000000) {
		t.Errorf("流通市值期望 5285000000，实际：%v", q2.FloatCap)
	}

	fmt.Println("解析测试通过！")
	for _, q := range quotes {
		totalCap := "N/A"
		if q.TotalCap != nil {
			totalCap = fmt.Sprintf("%.2f 亿", *q.TotalCap/100000000)
		}
		floatCap := "N/A"
		if q.FloatCap != nil {
			floatCap = fmt.Sprintf("%.2f 亿", *q.FloatCap/100000000)
		}
		totalShares := "N/A"
		if q.TotalShares != nil {
			totalShares = fmt.Sprintf("%.2f 万", *q.TotalShares/10000)
		}
		floatShares := "N/A"
		if q.FloatShares != nil {
			floatShares = fmt.Sprintf("%.2f 万", *q.FloatShares/10000)
		}
		change5d := "N/A"
		if q.Change5d != nil {
			change5d = fmt.Sprintf("%.2f%%", *q.Change5d)
		}
		fmt.Printf("%s %s: %.2f (总市值:%s 流通:%s 总股本:%s 流通股本:%s 5 日:%s)\n",
			q.StockCode, q.ShortName, q.Price,
			totalCap, floatCap, totalShares, floatShares, change5d)
	}
}

func TestParseWithUnit(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"968.34 万", 9683400},
		{"4.10 亿", 410000000},
		{"1.35 亿", 135000000},
		{"1825.99 万", 18259900},
		{"1000", 1000},
		{"1,000", 1000},
		{"5.5 千", 5500},
		{"", 0},
	}

	for _, tt := range tests {
		result := parseWithUnit(tt.input)
		if !almostEqual(result, tt.expected) {
			t.Errorf("parseWithUnit(%q) = %f, 期望 %f", tt.input, result, tt.expected)
		}
	}

	fmt.Println("单位解析测试通过！")
}

func TestParseF(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"43.48", 43.48},
		{"+25.10", 25.10},
		{"-0.02", -0.02},
		{"+136.56%", 136.56},
		{"", 0},
		{"  3.14  ", 3.14},
	}

	for _, tt := range tests {
		// 移除+和%符号
		cleaned := strings.ReplaceAll(tt.input, "+", "")
		cleaned = strings.ReplaceAll(cleaned, "%", "")
		result := parseF(cleaned)
		if result != tt.expected {
			t.Errorf("parseF(%q) = %f, 期望 %f", tt.input, result, tt.expected)
		}
	}

	fmt.Println("数值解析测试通过！")
}

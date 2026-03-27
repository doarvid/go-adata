package stockmarket

import (
	"context"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
)

// WithRemoteChrome 设置 Chrome 远程调试地址
func WithRemoteChrome(url string) MarketOpt {
	return func(m *Market) {
		m.remoteChrome = url
	}
}

// AllFutu 爬取富途牛牛 A 股行情数据，支持翻页
// 使用 Chrome 远程调试方式，需要预先启动带调试端口的 Chrome
// 启动命令示例：chrome --remote-debugging-port=9222
func (m *Market) AllFutu(ctx context.Context, callback func(stock *CurrentQuote)) error {
	debugURL := m.remoteChrome
	if debugURL == "" {
		debugURL = "http://localhost:9222"
	}

	allocCtx, cancel := chromedp.NewRemoteAllocator(context.Background(), debugURL)
	defer cancel()

	ctx, cancel = chromedp.NewContext(allocCtx)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	url := "https://www.futunn.com/quote/cn/stock-list/all-cn-stocks/top-market-cap"
	contentSelector := "#view-page > div.quote-page.router-page > div > section > div > div.stock-list > div.list-table > section > div.content-main"
	nextPageSelector := "#view-page > div.quote-page.router-page > div > section > div > div.stock-list > div.list-table > div > span.item.next > i"

	var totalProcessed int

	for page := 1; ; page++ {
		var htmlContent string

		if page == 1 {
			err := chromedp.Run(ctx,
				chromedp.Navigate(url),
				chromedp.Sleep(30*time.Second),
				chromedp.WaitVisible(contentSelector, chromedp.ByQuery),
				chromedp.OuterHTML(contentSelector, &htmlContent, chromedp.NodeVisible),
			)
			if err != nil {
				if totalProcessed == 0 {
					return err
				}
				break
			}
		} else {
			var hasNext bool
			err := chromedp.Run(ctx,
				chromedp.WaitVisible(nextPageSelector, chromedp.ByQuery, chromedp.AtLeast(0)),
				chromedp.Evaluate(`document.querySelector('`+nextPageSelector+`') !== null && document.querySelector('`+nextPageSelector+`').offsetParent !== null`, &hasNext),
			)
			if err != nil || !hasNext {
				break
			}

			err = chromedp.Run(ctx,
				chromedp.Click(nextPageSelector, chromedp.ByQuery),
				chromedp.Sleep(2*time.Second),
				chromedp.WaitVisible(contentSelector, chromedp.ByQuery),
				chromedp.OuterHTML(contentSelector, &htmlContent, chromedp.NodeVisible),
			)
			if err != nil {
				break
			}
		}

		quotes, err := parseFutuHTML(htmlContent)
		if err != nil {
			continue
		}

		if len(quotes) == 0 {
			break
		}

		for _, quote := range quotes {
			callback(quote)
			totalProcessed++
		}

		select {
		case <-ctx.Done():
			return nil
		default:
		}
		time.Sleep(time.Second * 30)
	}

	return nil
}

// parseFutuHTML 解析富途 HTML 内容
// HTML 结构：a.list-item > div.fix-left(span.code, span.name) + div.middle(span.value...)
// span.value 字段顺序 (根据实际 HTML):
// [0] 最新价 [1] 涨跌额 [2] 涨跌幅 [3] 成交量 [4] 成交额
// [5] 总市值 [6] 流通市值 [7] 总股本 [8] 流通股本
// [9] 5 日涨跌幅 [10] 10 日涨跌幅 [11] 20 日涨跌幅 [12] 60 日涨跌幅
// [13] 120 日涨跌幅 [14] 250 日涨跌幅 [15] 年初至今涨跌幅
func parseFutuHTML(htmlContent string) ([]*CurrentQuote, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		return nil, err
	}

	var quotes []*CurrentQuote

	// 每行数据是 a.list-item
	doc.Find("a.list-item").Each(func(i int, s *goquery.Selection) {
		quote := &CurrentQuote{}

		// 获取股票代码 - span.code 的 title 属性或文本
		code := strings.TrimSpace(s.Find("span.code").Text())
		if code == "" {
			return
		}
		quote.StockCode = code

		// 获取股票名称 - span.name 的 title 属性或文本
		quote.ShortName = strings.TrimSpace(s.Find("span.name").Text())

		// div.middle 包含多个 span.value，按顺序抽取所有可用字段
		values := s.Find("div.middle span.value")
		if values.Length() < 3 {
			return
		}

		// [0] 最新价
		quote.Price = parseF(strings.TrimSpace(values.Eq(0).Text()))

		// [1] 涨跌额
		quote.Change = parseF(strings.TrimSpace(values.Eq(1).Text()))

		// [2] 涨跌幅 (%)
		quote.ChangePct = parseF(strings.TrimSpace(values.Eq(2).Text()))

		// [3] 成交量 (带单位：万/亿)
		if values.Length() > 3 {
			quote.Volume = parseWithUnit(strings.TrimSpace(values.Eq(3).Text()))
		}

		// [4] 成交额 (带单位：万/亿)
		if values.Length() > 4 {
			quote.Amount = parseWithUnit(strings.TrimSpace(values.Eq(4).Text()))
		}

		// [5] 总市值 (带单位：万/亿)
		if values.Length() > 5 {
			quote.TotalCap = parseWithUnitPtr(strings.TrimSpace(values.Eq(5).Text()))
		}

		// [6] 流通市值 (带单位：万/亿)
		if values.Length() > 6 {
			quote.FloatCap = parseWithUnitPtr(strings.TrimSpace(values.Eq(6).Text()))
		}

		// [7] 总股本 (带单位：万/亿)
		if values.Length() > 7 {
			quote.TotalShares = parseWithUnitPtr(strings.TrimSpace(values.Eq(7).Text()))
		}

		// [8] 流通股本 (带单位：万/亿)
		if values.Length() > 8 {
			quote.FloatShares = parseWithUnitPtr(strings.TrimSpace(values.Eq(8).Text()))
		}

		// [9] 5 日涨跌幅 (%)
		if values.Length() > 9 {
			quote.Change5d = parsePercentPtr(strings.TrimSpace(values.Eq(9).Text()))
		}

		// [10] 10 日涨跌幅 (%)
		if values.Length() > 10 {
			quote.Change10d = parsePercentPtr(strings.TrimSpace(values.Eq(10).Text()))
		}

		// [11] 20 日涨跌幅 (%)
		if values.Length() > 11 {
			quote.Change20d = parsePercentPtr(strings.TrimSpace(values.Eq(11).Text()))
		}

		// [12] 60 日涨跌幅 (%)
		if values.Length() > 12 {
			quote.Change60d = parsePercentPtr(strings.TrimSpace(values.Eq(12).Text()))
		}

		// [13] 120 日涨跌幅 (%)
		if values.Length() > 13 {
			quote.Change120d = parsePercentPtr(strings.TrimSpace(values.Eq(13).Text()))
		}

		// [14] 250 日涨跌幅 (%)
		if values.Length() > 14 {
			quote.Change250d = parsePercentPtr(strings.TrimSpace(values.Eq(14).Text()))
		}

		// [15] 年初至今涨跌幅 (%)
		if values.Length() > 15 {
			quote.ChangeYtd = parsePercentPtr(strings.TrimSpace(values.Eq(15).Text()))
		}

		if quote.StockCode != "" {
			quotes = append(quotes, quote)
		}
	})

	return quotes, nil
}

// parseWithUnit 解析带单位的数值（支持%后缀）
func parseWithUnit(s string) float64 {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0
	}

	// 移除百分号
	s = strings.ReplaceAll(s, "%", "")

	multiplier := 1.0
	if strings.Contains(s, "万亿") {
		multiplier = 100000000000000000
		s = strings.ReplaceAll(s, "万亿", "")
	} else if strings.Contains(s, "亿") {
		multiplier = 100000000
		s = strings.ReplaceAll(s, "亿", "")
	} else if strings.Contains(s, "万") {
		multiplier = 10000
		s = strings.ReplaceAll(s, "万", "")
	} else if strings.Contains(s, "千") {
		multiplier = 1000
		s = strings.ReplaceAll(s, "千", "")
	}

	// 移除逗号
	s = strings.ReplaceAll(s, ",", "")

	return parseF(s) * multiplier
}

// parsePercent 解析百分比数值（移除%符号）
func parsePercent(s string) float64 {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0
	}
	s = strings.ReplaceAll(s, "%", "")
	return parseF(s)
}

// ptrFloat64 创建 float64 指针（仅当值有效时）
func ptrFloat64(val float64, hasData bool) *float64 {
	if !hasData || val == 0 {
		return nil
	}
	return &val
}

// parsePercentPtr 解析百分比数值并返回指针
func parsePercentPtr(s string) *float64 {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	s = strings.ReplaceAll(s, "%", "")
	val := parseF(s)
	if val == 0 {
		return nil
	}
	return &val
}

// parseWithUnitPtr 解析带单位的数值并返回指针
func parseWithUnitPtr(s string) *float64 {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	val := parseWithUnit(s)
	if val == 0 {
		return nil
	}
	return &val
}

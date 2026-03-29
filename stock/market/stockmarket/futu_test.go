package stockmarket

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

// TestAllFutu 测试富途牛牛 A 股行情数据爬取
// 使用前需要启动带调试端口的 Chrome:
//
//	chrome --remote-debugging-port=9222 --headless --no-sandbox --disable-gpu
func TestAllFutu(t *testing.T) {
	m := NewMarket(
		WithRemoteChrome("http://192.168.31.100:9222"),
		WithDebug(true),
	)

	count := 0
	startTime := time.Now()

	err := m.AllFutu(context.Background(), func(quote []*CurrentQuote) {
		count += len(quote)
		for _, q := range quote {
			fmt.Printf("[%d] %s - %s: %.2f (%.2f%%)\n",
				count,
				q.StockCode,
				q.ShortName,
				q.Price,
				q.ChangePct,
			)
		}
	})

	elapsed := time.Since(startTime)
	fmt.Printf("\n完成！共爬取 %d 条数据，耗时：%v\n", count, elapsed)

	if err != nil {
		t.Errorf("AllFutu 执行失败：%v", err)
	}
}

// TestExampleAllFutu 测试使用示例
func TestExampleAllFutu(t *testing.T) {
	// 创建 Market 实例，配置 Chrome 远程调试地址
	m := NewMarket(
		WithRemoteChrome("http://192.168.31.100:9222"), // Chrome 调试地址
	)

	// 爬取数据
	err := m.AllFutu(context.Background(), func(quote []*CurrentQuote) {
		// 回调处理每条数据
		data, _ := json.Marshal(quote)
		fmt.Printf("%d quotes: %s\n", len(quote), string(data))
	})

	if err != nil {
		fmt.Println("错误:", err)
	}
}

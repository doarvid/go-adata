## 目标

* 移除所有函数/方法签名中的 `wait time.Duration` 参数，将等待逻辑上移到各实体结构体的配置中统一管理。

* 统一支持 `WithWait(time.Duration)`、`WithRetries(int)`（必要时 `WithJitter`），由实例在内部循环与重试时使用，不再从调用方传入 wait。

## 设计规范

* 结构体字段：`Wait time.Duration`（默认 50ms）、`Retries int`（默认 2）

* 选项：`WithWait(d time.Duration)`、`WithRetries(n int)`、已有的 `WithTimeout/WithProxy/WithUserAgent/WithClient` 保留

* 方法签名统一：移除最后的 `wait time.Duration` 参数；内部在发起请求前 `time.Sleep(m.Wait)`，在重试之间也使用 `m.Wait`

* 测试：通过 `NewXxx(...WithWait(...))` 注入等待时间；不再在测试中传入 `wait`

## 涉及目录与文件（按模块）

* stock/market/stockmarket（已有 Market）

  * 移除 `GetDaily/GetMinute/GetBar/GetFive/ListCurrent` 及各数据源方法中的 `wait` 参数

  * 在 `Market` 增加 `Wait` 字段与 `WithWait` 选项；所有内部循环使用 `m.Wait`

  * 更新 \[market.go] 聚合入口、\[baidu.go]/\[east.go]/\[sina.go]/\[qq.go] 方法签名与实现

  * 更新 \[market\_test.go] 使用结构体选项，不再传 `wait`

* stock/market/indexmarket（已是结构体方法）

  * 在 `IndexMarketConfig` 增加 `Wait time.Duration` 与 `WithWait` 选项

  * 移除 `GetDaily/GetMinute/GetCurrent` 的 `wait` 参数，使用 `im.cfg.Wait`

  * 更新 \[market.go] 聚合入口（如有跨调用）与测试

* stock/market/capitalflow

  * 为 `Market`（或现有 Client）增加 `Wait` 字段与 `WithWait` 选项

  * 移除 `GetMin/GetDaily` 及底层 `GetStockCapitalFlow*` 相关方法的 `wait` 参数（或保留底层函数，但由 Market 在调用前 `Sleep`）

  * 更新 \[market.go]/\[baidu.go]/\[east.go] 实现与测试

* stock/market/conceptmarket

  * 引入/增强 `ConceptMarket` 结构体：增加 `Wait` 字段与选项

  * 移除 `GetConceptDaily/Minute/Current` 的 `wait` 参数，方法改为结构体方法使用 `m.Wait`

  * 更新 \[east.go]/\[ths.go] 与测试

* stock/market/conceptflow

  * 引入 `ConceptFlowClient` 配置 `Wait`；移除 `ListConceptCapitalFlowEast` 的 `wait` 参数

  * 更新 \[east.go] 与测试

* stock/info/stockindex

  * 引入/增强 `StockIndexClient` 配置 `Wait`；移除 `AllIndexCodeEast`、`IndexConstituentBaidu` 的 `wait` 参数

  * 更新 \[stockindex.go] 与（若有）测试

* stock/info/stockcode

  * 引入 `StockCodeClient` 配置 `Wait`；将 `marketRankEast/Baidu/Sina/newSubEast/AllCode` 的 `wait` 参数移除并改为方法

  * 更新 \[stockcode.go] 与测试

* sentiment/\*（alist、hot、lifting、mine、margin、northflow）

  * 所有 Client 增加 `Wait` 配置；移除方法中的 `wait` 参数，如 `Daily/Details/Popular/Stocks/Concepts/History/Minute/Current` 等

  * 更新各文件与（若存在）测试

* stock/stock.go 对外聚合

  * 删除所有对外方法的 `wait` 参数：`Daily/Minute/Bar/Five/Current/IndexDaily/IndexMinute/...`

  * 内部按需创建对应 Market/Client，并通过 `WithWait` 注入等待参数以维持行为一致

## 兼容策略

* 迁移期可提供薄包装（Deprecated）：保留旧签名但忽略传入的 `wait`，转调新方法；一个版本后删除

* 测试中对外网不稳定接口继续使用 `t.Skip` 保持稳定

## 实施步骤

1. 在每个模块引入/增强选项：`WithWait/WithRetries`，结构体加入 `Wait/Retries` 字段
2. 修改方法签名，移除 `wait` 参数；在内部循环与请求前使用结构体的 `Wait`
3. 更新聚合入口（market.go、stock/stock.go 等）与测试调用
4. 全仓 grep 检查：不再存在 `wait time.Duration` 作为方法参数
5. 构建与测试：`go build ./...`、`go test ./... -run Test`，逐项修复

## 风险与回退

* 若第三方接口限流导致失败，适当提升默认 `Wait` 或在测试用例 `Skip`

* 若需要快速回退，保留的薄包装可短期维持旧签名调用

## 交付

* 更新统一客户端规范文档，加入 `WithWait/WithRetries` 示例

* 提供示例：

```go
m := stockmarket.NewMarket(WithWait(100*time.Millisecond))
rows, _ := m.GetDaily("002926", "2025-01-01", "2025-01-31", stockmarket.KTypeDay, stockmarket.AdjustTypePre)

idx := indexmarket.NewIndexMarket(indexmarket.WithWait(100*time.Millisecond))
d, _ := idx.GetDailyEast(ctx, "000001", "2024-01-01", 1)
```


## 重构目标

* 全项目移除函数/方法签名中的 wait time.Duration 参数，统一由实体结构体的配置管理等待与重试（WithWait/WithRetries）。

* 将所有包级函数迁移为结构体方法，实例级共享 resty.Client；彻底删除 getHTTPClient。

* 保持命名与调用一致性：Daily/Minute/Current/List/Bar/Five 等。

## 审核与定位

* 通过全仓检索确认存在 wait 的位置（示例）：

  * stock/market/stockmarket：GetDaily/Minute/Bar/Five/ListCurrent 及各数据源方法（已部分改造，下一步移除签名中的 wait）

  * stock/market/indexmarket：GetDailyEast/MinuteEast/CurrentEast/Ths 等（IndexMarket 方法仍含 wait）

  * stock/market/conceptflow：east.go 的 ListConceptCapitalFlowEast(daysType, wait)

  * stock/market/conceptmarket：east.go、ths.go 多个 GetConcept\*(..., wait)

  * stock/market/capitalflow：GetStockCapitalFlow\*(..., wait) 及 Market/Client 封装

  * stock/info/stockindex：AllIndexCodeEast、IndexConstituentBaidu(..., wait)

  * stock/info/stockcode：AllCode、marketRankEast/Baidu/Sina/newSubEast(..., wait)

  * sentiment/\*：alist/hot/lifting/mine/margin/northflow 客户端方法普遍含 wait

  * stock/stock.go 聚合对外 API 也传递 wait

## 统一设计

* 结构体字段：Wait time.Duration（默认 50ms）、Retries int（默认 2）。

* 选项：WithWait(d)、WithRetries(n)、以及已存在的 WithTimeout/WithProxy/WithUserAgent/WithClient。

* 内部使用：所有发请求前与重试间隔使用 m.Wait；不再从调用方传入 wait。

## 分模块实施

### 概念资金流（conceptflow）

* 引入 ConceptFlowClient（字段：client、Wait、Retries、选项）。

* 将 ListConceptCapitalFlowEast(daysType, wait) 改为 (c \*ConceptFlowClient).ListEast(ctx, daysType)，内部用 c.Wait；删除 http\_client.go 并用实例 client。

* 并发页抓取保持不变，仅把 wait 控制改为 c.Wait；保持 Normalize 与解析逻辑。

### 概念市场（conceptmarket）

* 引入/增强 ConceptMarket 结构体（client、Wait、Retries）。

* east.go/ths.go：GetConceptDaily/Minute/Current 均改为结构体方法且移除 wait；内部用 m.Wait。

* 删除包级 http\_client.go。

### 资金流（capitalflow）

* 建立 Market 结构体（或增强现有 Client）：client、Wait、Retries、选项。

* 将 GetStockCapitalFlowMinBaidu/East、GetStockCapitalFlowBaidu/East 改为结构体方法 MinutesBaidu/MinutesEast、DailyBaidu/DailyEast；移除 wait。

### 指数行情（indexmarket）

* 在 IndexMarketConfig 增加 Wait、WithWait；

* 移除 GetDaily/Minute/Current 的 wait 参数，内部使用 cfg.Wait；

* market.go 聚合入口若涉及，同步更新。

### 股票信息（stockindex、stockcode）

* 引入/增强 StockIndexClient、StockCodeClient（client、Wait、Retries）。

* 将 AllIndexCodeEast、IndexConstituentBaidu、marketRankEast/Baidu/Sina、新股 newSubEast 统一改为结构体方法并移除 wait。

### 情绪与其他（sentiment/\*）

* 为 alist/hot/lifting/mine/margin/northflow 客户端增加 Wait/WithWait；

* 所有方法移除 wait 参数，内部使用客户端 Wait。

### 对外聚合（stock/stock.go）

* 删除对外 API 的 wait 形参：Daily/Minute/Bar/Five/Current/Index\*、Concept\*、CapitalFlow\* 等；

* 内部通过 NewXxx(WithWait(...)) 注入等待时间以维持行为一致；必要时提供 Deprecated 包装保留一版过渡。

## 验证步骤

* 执行模块改造后，逐包运行 go build 与单包测试；

* 对外网不稳定的测试保留 t.Skip 保障稳定；

* 全仓 go build ./... 与 go test ./... -run Test。

## 交付与兼容

* 文档补充：统一客户端规范增加 WithWait 示例；

* 过渡期保留旧签名的薄包装（Deprecated），调用新方法忽略传入 wait；最终删除。

## 示例（概念资金流改造后）

* 用法：

```go
c := conceptflow.New(WithWait(100*time.Millisecond))
rows, _ := c.ListEast(ctx, 1)
```


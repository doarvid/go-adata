## 总体目标
- 全仓统一为“实体结构体 + 上下文 + 实例级 resty.Client”模式；彻底移除各包的 getHTTPClient。
- 覆盖 stockmarket、capitalflow、conceptmarket、conceptflow、stockindex、stockcode 等目录，避免遗漏。

## 目录与文件清单
- stock/market/stockmarket（数据源→Market 方法）
  - east.go：GetMarketDailyEastCtx、GetMarketMinuteEastCtx → Market.GetDailyEast/MinuteEast
  - baidu.go：GetMarketDailyBaiduCtx、GetMarketMinuteBaiduCtx、GetMarketBarBaiduCtx、GetMarketFiveBaiduCtx → Market.GetDailyBaidu/MinuteBaidu/BarBaidu/FiveBaidu
  - sina.go：ListMarketCurrentSinaCtx → Market.ListCurrentSina
  - qq.go：已完成迁移
  - http_client.go：删除
  - market.go：聚合入口统一改为 m.* 方法
  - market_test.go：统一改为 Market 方法
- stock/market/capitalflow（数据源→Market 方法）
  - baidu.go、east.go：GetStockCapitalFlowMinBaidu/East、GetStockCapitalFlowBaidu/East → Market.MinutesBaidu/MinutesEast、Market.DailyBaidu/DailyEast
  - market.go：改为使用实例方法；注入 client 字段与选项
  - http_client.go：删除
  - baidu_test.go：改为 Market 方法
- stock/market/conceptmarket（数据源→ConceptMarket 方法）
  - east.go、ths.go：GetConceptDailyEast/MinuteEast/CurrentEast、GetConceptDailyThs/MinuteThs/CurrentThs → ConceptMarket.DailyEast/MinuteEast/CurrentEast、ConceptMarket.DailyThs/MinuteThs/CurrentThs
  - 新增 ConceptMarket 结构体及选项，内部共享 client
  - http_client.go：删除
  - 测试：改为结构体方法调用
- stock/market/conceptflow（数据源→ConceptFlowClient 方法）
  - east.go：ListConceptCapitalFlowEast → ConceptFlowClient.ListEast
  - 新增 ConceptFlowClient 结构体及选项；删除 http_client.go
  - 测试更新
- stock/market/indexmarket（已是结构体方法）
  - 已符合规范；补充选项一致性检查
- stock/info/concept（已是结构体方法）
  - 已符合规范
- stock/info/stockindex（数据源→StockIndexClient 方法）
  - stockindex.go：AllIndexCodeEast、IndexConstituentBaidu → StockIndexClient.AllCodesEast/ConstituentBaidu
  - 新增结构体选项；删除 http_client.go
  - 测试（若需要）补充
- stock/info/stockcode（数据源→StockCodeClient 方法）
  - stockcode.go：marketRankEast/Baidu/Sina、newSubEast → StockCodeClient.* 方法
  - 新增结构体选项；删除 http_client.go
  - stockcode_test.go：保持或调整为结构体方法

## 技术一致性
- 选项统一：WithTimeout/WithProxy/WithUserAgent/WithClient（必要时 WithHeaders）
- 命名统一：Daily/Minute/Current/List/Bar/Five
- 上下文统一：所有方法首参为 ctx context.Context

## 实施步骤
1. 在各包引入实体结构体与选项，初始化 resty.Client（实例级共享）
2. 将包级 Ctx 函数逐个迁移为结构体方法，内部复用 m.client；旧函数暂时保留薄包装（Deprecated）方便过渡，最终移除
3. 聚合入口与测试同步更新为结构体方法调用；对外网不稳定接口保留 t.Skip
4. 全仓清理：删除所有 http_client.go、移除 getHTTPClient 引用；grep 检查无残留
5. 构建与测试：go build ./...；go test ./... -run Test；逐项修复

## 验证与回退
- 模块化提交，出现外网失败则在测试用例中 Skip；如需回退，仅针对单目录的迁移提交回滚即可

## 交付
- 完成后提供“统一客户端规范”更新与使用示例，覆盖上述各目录的典型调用
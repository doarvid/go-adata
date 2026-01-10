## 目标
- 以业务实体为中心（如 AList）重构 HTTP 访问层：每个实体可独立配置客户端与数据源。
- 采用函数式选项（WithProxy 等）进行配置注入，避免到处显式构造客户端。
- 所有接口统一接收 ctx context.Context，支持取消与超时。

## 实体建模
### 核心结构
- 以实体为单位定义客户端封装：
  - 以 sentiment.AList 为例：
    - type AList struct { client *resty.Client; cfg AListConfig }
    - type AListConfig struct { Timeout time.Duration; Proxy string; UserAgent string; Source string; Headers map[string]string }
    - func NewAList(opts ...AListOption) *AList
    - type AListOption func(*AListConfig)
- 其他实体（如 IndexMarket、StockMarket、Concept、CapitalFlow 等）沿用同样模式，便于扩展和统一管理。

### 函数式选项（示例）
- WithProxy(url string)
- WithTimeout(d time.Duration)
- WithUserAgent(ua string)
- WithSource(src string)（如 "east"、"ths"、"baidu"、"sina"、"qq"）
- WithHeaders(h map[string]string)
- WithClient(c *resty.Client)（直接注入已配置的客户端）

### 方法与上下文
- 将原函数改为实体方法，首参 ctx：
  - AList.ListDaily(ctx context.Context, reportDate string, wait time.Duration) ([]AListDaily, error)
  - AList.GetInfo(ctx context.Context, stockCode string, reportDate string, wait time.Duration) ([]AListInfo, error)
- 发起请求时绑定上下文：client.R().SetContext(ctx).SetHeaders(cfg.Headers).Post/Get。
- 依据 cfg.Source 选择不同的数据源 URL 与解析逻辑（例如 east/ths 的差异）。

## 兼容适配
- 保留旧的包级函数作为过渡层：
  - ListAListDaily(reportDate string, wait time.Duration) 调用 NewAList() 后用 context.Background() 代理到 AList.ListDaily。
  - GetAListInfo(...同理)。
- 逐步迁移调用方到实体方法与 ctx 版本；完成后移除过渡层。

## 改动范围与示例
- 删除公共库：[common/http/http.go](file:///Users/duxiaoliang/project/stock_quant/go-adata/common/http/http.go)。统一直接依赖 resty。
- 在 sentiment 包新增 [client.go]（文件名示意）：定义 AListConfig/AListOption/NewAList/WithXXX 与内部构造 resty.Client。
- 修改 [alist.go](file:///Users/duxiaoliang/project/stock_quant/go-adata/sentiment/alist.go)：
  - 将现有函数改为 AList 的方法，切换到 getClient(ctx) + SetContext。
  - 根据 cfg.Source 决定 URL（east 默认；未来可扩展 ths 等）。
- 其他包含 httpc.NewClient 的包按同样模式引入对应实体（或聚合实体）与 WithXXX 选项，替换所有 httpc.NewClient 的直接使用。

## 数据源多样性
- 在每个实体的 cfg.Source 中支持多数据源：east、ths、baidu、sina、qq。
- 在方法实现中通过 switch cfg.Source 选择不同 URL 和解析函数（保证结构清晰，避免散落 if/else）。

## 测试与验证
- 构造 AList 实例，使用 WithTimeout/WithProxy 等验证选项生效。
- 使用 httptest 模拟慢响应与取消，验证 ctx 生效。
- 针对不同 Source 的解析测试，确保结果一致。

## 验收标准
- 所有直接 httpc.NewClient 的调用被移除，替换为实体层封装与可配置客户端。
- 所有对外接口具备 ctx 参数并在请求上绑定。
- 实体层的 WithXXX 选项可独立配置，支持多数据源；测试通过。
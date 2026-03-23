package finance

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-resty/resty/v2"

	"github.com/doarvid/go-adata/common/utils"
)

// CoreIndex 核心指标结构体
type CoreIndex struct {
	StockCode                 string  `json:"stock_code"`
	ShortName                 string  `json:"short_name"`
	ReportDate                string  `json:"report_date"`
	ReportType                string  `json:"report_type"`
	NoticeDate                string  `json:"notice_date"`
	BasicEPS                  float64 `json:"basic_eps"`
	DilutedEPS                float64 `json:"diluted_eps"`
	NonGaapEPS                float64 `json:"non_gaap_eps"`
	NetAssetPS                float64 `json:"net_asset_ps"`
	CapReservePS              float64 `json:"cap_reserve_ps"`
	UndistProfitPS            float64 `json:"undist_profit_ps"`
	OperCFPS                  float64 `json:"oper_cf_ps"`
	TotalRev                  float64 `json:"total_rev"`
	GrossProfit               float64 `json:"gross_profit"`
	NetProfitAttrSH           float64 `json:"net_profit_attr_sh"`
	NonGaapNetProfit          float64 `json:"non_gaap_net_profit"`
	TotalRevYoyGr             float64 `json:"total_rev_yoy_gr"`
	NetProfitYoyGr            float64 `json:"net_profit_yoy_gr"`
	NonGaapNetProfitYoyGr     float64 `json:"non_gaap_net_profit_yoy_gr"`
	TotalRevQoqGr             float64 `json:"total_rev_qoq_gr"`
	NetProfitQoqGr            float64 `json:"net_profit_qoq_gr"`
	NonGaapNetProfitQoqGr     float64 `json:"non_gaap_net_profit_qoq_gr"`
	RoeWtd                    float64 `json:"roe_wtd"`
	RoeNonGaapWtd             float64 `json:"roe_non_gaap_wtd"`
	RoaWtd                    float64 `json:"roa_wtd"`
	GrossMargin               float64 `json:"gross_margin"`
	NetMargin                 float64 `json:"net_margin"`
	AdvReceiptsToRev          float64 `json:"adv_receipts_to_rev"`
	NetCfSalesToRev           float64 `json:"net_cf_sales_to_rev"`
	OperCfToRev               float64 `json:"oper_cf_to_rev"`
	EffTaxRate                float64 `json:"eff_tax_rate"`
	CurrRatio                 float64 `json:"curr_ratio"`
	QuickRatio                float64 `json:"quick_ratio"`
	CashFlowRatio             float64 `json:"cash_flow_ratio"`
	AssetLiabRatio            float64 `json:"asset_liab_ratio"`
	EquityMultiplier          float64 `json:"equity_multiplier"`
	EquityRatio               float64 `json:"equity_ratio"`
	TotalAssetTurnDays        float64 `json:"total_asset_turn_days"`
	InvTurnDays               float64 `json:"inv_turn_days"`
	AcctRecvTurnDays          float64 `json:"acct_recv_turn_days"`
	TotalAssetTurnRate        float64 `json:"total_asset_turn_rate"`
	InvTurnRate               float64 `json:"inv_turn_rate"`
	AcctRecvTurnRate          float64 `json:"acct_recv_turn_rate"`
}

type Client struct {
	client *resty.Client
}

func New() *Client {
	return &Client{
		client: resty.New(),
	}
}

func (c *Client) GetCoreIndex(ctx context.Context, stockCode string) ([]CoreIndex, error) {
	return c.getCoreIndexEast(ctx, stockCode)
}

func (c *Client) getCoreIndexEast(ctx context.Context, stockCode string) ([]CoreIndex, error) {
	reportTypes := []string{"年报", "中报", "三季报", "一季报"}
	stockCode = utils.CompileExchangeByStockCode(stockCode)

	var data []map[string]any
	for _, reportType := range reportTypes {
		url := fmt.Sprintf("https://datacenter.eastmoney.com/securities/api/data/get?type=RPT_F10_FINANCE_MAINFINADATA&sty=APP_F10_MAINFINADATA&quoteColumns=&filter=(SECUCODE=\"%s\")(REPORT_TYPE=\"%s\")&p=1&ps=100&sr=-1&st=REPORT_DATE&source=HSF10&client=PC&v=03890754131799983",
			stockCode, reportType)

		resp, err := c.client.R().
			SetContext(ctx).
			Get(url)
		if err != nil {
			continue
		}

		if resp.StatusCode() != 200 {
			continue
		}

		var result struct {
			Code int `json:"code"`
			Result struct {
				Data []map[string]any `json:"data"`
			} `json:"result"`
		}

		if err := json.Unmarshal(resp.Body(), &result); err != nil {
			continue
		}

		if result.Code == 0 && result.Result.Data != nil {
			data = append(data, result.Result.Data...)
		}
	}


	var coreIndices []CoreIndex
	for _, item := range data {
		ci := CoreIndex{}

		// 处理字符串字段
		if v, ok := item["SECURITY_CODE"]; ok {
			ci.StockCode = toString(v)
		}
		if v, ok := item["SECURITY_NAME_ABBR"]; ok {
			ci.ShortName = toString(v)
		}
		if v, ok := item["REPORT_DATE"]; ok {
			ci.ReportDate = formatDate(toString(v))
		}
		if v, ok := item["REPORT_TYPE"]; ok {
			ci.ReportType = toString(v)
		}
		if v, ok := item["NOTICE_DATE"]; ok {
			ci.NoticeDate = formatDate(toString(v))
		}

		// 处理数值字段
		ci.BasicEPS = parseFloat(item["EPSJB"])
		ci.DilutedEPS = parseFloat(item["EPSKCJB"])
		ci.NonGaapEPS = parseFloat(item["EPSXS"])
		ci.NetAssetPS = parseFloat(item["BPS"])
		ci.CapReservePS = parseFloat(item["MGZBGJ"])
		ci.UndistProfitPS = parseFloat(item["MGWFPLR"])
		ci.OperCFPS = parseFloat(item["MGJYXJJE"])
		ci.TotalRev = parseFloat(item["TOTALOPERATEREVE"])
		ci.GrossProfit = parseFloat(item["MLR"])
		ci.NetProfitAttrSH = parseFloat(item["PARENTNETPROFIT"])
		ci.NonGaapNetProfit = parseFloat(item["KCFJCXSYJLR"])
		ci.TotalRevYoyGr = parseFloat(item["TOTALOPERATEREVETZ"])
		ci.NetProfitYoyGr = parseFloat(item["PARENTNETPROFITTZ"])
		ci.NonGaapNetProfitYoyGr = parseFloat(item["KCFJCXSYJLRTZ"])
		ci.TotalRevQoqGr = parseFloat(item["YYZSRGDHBZC"])
		ci.NetProfitQoqGr = parseFloat(item["NETPROFITRPHBZC"])
		ci.NonGaapNetProfitQoqGr = parseFloat(item["KFJLRGDHBZC"])
		ci.RoeWtd = parseFloat(item["ROEJQ"])
		ci.RoeNonGaapWtd = parseFloat(item["ROEKCJQ"])
		ci.RoaWtd = parseFloat(item["ZZCJLL"])
		ci.GrossMargin = parseFloat(item["XSMLL"])
		ci.NetMargin = parseFloat(item["XSJLL"])
		ci.AdvReceiptsToRev = parseFloat(item["YSZKYYSR"])
		ci.NetCfSalesToRev = parseFloat(item["XSJXLYYSR"])
		ci.OperCfToRev = parseFloat(item["JYXJLYYSR"])
		ci.EffTaxRate = parseFloat(item["TAXRATE"])
		ci.CurrRatio = parseFloat(item["LD"])
		ci.QuickRatio = parseFloat(item["SD"])
		ci.CashFlowRatio = parseFloat(item["XJLLB"])
		ci.AssetLiabRatio = parseFloat(item["ZCFZL"])
		ci.EquityMultiplier = parseFloat(item["QYCS"])
		ci.EquityRatio = parseFloat(item["CQBL"])
		ci.TotalAssetTurnDays = parseFloat(item["ZZCZZTS"])
		ci.InvTurnDays = parseFloat(item["CHZZTS"])
		ci.AcctRecvTurnDays = parseFloat(item["YSZKZZTS"])
		ci.TotalAssetTurnRate = parseFloat(item["TOAZZL"])
		ci.InvTurnRate = parseFloat(item["CHZZL"])
		ci.AcctRecvTurnRate = parseFloat(item["YSZKZZL"])

		coreIndices = append(coreIndices, ci)
	}

	// 按报告日期降序排序
	for i := range coreIndices {
		for j := i + 1; j < len(coreIndices); j++ {
			if coreIndices[i].ReportDate < coreIndices[j].ReportDate {
				coreIndices[i], coreIndices[j] = coreIndices[j], coreIndices[i]
			}
		}
	}

	return coreIndices, nil
}

func toString(v any) string {
	return strings.TrimSpace(fmt.Sprintf("%v", v))
}

func parseFloat(v any) float64 {
	s := toString(v)
	s = strings.TrimSpace(strings.ReplaceAll(s, "%", ""))
	if s == "" || s == "-" || s == "--" {
		return 0.0
	}
	f, _ := parseFloat64(s)
	return f
}

func parseFloat64(s string) (float64, error) {
	return parseFloat64WithDefault(s, 0.0)
}

func parseFloat64WithDefault(s string, def float64) (float64, error) {
	s = strings.TrimSpace(s)
	if s == "" || s == "-" || s == "--" {
		return def, nil
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return def, err
	}
	return f, nil
}

func formatDate(dateStr string) string {
	dateStr = strings.TrimSpace(dateStr)
	if len(dateStr) == 0 {
		return ""
	}

	// 处理可能的日期格式，如 "2024-09-03 00:00:00" 转换为 "2024-09-03"
	if idx := strings.Index(dateStr, " "); idx != -1 {
		dateStr = dateStr[:idx]
	}

	return dateStr
}

package capitalflow

import (
	"context"
)

type Client struct{}

func NewClient() *Client { return &Client{} }

func (c *Client) MinutesBaidu(ctx context.Context, stockCode string) ([]FlowMin, error) {
	return GetStockCapitalFlowMinBaidu(stockCode)
}

func (c *Client) MinutesEast(ctx context.Context, stockCode string) ([]FlowMin, error) {
	return GetStockCapitalFlowMinEast(stockCode)
}

func (c *Client) DailyBaidu(ctx context.Context, stockCode string, startDate, endDate string) ([]FlowDaily, error) {
	return GetStockCapitalFlowBaidu(stockCode, startDate, endDate)
}

func (c *Client) DailyEast(ctx context.Context, stockCode string, startDate, endDate string) ([]FlowDaily, error) {
	return GetStockCapitalFlowEast(stockCode, startDate, endDate)
}

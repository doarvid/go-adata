package capitalflow

import (
	"context"
	"time"
)

type Client struct{}

func NewClient() *Client { return &Client{} }

func (c *Client) MinutesBaidu(ctx context.Context, stockCode string, wait time.Duration) ([]FlowMin, error) {
	return GetStockCapitalFlowMinBaidu(stockCode, wait)
}

func (c *Client) MinutesEast(ctx context.Context, stockCode string, wait time.Duration) ([]FlowMin, error) {
	return GetStockCapitalFlowMinEast(stockCode, wait)
}

func (c *Client) DailyBaidu(ctx context.Context, stockCode string, startDate, endDate string, wait time.Duration) ([]FlowDaily, error) {
	return GetStockCapitalFlowBaidu(stockCode, startDate, endDate, wait)
}

func (c *Client) DailyEast(ctx context.Context, stockCode string, startDate, endDate string, wait time.Duration) ([]FlowDaily, error) {
	return GetStockCapitalFlowEast(stockCode, startDate, endDate, wait)
}

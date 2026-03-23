package finance

import (
	"context"
)

type Finance struct {
	client *Client
}

func NewFinance() *Finance {
	return &Finance{
		client: New(),
	}
}

func (f *Finance) GetCoreIndex(ctx context.Context, stockCode string) ([]CoreIndex, error) {
	return f.client.GetCoreIndex(ctx, stockCode)
}

var defaultFinance = NewFinance()

func GetCoreIndex(ctx context.Context, stockCode string) ([]CoreIndex, error) {
	return defaultFinance.GetCoreIndex(ctx, stockCode)
}

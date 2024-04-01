package binance

import (
	"context"
	"errors"
	"fmt"
	"github.com/aiviaio/go-binance/v2"
)

var ErrPriceNotFound = errors.New("price not found")

type Client interface {
	GetSymbols(ctx context.Context, count int) ([]string, error)
	GetLastPriceBySymbol(ctx context.Context, symbol string) (string, error)
}
type clientImpl struct {
	binanceClient *binance.Client
}

func NewClient() Client {
	return &clientImpl{
		binanceClient: binance.NewClient("", ""),
	}
}

func (c clientImpl) GetSymbols(ctx context.Context, count int) ([]string, error) {
	responseExchangeInfo, err := c.binanceClient.NewExchangeInfoService().Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("can`t exchange info: %w", err)
	}

	symbols := make([]string, count)
	for i, val := range responseExchangeInfo.Symbols[:count] {
		symbols[i] = val.Symbol
	}

	return symbols, nil
}

func (c clientImpl) GetLastPriceBySymbol(ctx context.Context, symbol string) (string, error) {
	resp, err := c.binanceClient.NewListSymbolTickerService().Symbol(symbol).Do(ctx)
	if err != nil {
		return "", fmt.Errorf("can`t exchange info: %w", err)
	}

	if len(resp) == 0 {
		return "", ErrPriceNotFound
	}

	return resp[0].LastPrice, nil
}

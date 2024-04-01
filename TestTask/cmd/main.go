package main

import (
	"TestTask/pkg/binance"
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"log"
)

const (
	symbolsCount = 5
)

func main() {
	ctx := context.Background()
	client := binance.NewClient()
	symbols, err := client.GetSymbols(ctx, symbolsCount)
	if err != nil {
		log.Fatal(err.Error())
	}
	if len(symbols) == 0 {
		log.Println("symbols not found")
		return
	}

	data := make(chan map[string]string)

	errG, ctxG := errgroup.WithContext(ctx)
	for _, val := range symbols {
		val := val
		errG.Go(func() error {
			price, err := client.GetLastPriceBySymbol(ctxG, val)
			if err != nil {
				return err
			}
			data <- map[string]string{
				val: price,
			}
			return nil
		})
	}

	errG.Go(func() error {
		i := symbolsCount
		for m := range data {
			for k, v := range m {
				fmt.Println(fmt.Sprintf("%s %s", k, v))
			}
			if i--; i == 0 {
				close(data)
			}
		}
		return nil
	})

	if err = errG.Wait(); err != nil {
		log.Fatal(err.Error())
	}
}

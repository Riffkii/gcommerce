package rpc

import (
	"context"
	"order/proto/compiled"
	"time"
)

func GetProducts(client compiled.ProductServiceClient, ids []int64) (*compiled.Products, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	event, err := client.GetProducts(ctx, &compiled.ProductIds{Ids: ids})
	if err != nil {
		return nil, err
	}
	return event, nil
}

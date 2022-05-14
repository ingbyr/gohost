package biz

import (
	"context"
	"fmt"

	"github.com/ingbyr/gohost/dal/query"
)

func Query(ctx context.Context) {
	tx := query.Product
	do := tx.WithContext(context.Background())
	products, err := do.Find()
	catchError("list all", err)
	for _, product := range products {
		fmt.Println(product)
	}
}

func catchError(detail string, err error) {
	if err != nil {
		fmt.Printf("%s: %v\n", detail, err)
	}
}

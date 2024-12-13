package main

import (
	"context"
	"log"

	"github.com/whosonfirst/go-whosonfirst-sql/app/index"
)

func main() {

	ctx := context.Background()
	err := index.Run(ctx)

	if err != nil {
		log.Fatalf("Failed to index, %v", err)
	}
}

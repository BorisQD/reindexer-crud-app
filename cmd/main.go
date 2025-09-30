package main

import (
	"context"
	"crud/internal/app"
	"log"
)

func main() {
	err := app.Start(context.Background())
	if err != nil {
		log.Fatal(err.Error())
	}
}

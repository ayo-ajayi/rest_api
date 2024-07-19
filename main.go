package main

import (
	"github.com/ayo-ajayi/rest_api_template/app"
	_ "github.com/lib/pq"
	"log"
)

func main() {
	a, err := app.NewApp()
	if err != nil {
		log.Fatalf("error initializing app: %v", err)
	}
	a.Start()
}

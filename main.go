package main

import (
	"log"

	"github.com/ayo-ajayi/rest_api_template/route"
	_ "github.com/lib/pq"
)

func main() {
	server := route.Router()
	log.Fatal(server.Run(":8000"))

}

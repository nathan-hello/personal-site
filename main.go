package main

import (
	"log"

	"github.com/nathan-hello/personal-site/src/db"
)

func main() {
	err := db.InitDb()
	if err != nil {
		log.Fatal(err)
	}
}

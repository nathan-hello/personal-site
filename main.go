package main

import (
	"log"

	"github.com/nathan-hello/personal-site/src/db"
	"github.com/nathan-hello/personal-site/src/utils"
)

func main() {
	err := utils.InitEnv(nil)
	if err != nil {
		log.Fatal(err)
	}
	err = db.InitDb()
	if err != nil {
		log.Fatal(err)
	}

}

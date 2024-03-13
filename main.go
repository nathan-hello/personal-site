package main

import "github.com/nathan-hello/personal-site/src/db"
import "github.com/nathan-hello/personal-site/src/utils"
import "log"

func main() {
	err := utils.InitEnv()
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.InitDb()

}

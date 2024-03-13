package db

import "github.com/nathan-hello/personal-site/src/utils"

var Conn *Queries

func InitDb() error {
	var d, err = sql.Open("sqlite", utils.Env().DB_URI)
	if err != nil {
		return err
	}
	Conn = New(d)

	return nil
}

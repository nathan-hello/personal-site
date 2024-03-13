package db

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/nathan-hello/personal-site/src/utils"
)

var conn *Queries

func InitDb() error {

	var d, err = sql.Open("sqlite3", utils.Env().DB_URI)
	if err != nil {
		return err
	}
	err = d.Ping()
	if err != nil {
		fmt.Print("ping")
		return err
	}
	conn = New(d)
	return nil
}

func Conn() *Queries {
	if conn == nil {
		panic("asdf")
	}
	return conn
}

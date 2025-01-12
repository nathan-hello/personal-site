package db

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

var conn *Queries

func InitDb() error {

	var d, err = sql.Open("sqlite3", "file://data.db")
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

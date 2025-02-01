package db

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

var Conn *Queries

func InitDb() (*Queries, error) {
	if Conn != nil {
		return Conn, nil
	}

	var d, err = sql.Open("sqlite3", "file:test.db?cache=shared&mode=memory")
	if err != nil {
		return nil, err
	}
	err = d.Ping()
	if err != nil {
		fmt.Print("ping")
		return nil, err
	}
	Conn = New(d)

	return Conn, nil
}

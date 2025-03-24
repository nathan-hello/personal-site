package db

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed input/schema.sql
var schema string

var Conn *Queries

func InitDb(input string) (*Queries, error) {
	if Conn != nil {
		return Conn, nil
	}

	d, err := sql.Open("sqlite3", input)
	if err != nil {
		return nil, err
	}

	if _, err := d.ExecContext(context.Background(), schema); err != nil {
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

func Db() *Queries {
        return Conn
}

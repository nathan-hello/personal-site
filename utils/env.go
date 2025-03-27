package utils

import (
	_ "embed"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var env struct {
	access_expiry_time_hours  int
	refresh_expiry_time_hours int
	jwt_secret                string
}

type parsed_env struct {
	ACCESS_EXPIRY_TIME  time.Duration
	REFRESH_EXPIRY_TIME time.Duration
	JWT_SECRET          string
	DATABASE_URI        string
}

var parsed = parsed_env{}

func testStruct() error {
	reflected := reflect.ValueOf(parsed)
	for i := range reflected.NumField() {
		good := reflected.Field(i).IsValid()
		if !good {
			return fmt.Errorf("there was a zero value int he dotenv struct: %#v", parsed)
		}
	}
	return nil
}

func ParseDotenv(dotenv string) error {
	if dotenv == "" {
		return fmt.Errorf(".env not loaded")
	}

	lines := strings.Split(dotenv, "\n")
	for i, v := range lines {
		asdf := strings.Split(v, "=")

		if len(asdf) != 2 {
            return nil
		}

		key := asdf[0]
		value := asdf[1]

		switch key {
		case "ACCESS_EXPIRY_TIME":
			num, err := strconv.Atoi(value)
			if err != nil {
				return fmt.Errorf("bad value for ACCESS_EXPIRY_TIME %s %w", v, err)
			}
			parsed.ACCESS_EXPIRY_TIME = time.Hour * time.Duration(num)
		case "REFRESH_EXPIRY_TIME":
			num, err := strconv.Atoi(value)
			if err != nil {
				return fmt.Errorf("bad value for REFRESH_EXPIRY_TIME %s %w", v, err)
			}
			parsed.REFRESH_EXPIRY_TIME = time.Hour * time.Duration(num)
		case "JWT_SECRET":
			parsed.JWT_SECRET = value
		case "DATABASE_URI":
			parsed.DATABASE_URI = value
		default:
			return fmt.Errorf("unknown key at line %d %s", i, v)
		}

	}
	err := testStruct()
	if err != nil {
		return err
	}
	return nil
}

func Env() *parsed_env {
	return &parsed
}

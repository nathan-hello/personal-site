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
	webhook_secret            string
}

type parsed_env struct {
	ACCESS_EXPIRY_TIME  time.Duration
	REFRESH_EXPIRY_TIME time.Duration
	JWT_SECRET          string
	DATABASE_URI        string
	ADMIN_PASS          string
	WEBHOOK_SECRET      string
	LOG_PATH            string
}

var parsed = parsed_env{}

func testStruct() error {
	reflected := reflect.ValueOf(parsed)
	for i := range reflected.NumField() {
		good := reflected.Field(i).IsValid()
		if !good {
			return fmt.Errorf("there was a zero value in the dotenv struct: %#v", parsed)
		}
	}
	return nil
}

func ParseDotenv(dotenv string) error {
	if dotenv == "" {
		return fmt.Errorf(".env not loaded")
	}

	lines := strings.Split(dotenv, "\n")
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid line %d: %s", i+1, line)
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "ACCESS_EXPIRY_TIME":
			num, err := strconv.Atoi(value)
			if err != nil {
				return fmt.Errorf("bad value for ACCESS_EXPIRY_TIME on line %d: %w", i+1, err)
			}
			parsed.ACCESS_EXPIRY_TIME = time.Hour * time.Duration(num)
		case "REFRESH_EXPIRY_TIME":
			num, err := strconv.Atoi(value)
			if err != nil {
				return fmt.Errorf("bad value for REFRESH_EXPIRY_TIME on line %d: %w", i+1, err)
			}
			parsed.REFRESH_EXPIRY_TIME = time.Hour * time.Duration(num)
		case "JWT_SECRET":
			parsed.JWT_SECRET = value
		case "DATABASE_URI":
			parsed.DATABASE_URI = value
		case "ADMIN_PASS":
			parsed.ADMIN_PASS = value
		case "WEBHOOK_SECRET":
			parsed.WEBHOOK_SECRET = value
		case "LOG_PATH":
			parsed.LOG_PATH = value
		default:
			return fmt.Errorf("unknown key on line %d: %s", i+1, key)
		}
	}
	if err := testStruct(); err != nil {
		return err
	}
	return nil
}

func Env() *parsed_env {
	return &parsed
}

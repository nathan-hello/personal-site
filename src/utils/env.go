package utils

import (
	"database/sql"
	"fmt"
	"reflect"
	"time"

	"github.com/joho/godotenv"
)

type Dotenv struct {
	DB_URI     string
	JWT_SECRET string
}

type FullConfig struct {
	DB_URI              string
	JWT_SECRET          string
	REFRESH_EXPIRY_TIME time.Duration
	ACCESS_EXPIRY_TIME  time.Duration
}

var g Dotenv

func InitEnv(override *Dotenv) error {
	if override {
		g.DB_URI = override.DB_URI
		g.JWT_SECRET = override.JWT_SECRET
	} else {
		dotenv, err := godotenv.Read(".env")
		if err != nil {
			return err
		}
		g.DB_URI = dotenv["DB_URI"]
		g.JWT_SECRET = dotenv["JWT_SECRET"]
	}

	if g.DB_URI == "" ||
		g.JWT_SECRET == "" {
		return errors.New("DB_URI or JWT_SECRET not provided")
	}

	return nil
}

var Config = FullConfig{
	DB_URI:              g.DB_URI,
	JWT_SECRET:          g.JWT_SECRET,
	REFRESH_EXPIRY_TIME: time.Hour * 72,
	ACCESS_EXPIRY_TIME:  time.Hour * 24,
}

func Env() *FullConfig {
	return &Config
}

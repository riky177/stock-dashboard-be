package config

import (
	"github.com/joho/godotenv"
)

func LoadEnv() error {
	err := godotenv.Load(".env")

	if err != nil {
		panic("Error loading .env file")
	}
	return err
}

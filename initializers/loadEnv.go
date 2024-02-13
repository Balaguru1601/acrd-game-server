package initializers

import (
	"github.com/joho/godotenv"
)

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		println("Error loading .env from godotenv")
	}
}

package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type DBConfig struct {
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
}

type Config struct {
	DB             DBConfig
	AppPort        string
	Workers        int
	CountOfTryings int
}

func Load() *Config {
	// Загружаем .env
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found, using system env variables")
	}

	port, _ := strconv.Atoi(os.Getenv("DB_PORT"))
	workers, _ := strconv.Atoi(os.Getenv("WORKERS"))
	countOfTryings, _ := strconv.Atoi(os.Getenv("COUNT_OF_TRYINGS"))

	return &Config{
		DB: DBConfig{
			DBHost:     os.Getenv("DB_HOST"),
			DBPort:     port,
			DBUser:     os.Getenv("DB_USER"),
			DBPassword: os.Getenv("DB_PASSWORD"),
			DBName:     os.Getenv("DB_NAME"),
		},

		AppPort:        os.Getenv("APP_PORT"),
		Workers:        workers,
		CountOfTryings: countOfTryings,
	}
}

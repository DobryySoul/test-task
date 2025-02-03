package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port           int
	DBHost         string
	DBPort         int
	DBUser         string
	DBPassword     string
	DBName         string
	MigrationsPath string
	APIURL         string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load("C:/Users/posei/test-task/.env")
	if err != nil {
		return nil, err
	}

	port, _ := strconv.Atoi(os.Getenv("PORT"))
	dbPort, _ := strconv.Atoi(os.Getenv("DB_PORT"))

	return &Config{
		Port:           port,
		DBHost:         os.Getenv("DB_HOST"),
		DBPort:         dbPort,
		DBUser:         os.Getenv("DB_USER"),
		DBPassword:     os.Getenv("DB_PASSWORD"),
		DBName:         os.Getenv("DB_NAME"),
		MigrationsPath: os.Getenv("MIGRATIONS_PATH"),
		APIURL:         os.Getenv("API_URL"),
	}, nil
}

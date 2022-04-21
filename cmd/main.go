package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	//Postgres database connection string
	dbSourceName := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=%s connect_timeout=5 statement_timeout=30",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_DBNAME"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_SSL"),
	)

	fmt.Println(dbSourceName) //Remove
}

package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/anoobz/dualread/auth/internal/httpserver"
	"github.com/anoobz/dualread/auth/internal/store"
	"github.com/anoobz/dualread/auth/internal/store/psqlstore"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	logger := log.New(os.Stdout, "Server ", log.Lshortfile|log.Ltime)

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

	db, err := store.NewDatabase(dbSourceName, 5*time.Second)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()

	statementBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		RunWith(db)
	store := psqlstore.NewSqlStore(db, statementBuilder)

	port, err := strconv.Atoi(os.Getenv("SERVER_PORT"))
	if err != nil {
		logger.Fatal(err)
	}

	server := httpserver.NewServer(store, logger, port)
	err = server.Start()
	if err != nil {
		logger.Fatal(err)
	}
}

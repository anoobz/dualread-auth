package psqlstore

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/anoobz/dualread/auth/internal/store"
)

func CreateTestStore(t *testing.T) (*SqlStore, func(...string)) {
	t.Helper()

	//Postgres database connection string
	dbSourceName := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_TEST_DBNAME"),
		os.Getenv("POSTGRES_PASSWORD"),
	)

	db, err := store.NewDatabase(dbSourceName, 5*time.Second)
	if err != nil {
		t.Fatal(err)
	}

	statementBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		RunWith(db)
	store := NewSqlStore(db, statementBuilder)

	return store, func(tables ...string) {
		if len(tables) > 0 {
			if _, err := db.Exec(fmt.Sprintf("TRUNCATE %s CASCADE", strings.Join(tables, ", "))); err != nil {
				t.Fatal(err)
			}
		}

		db.Close()
	}
}

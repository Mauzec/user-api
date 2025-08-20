package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mauzec/user-api/internal/config"
)

var testQueries *Queries
var testDB *pgxpool.Pool

func TestMain(m *testing.M) {
	config, err := config.LoadConfig("app", "env", "../../config")
	if err != nil {
		log.Fatalf("unable to load config:%+v", err)
	}

	ctx := context.Background()
	testDB, err := pgxpool.New(ctx, config.DBSource)
	if err != nil {
		log.Fatalf("unable to connect to db:%+v", err)
	}
	defer testDB.Close()

	testQueries = New(testDB)
	os.Exit(m.Run())
}

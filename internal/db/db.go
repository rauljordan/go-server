package db

import (
	"context"
	"errors"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rauljordan/go-server/models"
)

// Database interface definition for the application.
type Database interface {
	CreateUser(ctx context.Context, email string, passwordHash []byte) (uint64, error)
	User(ctx context.Context, email string) (*models.User, error)
}

// SQLDatabase struct definition.
type SQLDatabase struct {
	db *sqlx.DB
}

// StartDB connection to a db instance given a database url.
func StartDB(ctx context.Context, databaseUrl string) (*SQLDatabase, error) {
	if databaseUrl == "" {
		return nil, errors.New("nil database url")
	}
	db, err := sqlx.Open("postgres", databaseUrl)
	if err != nil {
		return nil, err
	}
	if db == nil {
		return nil, errors.New("nil db connection")
	}
	log.Println("Established db connection")
	apiDB := &SQLDatabase{
		db: db,
	}
	return apiDB, nil
}

// Close the db connection.
func (d *SQLDatabase) Close() error {
	return d.Close()
}

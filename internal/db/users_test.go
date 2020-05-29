package db

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

func TestUser_CRUD(t *testing.T) {
	dbConn, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Unexpected error creating db connection: %v", err)
	}
	sqlxDB := sqlx.NewDb(dbConn, "sqlmock")
	appDB := &SQLDatabase{db: sqlxDB}

	email := "someone@email.com"
	passwordHash := []byte("helloworld")
	userID := uint64(1)
	// Mocking the insert.
	mock.ExpectQuery(
		regexp.QuoteMeta(
			`INSERT INTO users`,
		),
	).WithArgs(email, passwordHash).
		WillReturnRows(
			sqlmock.NewRows([]string{"user_id"}).
				AddRow(1),
		)
	// Mocking the select.
	mock.ExpectQuery(
		regexp.QuoteMeta(
			`SELECT * FROM users WHERE (email = $1)`,
		),
	).WithArgs(email).
		WillReturnRows(
			sqlmock.NewRows(
				[]string{"user_id", "email", "password_hash"},
			).AddRow(userID, email, passwordHash),
		)
	if _, err := appDB.CreateUser(context.Background(), email, passwordHash); err != nil {
		t.Fatal(err)
	}

	if _, err := appDB.User(context.Background(), email); err != nil {
		t.Fatal(err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

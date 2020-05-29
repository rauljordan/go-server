package db

import (
	"context"
	"database/sql"
	"errors"

	"github.com/rauljordan/go-server/models"
)

// CreateUser and save it to the database.
func (d *SQLDatabase) CreateUser(ctx context.Context, email string, passwordHash []byte) (uint64, error) {
	sqlStatement := `
      INSERT INTO users (email, password_hash)
      VALUES ($1, $2)
      RETURNING user_id`
	id := uint64(0)
	err := d.db.QueryRow(sqlStatement, email, passwordHash).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// User retrieves a user object from the database given an email.
func (d *SQLDatabase) User(ctx context.Context, email string) (*models.User, error) {
	sqlStatement := `SELECT * FROM users WHERE (email = $1);`
	user := &models.User{}
	row := d.db.QueryRow(sqlStatement, email)
	err := row.Scan(&user.UserID, &user.Email, &user.PasswordHash)
	if err == sql.ErrNoRows {
		return nil, errors.New("no user found")
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

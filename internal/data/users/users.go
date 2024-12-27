package users

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/thesambayo/digillets-api/internal/constants"
)

// User represents the users table in the database.
type User struct {
	ID        string    `json:"-"`
	PublicID  string    `json:"public_id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  password  `json:"-"`
	Activated bool      `json:"activated"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Version   int       `json:"version"`
}

var AnonymousUser = &User{}

// Check if a user instance is the AnonymousUser.
func (user *User) IsAnonymous() bool {
	return user == AnonymousUser
}

type UserModel struct {
	DB *sql.DB
}

func (userModel UserModel) Insert(user *User) (*User, error) {
	query := `
    INSERT INTO users
      (public_id, name, email, password_hash, activated)
    VALUES ($1, $2, $3, $4, $5)
    RETURNING id, created_at, version`

	args := []interface{}{user.PublicID, user.Name, user.Email, user.Password.hash, user.Activated}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// If the table already contains a record with this email address, then when we try
	// to perform the insert there will be a violation of the UNIQUE "user_email_key"
	// constraint that we set up in the previous chapter. We check for this error
	// specifically, and return custom ErrDuplicateEmail error instead.
	err := userModel.DB.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.CreatedAt, &user.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return nil, constants.ErrDuplicateEmail
		default:
			return nil, err
		}
	}

	return user, nil
}

func (userModel UserModel) GetByPublicId(publicID string) (*User, error) {
	query := `
    SELECT
			users.id,
			users.public_id,
			users.name,
			users.email,
			users.password_hash,
			users.activated,
			users.created_at,
			users.updated_at,
			users.version
    FROM users
    WHERE users.public_id = $1
  `

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var user User

	err := userModel.DB.QueryRowContext(ctx, query, publicID).Scan(
		&user.ID,
		&user.PublicID,
		&user.Name,
		&user.Email,
		&user.Password.hash,
		&user.Activated,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, constants.ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

func (userModel UserModel) GetByEmail(email string) (*User, error) {
	query := `
    SELECT
			users.id,
			users.public_id,
			users.name,
			users.email,
			users.password_hash,
			users.activated,
			users.created_at,
			users.updated_at,
			users.version
    FROM users
    WHERE users.email = $1
  `

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var user User

	err := userModel.DB.QueryRowContext(ctx, query, strings.ToLower(email)).Scan(
		&user.ID,
		&user.PublicID,
		&user.Name,
		&user.Email,
		&user.Password.hash,
		&user.Activated,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, constants.ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

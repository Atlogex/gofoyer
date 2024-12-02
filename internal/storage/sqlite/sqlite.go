package sqlite

import (
	"atlogex/gofoyer/internal/domain/models"
	"atlogex/gofoyer/internal/storage"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

// App returns app by id.
func (s *Storage) App(ctx context.Context, id int) (models.App, error) {
	const op = "storage.sqlite.App"

	stmt, err := s.db.Prepare("SELECT id, name, secret FROM apps WHERE id = ?")
	if err != nil {
		return models.App{}, fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx, id)

	var app models.App
	err = row.Scan(&app.ID, &app.Name, &app.Secret)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.App{}, fmt.Errorf("%s: %w", op, storage.ErrAppNotFound)
		}

		return models.App{}, fmt.Errorf("%s: %w", op, err)
	}

	return app, nil
}

func New(storagePath string) (*Storage, error) {
	const operation = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database %s: %w", operation, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveUser(
	ctx context.Context,
	email string,
	passHash []byte,
) (int64, error) {
	const operation = "storage.sqlite.SaveUser"

	stmt, err := s.db.Prepare("INSERT INTO users (email, pass_hash) values (?, ?)")
	if err != nil {
		return 0, fmt.Errorf("failed to prepare statement %s: %w", operation, err)
	}
	defer stmt.Close()

	result, err := stmt.Exec(email, passHash)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && errors.Is(sqlite3.ErrConstraintUnique, sqliteErr.ExtendedCode) {
			return 0, fmt.Errorf("failed to save user %s: %w", operation, storage.ErrUserExists)
		}

		return 0, fmt.Errorf("failed to execute statement %s: %w", operation, err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert id %s: %w", operation, err)
	}

	return id, nil
}

func (s *Storage) User(ctx context.Context, email string) (models.User, error) {
	const operation = "storage.sqlite.User"

	stmt, err := s.db.Prepare("SELECT id, email, pass_hash FROM users WHERE email = ?")
	if err != nil {
		return models.User{}, fmt.Errorf("failed to prepare statement %s: %w", operation, err)
	}

	row := stmt.QueryRowContext(ctx, email)

	var user models.User
	err = row.Scan(&user.ID, &user.Email, &user.PassHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, fmt.Errorf("failed to get user %s: %w", operation, storage.ErrUserNotFound)
		}

		return models.User{}, fmt.Errorf("failed to scan user %s: %w", operation, err)
	}

	return user, nil
}

func (s *Storage) IsAdmin(ctx context.Context, userId int64) (bool, error) {
	const operation = "storage.sqlite.isAdmin"

	stmt, err := s.db.Prepare("SELECT is_admin FROM users WHERE id = ?")
	if err != nil {
		return false, fmt.Errorf("failed to prepare statement %s: %w", operation, err)
	}

	row := stmt.QueryRowContext(ctx, userId)

	var isAdmin bool
	err = row.Scan(&isAdmin)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("failed to get user %s: %w", operation, storage.ErrUserNotFound)
		}

		return false, fmt.Errorf("failed to scan user %s: %w", operation, err)
	}

	return isAdmin, nil
}

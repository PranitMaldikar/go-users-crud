package services

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"go-users-crud/models"
)

type UserService struct {
	DB *sql.DB
}

func NewUserService(db *sql.DB) *UserService {
	return &UserService{DB: db}
}

func (s *UserService) CreateUser(ctx context.Context, req models.CreateUserRequest) (models.User, error) {
	req.Name = strings.TrimSpace(req.Name)
	req.Email = strings.TrimSpace(req.Email)

	if req.Name == "" || req.Email == "" {
		return models.User{}, ErrBadRequest
	}

	cctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var u models.User
	err := s.DB.QueryRowContext(cctx, `
		INSERT INTO users (name, email)
		VALUES ($1, $2)
		RETURNING id, name, email, created_at
	`, req.Name, req.Email).Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt)

	if err != nil {
		if isDuplicate(err) {
			return models.User{}, ErrConflict
		}
		return models.User{}, err
	}
	return u, nil
}

func (s *UserService) ListUsers(ctx context.Context) ([]models.User, error) {
	cctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	rows, err := s.DB.QueryContext(cctx, `
		SELECT id, name, email, created_at
		FROM users
		ORDER BY id ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (s *UserService) GetUser(ctx context.Context, id int) (models.User, error) {
	cctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var u models.User
	err := s.DB.QueryRowContext(cctx, `
		SELECT id, name, email, created_at
		FROM users
		WHERE id = $1
	`, id).Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, ErrNotFound
		}
		return models.User{}, err
	}
	return u, nil
}

func (s *UserService) UpdateUser(ctx context.Context, id int, req models.UpdateUserRequest) (models.User, error) {
	req.Name = strings.TrimSpace(req.Name)
	req.Email = strings.TrimSpace(req.Email)

	if req.Name == "" || req.Email == "" {
		return models.User{}, ErrBadRequest

	}

	cctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var u models.User
	err := s.DB.QueryRowContext(cctx, `
		UPDATE users
		SET name = $1, email = $2
		WHERE id = $3
		RETURNING id, name, email, created_at
	`, req.Name, req.Email, id).Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, ErrNotFound
		}
		if isDuplicate(err) {
			return models.User{}, ErrConflict
		}
		return models.User{}, err
	}
	return u, nil
}

func (s *UserService) DeleteUser(ctx context.Context, id int) error {
	cctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	res, err := s.DB.ExecContext(cctx, `DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrNotFound
	}
	return nil
}

// ---- error mapping helpers (simple) ----

var (
	ErrBadRequest = errors.New("bad request")
	ErrNotFound   = errors.New("not found")
	ErrConflict   = errors.New("conflict")
)

func isDuplicate(err error) bool {
	// quick-and-dirty check; good enough for now
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "duplicate") || strings.Contains(msg, "unique")
}

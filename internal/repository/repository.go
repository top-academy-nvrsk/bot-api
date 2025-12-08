package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"bot-api/internal/domain"

	_ "modernc.org/sqlite"
)

// UserRepository - интерфейс с экспортированными именами функций
type UserRepository interface {
	CreateTables() error

	InsertUser(ctx context.Context, u domain.UserRequest) (int, error)
	GetUser(ctx context.Context, id int) (domain.User, error)
	UpdateUser(ctx context.Context, id int, u domain.UserRequest) error

	InsertAnquette(ctx context.Context, a domain.AnquetteRequest) (int, error)
	GetAnquette(ctx context.Context, id int) (domain.Anquette, error)
	UpdateAnquette(ctx context.Context, id int, a domain.AnquetteRequest) error
	DeleteAnquette(ctx context.Context, id int) error
}

type Storage struct {
	db *sql.DB
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{db: db}
}

func (s *Storage) CreateTables() error {
	// ... (логика создания таблиц)
	usersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		tg_id INTEGER NOT NULL UNIQUE,
		tg_username TEXT,
		anquette_id INTEGER
	);`
	_, err := s.db.Exec(usersTable)
	if err != nil {
		return fmt.Errorf("repository: failed to create users table: %w", err)
	}
	anquettesTable := `
	CREATE TABLE IF NOT EXISTS anquettes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		age INTEGER NOT NULL,
		city TEXT,
		gender TEXT,
		preferences TEXT,
		description TEXT NOT NULL
	);`
	_, err = s.db.Exec(anquettesTable)
	if err != nil {
		return fmt.Errorf("repository: failed to create anquettes table: %w", err)
	}

	log.Println("INFO: Таблицы БД успешно инициализированы.")
	return nil
}

// --- Методы User с экспортированными именами ---

func (s *Storage) InsertUser(ctx context.Context, u domain.UserRequest) (int, error) {
	res, err := s.db.ExecContext(ctx, "INSERT INTO users(tg_id, tg_username, anquette_id) values(?, ?, ?)",
		u.TgID, u.TgUsername, u.AnquetteID)
	if err != nil {
		return 0, fmt.Errorf("repository: failed to insert user: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("repository: failed to get last insert ID: %w", err)
	}
	return int(id), nil
}

func (s *Storage) GetUser(ctx context.Context, id int) (domain.User, error) {
	var u domain.User
	row := s.db.QueryRowContext(ctx, "SELECT id, tg_id, tg_username, anquette_id FROM users WHERE id = ?", id)

	err := row.Scan(&u.ID, &u.TgID, &u.TgUsername, &u.AnquetteID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, sql.ErrNoRows
		}
		return domain.User{}, fmt.Errorf("repository: failed scanning user: %w", err)
	}
	return u, nil
}

func (s *Storage) UpdateUser(ctx context.Context, id int, u domain.UserRequest) error {
	res, err := s.db.ExecContext(ctx, "UPDATE users SET tg_id = ?, tg_username = ?, anquette_id = ? WHERE id = ?",
		u.TgID, u.TgUsername, u.AnquetteID, id)
	if err != nil {
		return fmt.Errorf("repository: failed to execute update user: %w", err)
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// --- Методы Anquette с экспортированными именами ---

func (s *Storage) InsertAnquette(ctx context.Context, a domain.AnquetteRequest) (int, error) {
	res, err := s.db.ExecContext(ctx, "INSERT INTO anquettes(name, age, city, gender, preferences, description) values(?, ?, ?, ?, ?, ?)",
		a.Name, a.Age, a.City, a.Gender, a.Preferences, a.Description)
	if err != nil {
		return 0, fmt.Errorf("repository: failed to insert anquette: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("repository: failed to get last insert ID for anquette: %w", err)
	}
	return int(id), nil
}

func (s *Storage) GetAnquette(ctx context.Context, id int) (domain.Anquette, error) {
	var a domain.Anquette
	row := s.db.QueryRowContext(ctx, "SELECT id, name, age, city, gender, preferences, description FROM anquettes WHERE id = ?", id)
	err := row.Scan(&a.ID, &a.Name, &a.Age, &a.City, &a.Gender, &a.Preferences, &a.Description)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Anquette{}, sql.ErrNoRows
		}
		return domain.Anquette{}, fmt.Errorf("repository: failed scanning anquette: %w", err)
	}
	return a, nil
}

func (s *Storage) UpdateAnquette(ctx context.Context, id int, a domain.AnquetteRequest) error {
	res, err := s.db.ExecContext(ctx, "UPDATE anquettes SET name = ?, age = ?, city = ?, gender = ?, preferences = ?, description = ? WHERE id = ?",
		a.Name, a.Age, a.City, a.Gender, a.Preferences, a.Description, id)
	if err != nil {
		return fmt.Errorf("repository: failed to execute update anquette: %w", err)
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (s *Storage) DeleteAnquette(ctx context.Context, id int) error {
	res, err := s.db.ExecContext(ctx, "DELETE FROM anquettes WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("repository: failed to execute delete anquette: %w", err)
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

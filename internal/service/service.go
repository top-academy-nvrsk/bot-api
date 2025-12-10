package service

import (
	"bot-api/internal/domain"
	"bot-api/internal/repository"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// Стандартизированные доменные ошибки
var (
	ErrNotFound         = errors.New("item not found")
	ErrValidationFailed = errors.New("validation failed")
	ErrAlreadyExists    = errors.New("item already exists")
)

// UserService - интерфейс с экспортированными именами функций
type UserService interface {
	InsertUser(ctx context.Context, req domain.UserRequest) (int, error)  // Экспортировано
	GetUser(ctx context.Context, id int) (domain.User, error)             // Экспортировано
	UpdateUser(ctx context.Context, id int, req domain.UserRequest) error // Экспортировано

	InsertAnquette(ctx context.Context, req domain.AnquetteRequest) (int, error)  // Экспортировано
	GetAnquette(ctx context.Context, id int) (domain.Anquette, error)             // Экспортировано
	UpdateAnquette(ctx context.Context, id int, req domain.AnquetteRequest) error // Экспортировано
	DeleteAnquette(ctx context.Context, id int) error                             // Экспортировано
}

// ServiceImpl - реализация сервиса, зависит от Repository
type ServiceImpl struct {
	Repo repository.UserRepository
}

func NewService(repo repository.UserRepository) *ServiceImpl {
	return &ServiceImpl{Repo: repo}
}

// --- Методы User с экспортированными именами ---

func (s *ServiceImpl) InsertUser(ctx context.Context, req domain.UserRequest) (int, error) {
	// Вызов экспортированного метода
	newID, err := s.Repo.InsertUser(ctx, req)
	if err != nil {
		return 0, fmt.Errorf("service: failed to insert user: %w", err)
	}
	return newID, nil
}

func (s *ServiceImpl) GetUser(ctx context.Context, tg_id int) (domain.User, error) {
	// Вызов экспортированного метода
	u, err := s.Repo.GetUser(ctx, tg_id)
	if err != nil {
		// Преобразуем ошибку БД в доменную ошибку ErrNotFound
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, fmt.Errorf("service: user not found: %w", ErrNotFound)
		}
		return domain.User{}, fmt.Errorf("service: failed to get user: %w", err)
	}
	return u, nil
}

func (s *ServiceImpl) UpdateUser(ctx context.Context, tg_id int, req domain.UserRequest) error {
	// Вызов экспортированного метода
	err := s.Repo.UpdateUser(ctx, tg_id, req)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("service: user not found for update: %w", ErrNotFound)
		}
		return fmt.Errorf("service: failed to update user: %w", err)
	}
	return nil
}

// --- Методы Anquette с экспортированными именами ---

func (s *ServiceImpl) InsertAnquette(ctx context.Context, req domain.AnquetteRequest) (int, error) {
	// Бизнес-валидация
	if len(req.Description) < 50 {
		return 0, fmt.Errorf("service: validation failed: %w", ErrValidationFailed)
	}

	// Вызов экспортированного метода
	newID, err := s.Repo.InsertAnquette(ctx, req)
	if err != nil {
		return 0, fmt.Errorf("service: failed to insert anquette: %w", err)
	}
	return newID, nil
}

func (s *ServiceImpl) GetAnquette(ctx context.Context, id int) (domain.Anquette, error) {
	// Вызов экспортированного метода
	a, err := s.Repo.GetAnquette(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Anquette{}, fmt.Errorf("service: anquette not found: %w", ErrNotFound)
		}
		return domain.Anquette{}, fmt.Errorf("service: failed to get anquette: %w", err)
	}
	return a, nil
}

func (s *ServiceImpl) UpdateAnquette(ctx context.Context, id int, req domain.AnquetteRequest) error {
	// Бизнес-валидация
	if len(req.Description) < 50 {
		return fmt.Errorf("service: validation failed: %w", ErrValidationFailed)
	}

	// Вызов экспортированного метода
	err := s.Repo.UpdateAnquette(ctx, id, req)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("service: anquette not found for update: %w", ErrNotFound)
		}
		return fmt.Errorf("service: failed to update anquette: %w", err)
	}
	return nil
}

func (s *ServiceImpl) DeleteAnquette(ctx context.Context, id int) error {
	// Вызов экспортированного метода
	err := s.Repo.DeleteAnquette(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("service: anquette not found for delete: %w", ErrNotFound)
		}
		return fmt.Errorf("service: failed to delete anquette: %w", err)
	}
	return nil
}

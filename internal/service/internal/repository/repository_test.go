package repository_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"bot-api/internal/domain"
	"bot-api/internal/repository"

	_ "modernc.org/sqlite" // Используем modernc.org/sqlite
)

// newTestStorage - хелпер для инициализации временной БД
func newTestStorage(t *testing.T) *repository.Storage {
	// Открываем in-memory SQLite (БД существует только в памяти во время выполнения)
	db, err := sql.Open("sqlite", "file::memory:?cache=shared")
	if err != nil {
		t.Fatalf("Не удалось открыть БД: %v", err)
	}

	storage := repository.NewStorage(db)
	if err := storage.CreateTables(); err != nil {
		t.Fatalf("Не удалось создать таблицы: %v", err)
	}
	return storage
}

// --- ТЕСТЫ USER ---

func TestStorage_InsertAndGetUser_Success(t *testing.T) {
	s := newTestStorage(t)
	ctx := context.Background()

	userReq := domain.UserRequest{TgID: 456, TgUsername: "integration_test_user", AnquetteID: 0}

	// 1. Вставка
	newID, err := s.InsertUser(ctx, userReq)
	if err != nil {
		t.Fatalf("InsertUser провалился: %v", err)
	}

	// 2. Получение
	user, err := s.GetUser(ctx, newID)
	if err != nil {
		t.Fatalf("GetUser провалился: %v", err)
	}

	if user.TgUsername != userReq.TgUsername {
		t.Errorf("Несовпадение username: ожидали %s, получили %s", userReq.TgUsername, user.TgUsername)
	}
}

func TestStorage_UpdateUser_Success(t *testing.T) {
	s := newTestStorage(t)
	ctx := context.Background()

	// 1. Вставка начальной записи
	initialReq := domain.UserRequest{TgID: 100, TgUsername: "old_name"}
	id, _ := s.InsertUser(ctx, initialReq)

	// 2. Обновление
	updateReq := domain.UserRequest{TgID: 101, TgUsername: "new_name"}
	err := s.UpdateUser(ctx, id, updateReq)
	if err != nil {
		t.Fatalf("UpdateUser провалился: %v", err)
	}

	// 3. Проверка
	updatedUser, _ := s.GetUser(ctx, id)
	if updatedUser.TgUsername != "new_name" {
		t.Errorf("Username не был обновлен")
	}
}

// --- ТЕСТЫ ANQUETTE ---

func TestStorage_InsertAndGetAnquette_Success(t *testing.T) {
	s := newTestStorage(t)
	ctx := context.Background()

	anquetteReq := domain.AnquetteRequest{
		Name: "Тестовая анкета", Age: 25, Description: "Очень длинное описание для теста.",
	}

	// 1. Вставка
	newID, err := s.InsertAnquette(ctx, anquetteReq)
	if err != nil {
		t.Fatalf("InsertAnquette провалился: %v", err)
	}

	// 2. Получение
	ank, err := s.GetAnquette(ctx, newID)
	if err != nil {
		t.Fatalf("GetAnquette провалился: %v", err)
	}

	if ank.Name != "Тестовая анкета" || ank.Age != 25 {
		t.Errorf("Данные анкеты не совпадают")
	}
}

func TestStorage_DeleteAnquette_NotFound(t *testing.T) {
	s := newTestStorage(t)
	ctx := context.Background()

	// Пробуем удалить несуществующую запись
	err := s.DeleteAnquette(ctx, 999)

	// Ожидаем ошибку sql.ErrNoRows, так как ни одна строка не была затронута
	if !errors.Is(err, sql.ErrNoRows) {
		t.Errorf("Ожидали sql.ErrNoRows, получили %v", err)
	}
}

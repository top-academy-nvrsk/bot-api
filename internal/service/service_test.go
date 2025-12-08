package service_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"bot-api/internal/domain"
	"bot-api/internal/repository"
	"bot-api/internal/service"
)

// MockRepo - заглушка, реализующая repository.UserRepository.
// Используется для контроля возвращаемых значений и ошибок.
type MockRepo struct {
	repository.UserRepository // Встраиваем интерфейс, чтобы не писать все методы, а переопределять только нужные.

	InsertUserFunc     func(ctx context.Context, u domain.UserRequest) (int, error)
	GetUserFunc        func(ctx context.Context, id int) (domain.User, error)
	InsertAnquetteFunc func(ctx context.Context, a domain.AnquetteRequest) (int, error)
	DeleteAnquetteFunc func(ctx context.Context, id int) error
}

// Переопределяем только те методы, которые нам нужны для тестов
func (m *MockRepo) InsertUser(ctx context.Context, u domain.UserRequest) (int, error) {
	return m.InsertUserFunc(ctx, u)
}
func (m *MockRepo) GetUser(ctx context.Context, id int) (domain.User, error) {
	return m.GetUserFunc(ctx, id)
}
func (m *MockRepo) InsertAnquette(ctx context.Context, a domain.AnquetteRequest) (int, error) {
	return m.InsertAnquetteFunc(ctx, a)
}
func (m *MockRepo) DeleteAnquette(ctx context.Context, id int) error {
	return m.DeleteAnquetteFunc(ctx, id)
}

// Для остальных методов (UpdateUser, GetAnquette, UpdateAnquette) будет использована базовая реализация,
// если они не переопределены, но для чистоты теста можно определить все, чтобы не было nil-указателей
// при случайном вызове. В данном примере оставляем только нужные для демонстрации.

// --- ТЕСТЫ USER ---

func TestServiceImpl_GetUser_Success(t *testing.T) {
	expectedUser := domain.User{ID: 1, TgID: 123, TgUsername: "testuser"}
	mockRepo := &MockRepo{
		GetUserFunc: func(ctx context.Context, id int) (domain.User, error) {
			return expectedUser, nil
		},
	}
	svc := service.NewService(mockRepo)

	user, err := svc.GetUser(context.Background(), 1)

	if err != nil {
		t.Fatalf("Ожидали отсутствие ошибки, получили: %v", err)
	}
	if user.TgUsername != expectedUser.TgUsername {
		t.Errorf("Ожидали username %s, получили %s", expectedUser.TgUsername, user.TgUsername)
	}
}

func TestServiceImpl_GetUser_NotFound(t *testing.T) {
	mockRepo := &MockRepo{
		GetUserFunc: func(ctx context.Context, id int) (domain.User, error) {
			// Имитируем ошибку БД "запись не найдена"
			return domain.User{}, sql.ErrNoRows
		},
	}
	svc := service.NewService(mockRepo)

	_, err := svc.GetUser(context.Background(), 99)

	// Проверяем, что ошибка БД преобразована в доменную ошибку service.ErrNotFound
	if !errors.Is(err, service.ErrNotFound) {
		t.Errorf("Ожидали ошибку service.ErrNotFound, получили: %v", err)
	}
}

// --- ТЕСТЫ ANQUETTE ---

func TestServiceImpl_InsertAnquette_ValidationFail(t *testing.T) {
	// Проверяем бизнес-логику: описание должно быть длинным
	req := domain.AnquetteRequest{Description: "Короткое описание"}

	// MockRepo не вызывается, поэтому можно передать nil или пустой MockRepo
	svc := service.NewService(&MockRepo{})

	_, err := svc.InsertAnquette(context.Background(), req)

	// Проверяем, что сработала валидация
	if !errors.Is(err, service.ErrValidationFailed) {
		t.Errorf("Ожидали ошибку service.ErrValidationFailed, получили: %v", err)
	}
}

func TestServiceImpl_InsertAnquette_Success(t *testing.T) {
	req := domain.AnquetteRequest{
		Description: "Это очень длинное описание, которое точно пройдет проверку валидации и будет вставлено.", // > 50 символов
	}
	mockRepo := &MockRepo{
		InsertAnquetteFunc: func(ctx context.Context, a domain.AnquetteRequest) (int, error) {
			return 5, nil // Имитируем успешную вставку с ID 5
		},
	}
	svc := service.NewService(mockRepo)

	newID, err := svc.InsertAnquette(context.Background(), req)

	if err != nil {
		t.Fatalf("Ожидали отсутствие ошибки, получили: %v", err)
	}
	if newID != 5 {
		t.Errorf("Ожидали ID 5, получили %d", newID)
	}
}

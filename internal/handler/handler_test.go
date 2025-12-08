package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"bot-api/internal/domain"
	"bot-api/internal/handler"
	"bot-api/internal/service"
)

// MockService - заглушка, реализующая service.UserService
type MockService struct {
	service.UserService // Встраиваем интерфейс

	InsertUserFunc     func(ctx context.Context, req domain.UserRequest) (int, error)
	GetUserFunc        func(ctx context.Context, id int) (domain.User, error)
	UpdateUserFunc     func(ctx context.Context, id int, req domain.UserRequest) error
	InsertAnquetteFunc func(ctx context.Context, req domain.AnquetteRequest) (int, error)
	GetAnquetteFunc    func(ctx context.Context, id int) (domain.Anquette, error)
	UpdateAnquetteFunc func(ctx context.Context, id int, req domain.AnquetteRequest) error
	DeleteAnquetteFunc func(ctx context.Context, id int) error
}

func (m *MockService) InsertUser(ctx context.Context, req domain.UserRequest) (int, error) {
	return m.InsertUserFunc(ctx, req)
}
func (m *MockService) GetUser(ctx context.Context, id int) (domain.User, error) {
	return m.GetUserFunc(ctx, id)
}
func (m *MockService) UpdateUser(ctx context.Context, id int, req domain.UserRequest) error {
	return m.UpdateUserFunc(ctx, id, req)
}
func (m *MockService) InsertAnquette(ctx context.Context, req domain.AnquetteRequest) (int, error) {
	return m.InsertAnquetteFunc(ctx, req)
}
func (m *MockService) GetAnquette(ctx context.Context, id int) (domain.Anquette, error) {
	return m.GetAnquetteFunc(ctx, id)
}
func (m *MockService) UpdateAnquette(ctx context.Context, id int, req domain.AnquetteRequest) error {
	return m.UpdateAnquetteFunc(ctx, id, req)
}
func (m *MockService) DeleteAnquette(ctx context.Context, id int) error {
	return m.DeleteAnquetteFunc(ctx, id)
}

// checkResponseCode - Хелпер для проверки HTTP-кода
func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Ожидали код ответа %d. Получили %d", expected, actual)
	}
}

// --- ТЕСТЫ USER ---

func TestCreateUserHandler_Success(t *testing.T) {
	reqBody := `{"tg_id": 12345, "tg_username": "test_user", "anquette_id": 0}`
	mockSvc := &MockService{
		InsertUserFunc: func(ctx context.Context, req domain.UserRequest) (int, error) {
			return 5, nil // Имитируем успешное создание с ID 5
		},
	}
	h := handler.NewHandler(mockSvc)

	req, _ := http.NewRequest("POST", "/api/v1/users", bytes.NewBufferString(reqBody))
	rr := httptest.NewRecorder()

	// 1. Инициализация роутера и регистрация хендлера
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v1/users", h.CreateUserHandler)

	// 2. Выполнение запроса
	mux.ServeHTTP(rr, req)

	checkResponseCode(t, http.StatusCreated, rr.Code)

	var resp domain.APIResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Ошибка парсинга JSON: %v", err)
	}

	if resp.Status != "created" || resp.ID != 5 {
		t.Errorf("Ожидали статус 'created' и ID 5, получили %s и ID %d", resp.Status, resp.ID)
	}
}

func TestGetUserHandler_NotFound(t *testing.T) {
	mockSvc := &MockService{
		GetUserFunc: func(ctx context.Context, id int) (domain.User, error) {
			return domain.User{}, service.ErrNotFound
		},
	}
	h := handler.NewHandler(mockSvc)

	// Создаем запрос GET /api/v1/users/999
	req, _ := http.NewRequest("GET", "/api/v1/users/999", nil)
	rr := httptest.NewRecorder()

	// 1. Инициализация роутера
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/users/{id}", h.GetUserHandler) // Роут с переменной

	// 2. Выполнение запроса. mux корректно установит r.PathValue("id") = "999"
	mux.ServeHTTP(rr, req)

	checkResponseCode(t, http.StatusNotFound, rr.Code)

	var resp domain.APIResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Ошибка парсинга JSON: %v", err)
	}

	if resp.Error == "" {
		t.Error("Ожидали сообщение об ошибке, получили пустое")
	}
}

// --- ТЕСТЫ ANQUETTE ---

func TestCreateAnquetteHandler_ValidationFail(t *testing.T) {
	// Service вернет ошибку валидации
	reqBody := `{"name": "Short", "description": "Too Short"}`
	mockSvc := &MockService{
		InsertAnquetteFunc: func(ctx context.Context, req domain.AnquetteRequest) (int, error) {
			return 0, service.ErrValidationFailed
		},
	}
	h := handler.NewHandler(mockSvc)

	req, _ := http.NewRequest("POST", "/api/v1/anquettes", bytes.NewBufferString(reqBody))
	rr := httptest.NewRecorder()

	// 1. Инициализация роутера
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v1/anquettes", h.CreateAnquetteHandler)

	// 2. Выполнение запроса
	mux.ServeHTTP(rr, req)

	// Service Layer: ErrValidationFailed -> Handler Layer: 400 Bad Request
	checkResponseCode(t, http.StatusBadRequest, rr.Code)
}

func TestDeleteAnquetteHandler_Success(t *testing.T) {
	mockSvc := &MockService{
		DeleteAnquetteFunc: func(ctx context.Context, id int) error {
			return nil // Имитируем успешное удаление
		},
	}
	h := handler.NewHandler(mockSvc)

	req, _ := http.NewRequest("DELETE", "/api/v1/anquettes/1", nil)
	rr := httptest.NewRecorder()

	// 1. Инициализация роутера
	mux := http.NewServeMux()
	mux.HandleFunc("DELETE /api/v1/anquettes/{id}", h.DeleteAnquetteHandler)

	// 2. Выполнение запроса
	mux.ServeHTTP(rr, req)

	checkResponseCode(t, http.StatusOK, rr.Code)

	var resp domain.APIResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Ошибка парсинга JSON: %v", err)
	}

	if resp.Status != "deleted" {
		t.Errorf("Ожидали статус 'deleted', получили %s", resp.Status)
	}
}

func TestDeleteAnquetteHandler_InternalError(t *testing.T) {
	mockSvc := &MockService{
		DeleteAnquetteFunc: func(ctx context.Context, id int) error {
			return errors.New("db connection lost") // Имитируем внутреннюю ошибку, НЕ ErrNotFound
		},
	}
	h := handler.NewHandler(mockSvc)

	req, _ := http.NewRequest("DELETE", "/api/v1/anquettes/1", nil)
	rr := httptest.NewRecorder()

	// 1. Инициализация роутера
	mux := http.NewServeMux()
	mux.HandleFunc("DELETE /api/v1/anquettes/{id}", h.DeleteAnquetteHandler)

	// 2. Выполнение запроса
	mux.ServeHTTP(rr, req)

	// Ожидаем 500 Internal Server Error, так как это не ErrNotFound и не ErrValidationFailed
	checkResponseCode(t, http.StatusInternalServerError, rr.Code)
}

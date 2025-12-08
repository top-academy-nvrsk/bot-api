package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"bot-api/internal/domain"
	"bot-api/internal/service"
)

// Handler - структура для DI
type Handler struct {
	Service service.UserService
}

func NewHandler(svc service.UserService) *Handler {
	return &Handler{Service: svc}
}

// sendJSON - Хелпер для отправки JSON-ответа
func sendJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
	}
}

// handleServiceError - централизованная функция для обработки ошибок Service
func handleServiceError(w http.ResponseWriter, err error, resourceName string) {
	log.Printf("WARNING: %s operation failed: %v", resourceName, err)

	if errors.Is(err, service.ErrNotFound) {
		sendJSON(w, http.StatusNotFound, domain.APIResponse{
			Status: "error", Error: fmt.Sprintf("%s не найден", resourceName),
		})
		return
	}
	if errors.Is(err, service.ErrValidationFailed) {
		sendJSON(w, http.StatusBadRequest, domain.APIResponse{
			Status: "error", Error: "Ошибка валидации: " + err.Error(),
		})
		return
	}

	sendJSON(w, http.StatusInternalServerError, domain.APIResponse{
		Status: "error", Error: "Внутренняя ошибка сервера",
	})
}

// --- Методы User с экспортированными именами ---

func (h *Handler) CreateUserHandler(w http.ResponseWriter, r *http.Request) { // Изменено
	var req domain.UserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendJSON(w, http.StatusBadRequest, domain.APIResponse{Status: "error", Error: "Неверный JSON формат"})
		return
	}

	newID, err := h.Service.InsertUser(r.Context(), req)
	if err != nil {
		handleServiceError(w, err, "юзер")
		return
	}

	log.Printf("INFO: User created ID: %d", newID)
	sendJSON(w, http.StatusCreated, domain.APIResponse{Status: "created", ID: newID})
}

func (h *Handler) GetUserHandler(w http.ResponseWriter, r *http.Request) { // Изменено
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		sendJSON(w, http.StatusBadRequest, domain.APIResponse{Status: "error", Error: "ID должен быть числом"})
		return
	}

	u, err := h.Service.GetUser(r.Context(), id)
	if err != nil {
		handleServiceError(w, err, "юзер")
		return
	}

	sendJSON(w, http.StatusOK, domain.APIResponse{Status: "ok", Data: u})
}

func (h *Handler) UpdateUserHandler(w http.ResponseWriter, r *http.Request) { // Изменено
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		sendJSON(w, http.StatusBadRequest, domain.APIResponse{Status: "error", Error: "Неверный ID"})
		return
	}

	var req domain.UserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendJSON(w, http.StatusBadRequest, domain.APIResponse{Status: "error", Error: "JSON ошибка"})
		return
	}

	if err := h.Service.UpdateUser(r.Context(), id, req); err != nil {
		handleServiceError(w, err, "юзер")
		return
	}

	sendJSON(w, http.StatusOK, domain.APIResponse{Status: "updated"})
}

// --- Методы Anquette с экспортированными именами ---

func (h *Handler) CreateAnquetteHandler(w http.ResponseWriter, r *http.Request) { // Изменено
	var req domain.AnquetteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendJSON(w, http.StatusBadRequest, domain.APIResponse{Status: "error", Error: "Неверный JSON"})
		return
	}

	newID, err := h.Service.InsertAnquette(r.Context(), req)
	if err != nil {
		handleServiceError(w, err, "анкета")
		return
	}

	log.Printf("INFO: Anquette created ID: %d", newID)
	sendJSON(w, http.StatusCreated, domain.APIResponse{Status: "created", ID: newID})
}

func (h *Handler) GetAnquetteHandler(w http.ResponseWriter, r *http.Request) { // Изменено
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		sendJSON(w, http.StatusBadRequest, domain.APIResponse{Status: "error", Error: "ID должен быть числом"})
		return
	}

	a, err := h.Service.GetAnquette(r.Context(), id)
	if err != nil {
		handleServiceError(w, err, "анкета")
		return
	}
	sendJSON(w, http.StatusOK, domain.APIResponse{Status: "ok", Data: a})
}

func (h *Handler) UpdateAnquetteHandler(w http.ResponseWriter, r *http.Request) { // Изменено
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		sendJSON(w, http.StatusBadRequest, domain.APIResponse{Status: "error", Error: "Неверный ID"})
		return
	}

	var req domain.AnquetteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendJSON(w, http.StatusBadRequest, domain.APIResponse{Status: "error", Error: "Ошибка JSON"})
		return
	}

	if err := h.Service.UpdateAnquette(r.Context(), id, req); err != nil {
		handleServiceError(w, err, "анкета")
		return
	}
	sendJSON(w, http.StatusOK, domain.APIResponse{Status: "updated"})
}

func (h *Handler) DeleteAnquetteHandler(w http.ResponseWriter, r *http.Request) { // Изменено
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		sendJSON(w, http.StatusBadRequest, domain.APIResponse{Status: "error", Error: "ID должен быть числом"})
		return
	}

	if err := h.Service.DeleteAnquette(r.Context(), id); err != nil {
		handleServiceError(w, err, "анкета")
		return
	}
	sendJSON(w, http.StatusOK, domain.APIResponse{Status: "deleted"})
}

package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

func sendJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
	}
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	var req UserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendJSON(w, http.StatusBadRequest, APIResponse{Status: "error", Error: "Кривой JSON"})
		return
	}

	newID, err := insertUser(req)
	if err != nil {
		log.Printf("ERROR: Failed to create user: %v", err)
		sendJSON(w, http.StatusInternalServerError, APIResponse{Status: "error", Error: "Ошибка создания юзера в БД"})
		return
	}

	log.Printf("INFO: User created ID: %d", newID)
	sendJSON(w, http.StatusCreated, APIResponse{Status: "created", ID: newID})
}

func getUserHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		sendJSON(w, http.StatusBadRequest, APIResponse{Status: "error", Error: "ID должен быть числом"})
		return
	}

	u, err := getUser(id)
	if err != nil {
		log.Printf("WARNING: Failed to get user ID %d: %v", id, err)
		sendJSON(w, http.StatusNotFound, APIResponse{Status: "error", Error: "Нет такого юзера"})
		return
	}

	sendJSON(w, http.StatusOK, APIResponse{Status: "ok", Data: u})
}

func updateUserHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		sendJSON(w, http.StatusBadRequest, APIResponse{Status: "error", Error: "Неверный ID"})
		return
	}

	var req UserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendJSON(w, http.StatusBadRequest, APIResponse{Status: "error", Error: "JSON ошибка"})
		return
	}

	if err := updateUser(id, req); err != nil {
		log.Printf("WARNING: Failed to update user ID %d: %v", id, err)
		sendJSON(w, http.StatusNotFound, APIResponse{Status: "error", Error: err.Error()})
		return
	}

	sendJSON(w, http.StatusOK, APIResponse{Status: "updated"})
}

func createAnquetteHandler(w http.ResponseWriter, r *http.Request) {
	var req AnquetteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendJSON(w, http.StatusBadRequest, APIResponse{Status: "error", Error: "Неверный JSON"})
		return
	}

	if len(req.Description) < 50 {
		sendJSON(w, http.StatusBadRequest, APIResponse{
			Status: "error",
			Error:  "Описание слишком короткое! Минимум 50 символов.",
		})
		return
	}

	newID, err := insertAnquette(req)
	if err != nil {
		log.Printf("ERROR: Failed to create anquette: %v", err)
		sendJSON(w, http.StatusInternalServerError, APIResponse{Status: "error", Error: "Ошибка создания анкеты в БД"})
		return
	}

	log.Printf("INFO: Anquette created ID: %d", newID)
	sendJSON(w, http.StatusCreated, APIResponse{Status: "created", ID: newID})
}

func getAnquetteHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		sendJSON(w, http.StatusBadRequest, APIResponse{Status: "error", Error: "ID должен быть числом"})
		return
	}

	a, err := getAnquette(id)
	if err != nil {
		log.Printf("WARNING: Failed to get anquette ID %d: %v", id, err)
		sendJSON(w, http.StatusNotFound, APIResponse{Status: "error", Error: "Анкета не найдена"})
		return
	}
	sendJSON(w, http.StatusOK, APIResponse{Status: "ok", Data: a})
}

func updateAnquetteHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		sendJSON(w, http.StatusBadRequest, APIResponse{Status: "error", Error: "Неверный ID"})
		return
	}

	var req AnquetteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendJSON(w, http.StatusBadRequest, APIResponse{Status: "error", Error: "Ошибка JSON"})
		return
	}

	if len(req.Description) < 50 {
		sendJSON(w, http.StatusBadRequest, APIResponse{Status: "error", Error: "Описание должно быть 50+ символов"})
		return
	}

	if err := updateAnquette(id, req); err != nil {
		log.Printf("WARNING: Failed to update anquette ID %d: %v", id, err)
		sendJSON(w, http.StatusNotFound, APIResponse{Status: "error", Error: err.Error()})
		return
	}
	sendJSON(w, http.StatusOK, APIResponse{Status: "updated"})
}

func deleteAnquetteHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		sendJSON(w, http.StatusBadRequest, APIResponse{Status: "error", Error: "ID должен быть числом"})
		return
	}

	if err := deleteAnquette(id); err != nil {
		log.Printf("WARNING: Failed to delete anquette ID %d: %v", id, err)
		sendJSON(w, http.StatusNotFound, APIResponse{Status: "error", Error: err.Error()})
		return
	}
	sendJSON(w, http.StatusOK, APIResponse{Status: "deleted"})
}

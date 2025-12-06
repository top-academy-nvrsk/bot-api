package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func newTestRequest(method, path string, body interface{}) *http.Request {
	var buf bytes.Buffer
	if body != nil {
		json.NewEncoder(&buf).Encode(body)
	}
	req := httptest.NewRequest(method, path, &buf)
	req.Header.Set("Content-Type", "application/json")
	return req
}

func TestCreateUserHandler_Success(t *testing.T) {
	testUser := UserRequest{
		TgID:       987654,
		TgUsername: "test_user_ok",
		AnquetteID: 1,
	}

	req := newTestRequest("POST", "/api/v1/users", testUser)
	w := httptest.NewRecorder()

	createUserHandler(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Ожидали код 201 Created, но получили %d", w.Code)
	}

	var response APIResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Не смогли распарсить JSON-ответ: %v", err)
	}

	if response.Status != "created" {
		t.Errorf("Ожидали статус 'created', получили '%s'", response.Status)
	}
}

func TestCreateAnquetteHandler_Success(t *testing.T) {
	validAnquette := AnquetteRequest{
		Name:        "Тест",
		Age:         30,
		City:        "Тестград",
		Gender:      "male",
		Preferences: "Любые",
		Description: "Это тестовое описание, которое точно имеет больше пятидесяти символов, чтобы пройти проверку. Ура! Тест должен быть зеленым! 💚",
	}

	req := newTestRequest("POST", "/api/v1/anquettes", validAnquette)
	w := httptest.NewRecorder()
	createAnquetteHandler(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Ожидали код 201 Created, получили %d. Ответ: %s", w.Code, w.Body.String())
	}
}

func TestCreateAnquetteHandler_ShortDescription_Fail(t *testing.T) {
	invalidAnquette := AnquetteRequest{
		Name:        "Тест",
		Age:         30,
		City:        "Тестград",
		Description: "Короткое.",
	}

	req := newTestRequest("POST", "/api/v1/anquettes", invalidAnquette)
	w := httptest.NewRecorder()
	createAnquetteHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Ожидали код 400 Bad Request, но получили %d", w.Code)
	}
}

func TestGetAnquetteHandler_Success(t *testing.T) {
	testAnquette := AnquetteRequest{
		Name:        "GetTest",
		Age:         40,
		City:        "GetCity",
		Description: "Описание должно быть очень длинным, чтобы тест прошел, иначе мы получим ошибку длины при вставке в хранилище. Длина должна быть больше 50 символов.",
	}
	id := insertAnquette(testAnquette)

	req := newTestRequest("GET", "/api/v1/anquettes", nil)

	req.SetPathValue("id", strconv.Itoa(id))

	w := httptest.NewRecorder()
	getAnquetteHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Ожидали код 200 OK, но получили %d", w.Code)
	}

	var response APIResponse
	json.NewDecoder(w.Body).Decode(&response)

	dataJSON, _ := json.Marshal(response.Data)
	var retrievedAnquette AnquetteRequest
	json.Unmarshal(dataJSON, &retrievedAnquette)

	if retrievedAnquette.Name != testAnquette.Name {
		t.Errorf("Имена не совпадают! Ожидали '%s', получили '%s'", testAnquette.Name, retrievedAnquette.Name)
	}
}

func TestDeleteAnquetteHandler_Success(t *testing.T) {
	testAnquette := AnquetteRequest{
		Name:        "DeleteTest",
		Age:         20,
		Description: "Это длинное описание для тестовой анкеты, которую мы собираемся немедленно удалить. Просто проверка функционала!",
	}
	id := insertAnquette(testAnquette)

	req := newTestRequest("DELETE", "/api/v1/anquettes", nil)
	req.SetPathValue("id", strconv.Itoa(id))

	w := httptest.NewRecorder()
	deleteAnquetteHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Ожидали код 200 OK после удаления, но получили %d. Ответ: %s", w.Code, w.Body.String())
	}

	_, err := getAnquette(id)
	if err == nil {
		t.Error("Анкета должна быть удалена, но функция getAnquette ее нашла!")
	}
}

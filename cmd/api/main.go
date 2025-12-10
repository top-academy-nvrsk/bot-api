package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "modernc.org/sqlite"

	"bot-api/internal/handler"
	"bot-api/internal/repository"
	"bot-api/internal/service"
)

func main() {
	// 1. Инициализация БД
	db, err := sql.Open("sqlite", "/home/creepy0964/bot-api/cmd/api/db/dating_app.db")
	if err != nil {
		log.Fatalf("FATAL: Ошибка открытия БД: %v", err)
	}
	defer db.Close()

	// 2. Сборка Зависимостей (Dependency Injection)

	// Repository: работает с БД
	repo := repository.NewStorage(db)
	if err := repo.CreateTables(); err != nil {
		log.Fatalf("FATAL: Ошибка создания таблиц: %v", err)
	}

	// Service: содержит бизнес-логику и зависит от Repository
	svc := service.NewService(repo)

	// Handler: обрабатывает HTTP и зависит от Service
	h := handler.NewHandler(svc)

	// 3. Настройка Роутера
	mux := http.NewServeMux()

	// Регистрация роутов (используем экспортированные методы h)
	mux.HandleFunc("POST /api/v1/users", h.CreateUserHandler)
	mux.HandleFunc("GET /api/v1/users/{id}", h.GetUserHandler)
	mux.HandleFunc("PUT /api/v1/users/{id}", h.UpdateUserHandler)
	mux.HandleFunc("POST /api/v1/anquettes", h.CreateAnquetteHandler)
	mux.HandleFunc("GET /api/v1/anquettes/{id}", h.GetAnquetteHandler)
	mux.HandleFunc("PUT /api/v1/anquettes/{id}", h.UpdateAnquetteHandler)
	mux.HandleFunc("DELETE /api/v1/anquettes/{id}", h.DeleteAnquetteHandler)

	// 4. Запуск Сервера
	port := ":8080"
	fmt.Printf("Сервер запущен на порту %s\n", port)
	log.Printf("INFO: Starting server on %s", port)

	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatal(err)
	}
}

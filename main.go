package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	if err := InitStorage("./dating_app.db"); err != nil {
		log.Fatalf("Ошибка инициализации хранилища: %v", err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/v1/users", createUserHandler)
	mux.HandleFunc("GET /api/v1/users/{id}", getUserHandler)
	mux.HandleFunc("PUT /api/v1/users/{id}", updateUserHandler)
	mux.HandleFunc("POST /api/v1/anquettes", createAnquetteHandler)
	mux.HandleFunc("GET /api/v1/anquettes/{id}", getAnquetteHandler)
	mux.HandleFunc("PUT /api/v1/anquettes/{id}", updateAnquetteHandler)
	mux.HandleFunc("DELETE /api/v1/anquettes/{id}", deleteAnquetteHandler)

	fmt.Println("Сервер запущен на порту :8080")
	log.Print("Starting server on :8080")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}

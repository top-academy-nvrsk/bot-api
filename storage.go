package main

import (
	"errors"
	"sync"
)

type Storage struct {
	mu          sync.Mutex
	users       map[int]UserRequest
	anquettes   map[int]AnquetteRequest
	userCounter int
	anqCounter  int
}

var storage = Storage{
	users:     make(map[int]UserRequest),
	anquettes: make(map[int]AnquetteRequest),
}

func insertUser(u UserRequest) int {
	storage.mu.Lock()
	defer storage.mu.Unlock()
	storage.userCounter++
	storage.users[storage.userCounter] = u
	return storage.userCounter
}

func getUser(id int) (UserRequest, error) {
	storage.mu.Lock()
	defer storage.mu.Unlock()
	u, ok := storage.users[id]
	if !ok {
		return UserRequest{}, errors.New("юзер не найден")
	}
	return u, nil
}

func updateUser(id int, u UserRequest) error {
	storage.mu.Lock()
	defer storage.mu.Unlock()

	if _, ok := storage.users[id]; !ok {
		return errors.New("юзер не найден, нельзя обновить")
	}
	storage.users[id] = u
	return nil
}

func insertAnquette(a AnquetteRequest) int {
	storage.mu.Lock()
	defer storage.mu.Unlock()
	storage.anqCounter++
	storage.anquettes[storage.anqCounter] = a
	return storage.anqCounter
}

func getAnquette(id int) (AnquetteRequest, error) {
	storage.mu.Lock()
	defer storage.mu.Unlock()
	a, ok := storage.anquettes[id]
	if !ok {
		return AnquetteRequest{}, errors.New("анкета не найдена")
	}
	return a, nil
}

func updateAnquette(id int, a AnquetteRequest) error {
	storage.mu.Lock()
	defer storage.mu.Unlock()

	if _, ok := storage.anquettes[id]; !ok {
		return errors.New("анкета не найдена")
	}
	storage.anquettes[id] = a
	return nil
}

func deleteAnquette(id int) error {
	storage.mu.Lock()
	defer storage.mu.Unlock()

	if _, ok := storage.anquettes[id]; !ok {
		return errors.New("нет такой анкеты")
	}
	delete(storage.anquettes, id)
	return nil
}

package main

import (
	"database/sql"
	"errors"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

var storage *Storage

func InitStorage(dbPath string) error {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}

	storage = &Storage{db: db}

	if err = db.Ping(); err != nil {
		return err
	}

	return storage.createTables()
}

func (s *Storage) createTables() error {
	usersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		tg_id INTEGER NOT NULL UNIQUE,
		tg_username TEXT,
		anquette_id INTEGER
	);`
	_, err := s.db.Exec(usersTable)
	if err != nil {
		return errors.New("ошибка создания таблицы users: " + err.Error())
	}

	anquettesTable := `
	CREATE TABLE IF NOT EXISTS anquettes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		age INTEGER NOT NULL,
		city TEXT,
		gender TEXT,
		preferences TEXT,
		description TEXT NOT NULL
	);`
	_, err = s.db.Exec(anquettesTable)
	if err != nil {
		return errors.New("ошибка создания таблицы anquettes: " + err.Error())
	}

	log.Println("INFO: Таблицы БД успешно инициализированы.")
	return nil
}

func insertUser(u UserRequest) (int, error) {
	stmt, err := storage.db.Prepare("INSERT INTO users(tg_id, tg_username, anquette_id) values(?, ?, ?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(u.TgID, u.TgUsername, u.AnquetteID)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func getUser(id int) (User, error) {
	row := storage.db.QueryRow("SELECT id, tg_id, tg_username, anquette_id FROM users WHERE id = ?", id)
	var u User
	err := row.Scan(&u.ID, &u.TgID, &u.TgUsername, &u.AnquetteID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, errors.New("юзер не найден")
		}
		return User{}, err
	}
	return u, nil
}

func updateUser(id int, u UserRequest) error {
	res, err := storage.db.Exec("UPDATE users SET tg_id = ?, tg_username = ?, anquette_id = ? WHERE id = ?",
		u.TgID, u.TgUsername, u.AnquetteID, id)
	if err != nil {
		return err
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("юзер не найден, нельзя обновить")
	}
	return nil
}

func insertAnquette(a AnquetteRequest) (int, error) {
	stmt, err := storage.db.Prepare("INSERT INTO anquettes(name, age, city, gender, preferences, description) values(?, ?, ?, ?, ?, ?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(a.Name, a.Age, a.City, a.Gender, a.Preferences, a.Description)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func getAnquette(id int) (Anquette, error) {
	row := storage.db.QueryRow("SELECT id, name, age, city, gender, preferences, description FROM anquettes WHERE id = ?", id)
	var a Anquette
	err := row.Scan(&a.ID, &a.Name, &a.Age, &a.City, &a.Gender, &a.Preferences, &a.Description)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Anquette{}, errors.New("анкета не найдена")
		}
		return Anquette{}, err
	}
	return a, nil
}

func updateAnquette(id int, a AnquetteRequest) error {
	res, err := storage.db.Exec("UPDATE anquettes SET name = ?, age = ?, city = ?, gender = ?, preferences = ?, description = ? WHERE id = ?",
		a.Name, a.Age, a.City, a.Gender, a.Preferences, a.Description, id)
	if err != nil {
		return err
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("анкета не найдена")
	}
	return nil
}

func deleteAnquette(id int) error {
	res, err := storage.db.Exec("DELETE FROM anquettes WHERE id = ?", id)
	if err != nil {
		return err
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("нет такой анкеты")
	}
	return nil
}

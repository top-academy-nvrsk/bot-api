package domain

// === Модели БД ===

type User struct {
	ID         int    `json:"id,omitempty"`
	TgID       int64  `json:"tg_id"`
	TgUsername string `json:"tg_username"`
	AnquetteID int    `json:"anquette_id"`
}

type Anquette struct {
	ID          int    `json:"id,omitempty"`
	Name        string `json:"name"`
	Age         int    `json:"age"`
	City        string `json:"city"`
	Gender      string `json:"gender"`
	Preferences string `json:"preferences"`
	Description string `json:"description"`
}

// === Структуры Запросов ===

type UserRequest struct {
	TgID       int64  `json:"tg_id"`
	TgUsername string `json:"tg_username"`
	AnquetteID int    `json:"anquette_id"`
}

type AnquetteRequest struct {
	Name        string `json:"name"`
	Age         int    `json:"age"`
	City        string `json:"city"`
	Gender      string `json:"gender"`
	Preferences string `json:"preferences"`
	Description string `json:"description"`
}

// === Структура Ответа API ===

type APIResponse struct {
	Status string      `json:"status"`
	ID     int         `json:"id,omitempty"`
	Data   interface{} `json:"data,omitempty"`
	Error  string      `json:"error,omitempty"`
}

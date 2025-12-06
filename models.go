package main

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

type APIResponse struct {
	Status string      `json:"status"`
	ID     int         `json:"id,omitempty"`
	Data   interface{} `json:"data,omitempty"`
	Error  string      `json:"error,omitempty"`
}

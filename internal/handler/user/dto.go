package user

// User - модель пользователя для ответа
type User struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Phone  string `json:"phone"`
	Rating int    `json:"rating"`
}

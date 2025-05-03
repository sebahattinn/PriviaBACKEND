package models

// Kullanıcı bilgileri
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// Login isteği
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

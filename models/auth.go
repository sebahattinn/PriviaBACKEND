package models

// Kullanıcı bilgileri
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

// Login isteği
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

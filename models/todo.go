package models

import "time"

// TodoItem modelindeki silme işlemi
type TodoItem struct {
	ID        int        `json:"id"`
	ListID    int        `json:"list_id"`
	Content   string     `json:"content" binding:"required"`
	IsDone    bool       `json:"is_done"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"` // Soft delete alanı
}

// TodoList modelindeki silme işlemi
type TodoList struct {
	ID         int        `json:"id"`
	Name       string     `json:"name"`
	Items      []TodoItem `json:"items"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty"` // Soft delete alanı
	Completion float32    `json:"completion"`
}

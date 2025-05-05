package models

import "time"

// TodoList represents a collection of todo items
type TodoList struct {
	ID         int         `json:"id"`
	UserID     int         `json:"owner_id"`
	Name       string      `json:"name"`
	Items      []*TodoItem `json:"items"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
	DeletedAt  *time.Time  `json:"deleted_at"`
	Completion float32     `json:"completion" default:"0"`
}

// TodoItem represents a single task in a todo list
type TodoItem struct {
	ID        int        `json:"id"`
	ListID    int        `json:"list_id"`
	Content   string     `json:"content" default:""`
	IsDone    bool       `json:"is_done" default:"false"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

type TodoItemUpdate struct {
	Content string `json:"content"`
	IsDone  bool   `json:"is_done"`
}

type TodoItemCreate struct {
	Content string `json:"content"`
	IsDone  bool   `json:"is_done"`
}

type TodoListCreate struct {
	Name string `json:"name" binding:"required"`
}
type TodoListUpdate struct {
	Name string `json:"name"`
}

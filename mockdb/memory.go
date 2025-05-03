package mockdb

import (
	"priviatodolist/models"
	"time"
)

var Users = map[string]string{
	"user1":  "1234",
	"admin1": "admin",
}

var UserRoles = map[string]string{
	"user1":  "user",
	"admin1": "admin",
}

var TodoLists = map[int]*models.TodoList{}
var TodoItems = map[int]*models.TodoItem{}

var TodoListIDCounter = 1
var TodoItemIDCounter = 1

func GetCurrentTime() time.Time {
	return time.Now().UTC()
}

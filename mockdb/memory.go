package mockdb

import (
	"priviatodolist/models"
	"time"
)

// Kullanıcı adı ve şifrelerin bulunduğu veritabanı
var Users = map[string]string{
	"user1":  "1234",
	"admin1": "admin",
	"user2":  "abcd",
	"user3":  "pass123",
	"user4":  "qwerty",
	"user5":  "zxcvbn",
}

// Kullanıcı rollerinin bulunduğu veritabanı
var UserRoles = map[string]string{
	"user1":  "user",
	"admin1": "admin",
	"user2":  "user",
	"user3":  "user",
	"user4":  "user",
	"user5":  "user",
}

// Kullanıcı ID'lerini içeren veritabanı (her kullanıcıya bir benzersiz ID atamak için)
var UserIDs = map[string]int{
	"user1":  1,
	"admin1": 2,
	"user2":  3,
	"user3":  4,
	"user4":  5,
	"user5":  6,
}

// Todo listelerini temsil eden veritabanı
var TodoLists = map[int]*models.TodoList{
	1: {
		ID:     1,
		Name:   "Default List",
		UserID: 1,
		Items: []*models.TodoItem{
			{
				ID:        1,
				ListID:    1,
				Content:   "Default Description",
				IsDone:    false,
				CreatedAt: GetCurrentTime(),
				UpdatedAt: GetCurrentTime(),
				DeletedAt: nil,
			},
		},
		Completion: 0,
		CreatedAt:  GetCurrentTime(),
		UpdatedAt:  GetCurrentTime(),
		DeletedAt:  nil,
	},
	2: {
		ID:     2,
		Name:   "Work Tasks",
		UserID: 2,
		Items: []*models.TodoItem{
			{
				ID:        2,
				ListID:    2,
				Content:   "Prepare presentation",
				IsDone:    false,
				CreatedAt: GetCurrentTime(),
				UpdatedAt: GetCurrentTime(),
				DeletedAt: nil,
			},
			{
				ID:        3,
				ListID:    2,
				Content:   "Send emails",
				IsDone:    true,
				CreatedAt: GetCurrentTime(),
				UpdatedAt: GetCurrentTime(),
				DeletedAt: nil,
			},
		},
		Completion: 50,
		CreatedAt:  GetCurrentTime(),
		UpdatedAt:  GetCurrentTime(),
		DeletedAt:  nil,
	},
	3: {
		ID:     3,
		Name:   "Shopping List",
		UserID: 1,
		Items: []*models.TodoItem{
			{
				ID:        4,
				ListID:    3,
				Content:   "Buy milk",
				IsDone:    false,
				CreatedAt: GetCurrentTime(),
				UpdatedAt: GetCurrentTime(),
				DeletedAt: nil,
			},
			{
				ID:        5,
				ListID:    3,
				Content:   "Buy bread",
				IsDone:    true,
				CreatedAt: GetCurrentTime(),
				UpdatedAt: GetCurrentTime(),
				DeletedAt: nil,
			},
		},
		Completion: 20,
		CreatedAt:  GetCurrentTime(),
		UpdatedAt:  GetCurrentTime(),
		DeletedAt:  nil,
	},
}

// Todo maddelerini temsil eden veritabanı
var TodoItems = map[int]*models.TodoItem{
	1: {
		ID:        1,
		ListID:    1,
		Content:   "Default Description",
		IsDone:    false,
		CreatedAt: GetCurrentTime(),
		UpdatedAt: GetCurrentTime(),
		DeletedAt: nil,
	},
	2: {
		ID:        2,
		ListID:    2,
		Content:   "Prepare presentation",
		IsDone:    false,
		CreatedAt: GetCurrentTime(),
		UpdatedAt: GetCurrentTime(),
		DeletedAt: nil,
	},
	3: {
		ID:        3,
		ListID:    2,
		Content:   "Send emails",
		IsDone:    true,
		CreatedAt: GetCurrentTime(),
		UpdatedAt: GetCurrentTime(),
		DeletedAt: nil,
	},
	4: {
		ID:        4,
		ListID:    3,
		Content:   "Buy milk",
		IsDone:    false,
		CreatedAt: GetCurrentTime(),
		UpdatedAt: GetCurrentTime(),
		DeletedAt: nil,
	},
	5: {
		ID:        5,
		ListID:    3,
		Content:   "Buy bread",
		IsDone:    true,
		CreatedAt: GetCurrentTime(),
		UpdatedAt: GetCurrentTime(),
		DeletedAt: nil,
	},
}

// ID sayaçları
var TodoListIDCounter = 4
var TodoItemIDCounter = 6

// Şu anki UTC zamanını döndüren fonksiyon
func GetCurrentTime() time.Time {
	return time.Now()
}

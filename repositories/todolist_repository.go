package repositories

import (
	"errors"
	"priviatodolist/mockdb"
	"priviatodolist/models"
)

// TodoList'i ID ile bul
func GetTodoListByID(listID int) (*models.TodoList, error) {
	list, exists := mockdb.TodoLists[listID]
	if !exists {
		return nil, errors.New("list not found")
	}
	return list, nil
}

// Yeni bir TodoList oluştur
func CreateTodoList(newList *models.TodoList) (*models.TodoList, error) {
	mockdb.TodoLists[newList.ID] = newList
	return newList, nil
}

// TodoList güncelle
func UpdateTodoList(listID int, updatedList *models.TodoList) (*models.TodoList, error) {
	mockdb.TodoLists[listID] = updatedList
	return updatedList, nil
}

// Kullanıcıya ait TodoList'leri getir
func GetTodoListsByUserID(userID int, includeDeleted bool) ([]*models.TodoList, error) {
	var lists []*models.TodoList

	for _, list := range mockdb.TodoLists {
		if list.UserID == userID && (includeDeleted || list.DeletedAt == nil) {
			listCopy := *list

			var activeItems []*models.TodoItem
			for _, item := range list.Items {
				if item.DeletedAt == nil {
					activeItems = append(activeItems, item)
				}
			}

			listCopy.Items = activeItems

			lists = append(lists, &listCopy)
		}
	}

	return lists, nil
}

// Tüm TodoList'leri getir (Admin için)
func GetAllTodoLists(includeDeleted bool) ([]*models.TodoList, error) {
	var lists []*models.TodoList
	for _, list := range mockdb.TodoLists {
		if includeDeleted || list.DeletedAt == nil {
			lists = append(lists, list)
		}
	}
	return lists, nil
}

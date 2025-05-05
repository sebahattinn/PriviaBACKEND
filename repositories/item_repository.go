package repositories

import (
	"errors"
	"priviatodolist/mockdb"
	"priviatodolist/models"
	"time"
)

func CreateItem(item *models.TodoItem) (*models.TodoItem, error) {
	item.ID = mockdb.TodoItemIDCounter
	mockdb.TodoItemIDCounter++

	item.CreatedAt = time.Now()
	item.UpdatedAt = time.Now()

	mockdb.TodoItems[item.ID] = item

	// Aynı zamanda ilgili listeye de ekleyelim
	if list, exists := mockdb.TodoLists[item.ListID]; exists {
		list.Items = append(list.Items, item)
	}

	return item, nil
}

func UpdateItem(itemID int, updated *models.TodoItem) (*models.TodoItem, error) {
	item, exists := mockdb.TodoItems[itemID]
	if !exists || item.DeletedAt != nil {
		return nil, errors.New("item not found")
	}

	item.Content = updated.Content
	item.IsDone = updated.IsDone
	item.UpdatedAt = time.Now()

	return item, nil
}

func DeleteItem(itemID int) error {
	item, exists := mockdb.TodoItems[itemID]
	if !exists || item.DeletedAt != nil {
		return errors.New("item not found")
	}

	now := time.Now()
	item.DeletedAt = &now
	item.UpdatedAt = now

	mockdb.TodoItems[itemID] = item

	if list, exists := mockdb.TodoLists[item.ListID]; exists {
		for i, listItem := range list.Items {
			if listItem.ID == itemID {
				list.Items[i].DeletedAt = &now
				list.Items[i].UpdatedAt = now
				break
			}
		}
	}

	return nil
}
func GetItemsByListID(listID int, includeDeleted bool) ([]*models.TodoItem, error) {
	var result []*models.TodoItem

	if _, exists := mockdb.TodoLists[listID]; !exists {
		return nil, errors.New("list not found")
	}

	for _, item := range mockdb.TodoItems {
		if item.ListID == listID {
			if item.DeletedAt != nil && !includeDeleted {
				continue
			}
			result = append(result, item)
		}
	}
	return result, nil
}

func GetItemByID(itemID int) (*models.TodoItem, error) {
	item, exists := mockdb.TodoItems[itemID]
	if !exists || item.DeletedAt != nil {
		return nil, errors.New("item not found")
	}
	return item, nil
}

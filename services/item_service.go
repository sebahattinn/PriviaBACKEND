package services

import (
	"errors"
	"priviatodolist/models"
	"priviatodolist/repositories"
)

func AddItemToList(listID int, userID int, item *models.TodoItem) (*models.TodoItem, error) {
	if !ownsList(userID, listID) {
		return nil, errors.New("unauthorized: list does not belong to user")
	}
	item.ListID = listID
	return repositories.CreateItem(item)
}

func UpdateItem(itemID int, userID int, updatedItem *models.TodoItem) (*models.TodoItem, error) {
	item, err := repositories.GetItemByID(itemID)
	if err != nil {
		return nil, err
	}
	if !ownsList(userID, item.ListID) {
		return nil, errors.New("unauthorized")
	}
	return repositories.UpdateItem(itemID, updatedItem)
}

func DeleteItem(itemID int, userID int) error {
	item, err := repositories.GetItemByID(itemID)
	if err != nil {
		return err
	}
	if !ownsList(userID, item.ListID) {
		return errors.New("unauthorized")
	}
	return repositories.DeleteItem(itemID)
}

func GetItems(listID int, userID int) ([]*models.TodoItem, error) {

	list, err := repositories.GetTodoListByID(listID)
	if err != nil {
		return nil, errors.New("list not found")
	}

	if list.UserID != userID {
		return nil, errors.New("forbidden")
	}

	return repositories.GetItemsByListID(listID, false)
}

func GetAllItemsForAdmin(listID int) ([]*models.TodoItem, error) {
	_, err := repositories.GetItemsByListID(listID, true)
	if err != nil {
		return nil, errors.New("list not found")
	}
	items, err := repositories.GetItemsByListID(listID, true)
	if err != nil {
		return nil, err
	}

	return items, nil
}

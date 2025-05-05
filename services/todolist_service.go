package services

import (
	"errors"
	"priviatodolist/models"
	"priviatodolist/repositories"
	"time"
)

// Yardımcı fonksiyon: Kullanıcı verilen listeye sahip mi?
func ownsList(userID, listID int) bool {
	list, err := repositories.GetTodoListByID(listID)
	if err != nil {
		return false
	}
	return list.UserID == userID
}

// Yeni bir todo listesi ekle
func CreateTodoList(userID int, newList *models.TodoList) (*models.TodoList, error) {
	if len(newList.Name) < 3 {
		return nil, errors.New("title must be at least 3 characters")
	}

	// Kullanıcının listesine ait olmayan bir liste yaratılmasını engelle
	newList.UserID = userID
	newList.CreatedAt = time.Now()
	newList.UpdatedAt = time.Now()

	// Listeni kaydet
	createdList, err := repositories.CreateTodoList(newList)
	if err != nil {
		return nil, err
	}

	// Completion oranını hesapla
	CalculateListCompletion(createdList)

	return createdList, nil
}

// Bir todo listesi güncelle
func UpdateTodoList(listID int, userID int, updatedList *models.TodoList) (*models.TodoList, error) {
	list, err := repositories.GetTodoListByID(listID)
	if err != nil {
		return nil, errors.New("list not found")
	}

	// Kullanıcının sadece kendi listelerini güncellemesine izin veriyoruz
	if !ownsList(userID, listID) {
		return nil, errors.New("unauthorized")
	}

	// Güncelleme işlemi
	list.Name = updatedList.Name
	list.UpdatedAt = time.Now()

	// Completion oranını tekrar hesapla
	CalculateListCompletion(list)

	updatedList, err = repositories.UpdateTodoList(listID, list)
	if err != nil {
		return nil, err
	}

	return updatedList, nil
}

// Todo listesini sil (soft delete)
func DeleteTodoList(listID int, userID int) error {
	list, err := repositories.GetTodoListByID(listID)
	if err != nil {
		return errors.New("list not found")
	}

	// Kullanıcının sadece kendi listelerini silmesine izin veriyoruz
	if !ownsList(userID, listID) {
		return errors.New("unauthorized")
	}

	// Silme işlemi
	now := time.Now()
	list.DeletedAt = &now
	list.UpdatedAt = now

	// Listedeki tüm item'ları da sil
	for _, item := range list.Items {
		if item.DeletedAt == nil { // Zaten silinmemiş item'ları sil
			item.DeletedAt = &now
			item.UpdatedAt = now
			_, err := repositories.UpdateItem(item.ID, item)
			if err != nil {
				return err
			}
		}
	}

	// Listemi repository'e kaydet
	_, err = repositories.UpdateTodoList(listID, list)
	return err
}

// Kullanıcıya ait tüm aktif todo listelerini getir
func GetMyTodoLists(userID int) ([]*models.TodoList, error) {
	lists, err := repositories.GetTodoListsByUserID(userID, false)
	if err != nil {
		return nil, err
	}

	// Her bir liste için completion oranını hesapla
	for _, list := range lists {
		CalculateListCompletion(list)
	}

	return lists, nil
}

// Admin için: Silinmiş dahil tüm todo listelerini getir
func GetAllTodoListsForAdmin() ([]*models.TodoList, error) {
	lists, err := repositories.GetAllTodoLists(true)
	if err != nil {
		return nil, err
	}

	// Her bir liste için completion oranını hesapla
	for _, list := range lists {
		CalculateListCompletion(list)
	}

	return lists, nil
}

// Todo listesi için tamamlanma oranını hesapla
func CalculateListCompletion(list *models.TodoList) {
	total := len(list.Items)
	if total == 0 {
		list.Completion = 0
		return
	}

	doneCount := 0
	for _, item := range list.Items {
		if item.IsDone {
			doneCount++
		}
	}
	list.Completion = float32(doneCount) / float32(total) * 100
}

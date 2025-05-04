package controllers

import (
	"net/http"
	"priviatodolist/mockdb"
	"priviatodolist/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// CreateTodoList godoc
// @Summary Create a new todo list
// @Description Create a new todo list with a title
// @Tags TodoLists
// @Accept json
// @Produce json
// @Param todo body models.TodoList true "Todo List"
// @Success 201 {object} models.TodoList
// @Failure 400 {object} map[string]interface{}
// @Router /todolists [post]
func CreateTodoList(c *gin.Context) {
	var newList models.TodoList
	if err := c.ShouldBindJSON(&newList); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title is required and must be at least 3 characters"})
		return
	}
	if len(newList.Name) < 3 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title must be at least 3 characters"})
		return
	}

	// Token'den UserID alınır
	UserID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	newList.UserID = UserID.(int) // token'den gelen userID burada UserID olarak set edilir
	newList.ID = mockdb.TodoListIDCounter
	newList.CreatedAt = time.Now()
	newList.UpdatedAt = time.Now()

	// Completion hesapla
	CalculateCompletion(&newList)

	mockdb.TodoLists[newList.ID] = &newList
	mockdb.TodoListIDCounter++

	c.JSON(http.StatusCreated, newList)
}

// GetTodoListsForAdmin godoc
// @Summary Get all todo lists
// @Description Retrieve all todo lists
// @Tags TodoLists
// @Produce json
// @Success 200 {array} models.TodoList
// @Router /todolists [get]
func GetTodoListsForAdmin(c *gin.Context) {
	var lists []models.TodoList
	for _, list := range mockdb.TodoLists {
		lists = append(lists, *list)
	}
	c.JSON(http.StatusOK, lists)
}

// UpdateTodoList godoc
// @Summary Update a todo list
// @Description Update a todo list by ID
// @Tags TodoLists
// @Accept json
// @Produce json
// @Param id path int true "Todo List ID"
// @Param todo body models.TodoList true "Updated Todo List"
// @Success 200 {object} models.TodoList
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /todolists/{id} [put]
func UpdateTodoList(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	list, exists := mockdb.TodoLists[id]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "List not found"})
		return
	}

	// Kullanıcının sadece kendi listelerini güncellemesine izin veriyoruz
	userID, _ := c.Get("userID")
	if list.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not allowed to update this list"})
		return
	}

	if err := c.ShouldBindJSON(&list); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	list.UpdatedAt = time.Now()

	// Completion hesapla
	CalculateCompletion(list)

	mockdb.TodoLists[id] = list
	c.JSON(http.StatusOK, list)
}

// DeleteTodoList godoc
// @Summary Soft delete a todo list
// @Description Soft delete a todo list by setting DeletedAt
// @Tags TodoLists
// @Produce json
// @Param id path int true "Todo List ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /todolists/{id} [delete]
func DeleteTodoList(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	list, exists := mockdb.TodoLists[id]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "List not found"})
		return
	}

	// Kullanıcının sadece kendi listelerini silmesine izin veriyoruz
	userID, _ := c.Get("userID")
	if list.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not allowed to delete this list"})
		return
	}

	now := time.Now()
	list.DeletedAt = &now
	list.UpdatedAt = now

	mockdb.TodoLists[id] = list
	c.JSON(http.StatusOK, gin.H{"message": "List marked as deleted"})
}

// AddItemToList godoc
// @Summary Add item to a todo list
// @Description Add a new item to a specific todo list
// @Tags TodoItems
// @Accept json
// @Produce json
// @Param id path int true "Todo List ID"
// @Param item body models.TodoItem true "New Item"
// @Success 201 {object} models.TodoItem
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /todolists/{id}/items [post]
func AddItemToList(c *gin.Context) {
	listID, _ := strconv.Atoi(c.Param("id"))
	list, exists := mockdb.TodoLists[listID]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "List not found"})
		return
	}

	// Kullanıcının sadece kendi listesine item eklemesine izin veriyoruz
	userID, _ := c.Get("userID")
	if list.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not allowed to add items to this list"})
		return
	}

	var newItem models.TodoItem
	if err := c.ShouldBindJSON(&newItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}
	newItem.ID = mockdb.TodoItemIDCounter
	newItem.ListID = listID
	newItem.CreatedAt = time.Now()
	newItem.UpdatedAt = time.Now()

	mockdb.TodoItemIDCounter++
	mockdb.TodoItems[newItem.ID] = &newItem
	list.Items = append(list.Items, &newItem)
	list.UpdatedAt = time.Now()

	// Completion hesapla
	CalculateCompletion(list)

	mockdb.TodoLists[listID] = list

	c.JSON(http.StatusCreated, newItem)
}

// UpdateItem godoc
// @Summary Update a todo item
// @Description Update an item in the list
// @Tags TodoItems
// @Accept json
// @Produce json
// @Param id path int true "Item ID"
// @Param item body models.TodoItem true "Updated Item"
// @Success 200 {object} models.TodoItem
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /items/{id} [put]
func UpdateItem(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	item, exists := mockdb.TodoItems[id]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	}

	// Kullanıcının sadece kendi itemını güncellemesine izin veriyoruz
	list, exists := mockdb.TodoLists[item.ListID]
	if !exists || list.UserID != c.MustGet("userID").(int) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not allowed to update this item"})
		return
	}

	// Gelen JSON verisini item'a bind et
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	// item'ı mockdb.TodoItems'a kaydet
	item.UpdatedAt = time.Now()
	mockdb.TodoItems[id] = item // Burada &item yerine item kullanılıyor

	// Eğer gerekiyorsa, TodoList'i güncelle
	for i, todoItem := range list.Items {
		if todoItem.ID == item.ID {
			list.Items[i] = item // item'ı listeye kaydet
			break
		}
	}

	// Listeyi de güncellemeyi unutma
	list.UpdatedAt = time.Now()

	// Completion hesapla
	CalculateCompletion(list)

	mockdb.TodoLists[item.ListID] = list

	c.JSON(http.StatusOK, item)
}

// DeleteItem godoc
// @Summary Soft delete an item
// @Description Mark an item as deleted by setting DeletedAt
// @Tags TodoItems
// @Produce json
// @Param id path int true "Item ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /items/{id} [delete]
func DeleteItem(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	item, exists := mockdb.TodoItems[id]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	}

	// Kullanıcının sadece kendi itemını silmesine izin veriyoruz
	list, exists := mockdb.TodoLists[item.ListID]
	if !exists || list.UserID != c.MustGet("userID").(int) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not allowed to delete this item"})
		return
	}

	// Daha önce silinmiş bir item tekrar silinemez
	if item.DeletedAt != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Item is already deleted"})
		return
	}

	now := time.Now()
	item.DeletedAt = &now
	item.UpdatedAt = now
	mockdb.TodoItems[id] = item

	// Listeyi güncelle
	list.UpdatedAt = now
	CalculateCompletion(list)
	mockdb.TodoLists[item.ListID] = list

	c.JSON(http.StatusOK, gin.H{"message": "Item marked as deleted"})
}

// GetItems godoc
// @Summary Get items for a specific todo list
// @Description Retrieve all items for a specific todo list
// @Tags TodoItems
// @Produce json
// @Param id path int true "Todo List ID"
// @Success 200 {array} models.TodoItem
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /todolists/{id}/items [get]
func GetItems(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	item, exists := mockdb.TodoItems[id]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	}

	// Kullanıcının sadece kendi itemını güncellemesine izin veriyoruz
	list, exists := mockdb.TodoLists[item.ListID]
	if !exists || list.UserID != c.MustGet("userID").(int) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not allowed to update this item"})
		return
	}

	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}
	item.UpdatedAt = time.Now()
	c.JSON(http.StatusOK, item)
}

// GetMyTodoLists godoc
// @Summary Get all todo lists for the authenticated user
// @Description Retrieve all todo lists for the authenticated user
// @Tags TodoLists
// @Produce json
// @Success 200 {array} models.TodoList
// @Failure 401 {object} map[string]interface{}
// @Router /todolists [get]
func GetMyTodoLists(c *gin.Context) {
	var lists []models.TodoList
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Kullanıcının sadece kendi todo listelerini alması
	for _, list := range mockdb.TodoLists {
		if list.UserID == userID {
			lists = append(lists, *list)
		}
	}

	if len(lists) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No todo lists found"})
		return
	}

	c.JSON(http.StatusOK, lists)
}

func CalculateCompletion(list *models.TodoList) {
	total := len(list.Items)
	if total == 0 {
		// Eğer hiç item yoksa, completion %0 olsun
		list.Completion = 0
		return
	}

	doneCount := 0
	for _, item := range list.Items {
		if item.IsDone {
			doneCount++
		}
	}

	// Tamamlanan item'ların oranını hesapla
	// Eğer doneCount = 0 ise, completion %0 olacak
	// Eğer doneCount = total ise, completion %100 olacak
	list.Completion = float32(doneCount) / float32(total) * 100
}

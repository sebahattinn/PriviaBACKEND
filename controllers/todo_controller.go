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
	newList.ID = mockdb.TodoListIDCounter
	newList.CreatedAt = time.Now()
	newList.UpdatedAt = time.Now()

	mockdb.TodoLists[newList.ID] = &newList
	mockdb.TodoListIDCounter++

	c.JSON(http.StatusCreated, newList)
}

// GetTodoLists godoc
// @Summary Get all todo lists
// @Description Retrieve all todo lists
// @Tags TodoLists
// @Produce json
// @Success 200 {array} models.TodoList
// @Router /todolists [get]
func GetTodoLists(c *gin.Context) {
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
	if err := c.ShouldBindJSON(&list); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}
	list.UpdatedAt = time.Now()
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
	now := time.Now()
	list.DeletedAt = &now
	list.UpdatedAt = now
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
	list.Items = append(list.Items, newItem)
	list.UpdatedAt = time.Now()

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

	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}
	item.UpdatedAt = time.Now()
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
	now := time.Now()
	item.DeletedAt = &now
	item.UpdatedAt = now
	c.JSON(http.StatusOK, gin.H{"message": "Item marked as deleted"})
}

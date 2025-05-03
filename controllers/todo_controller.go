package controllers

import (
	"net/http"
	"priviatodolist/mockdb"
	"priviatodolist/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func CreateTodoList(c *gin.Context) {
	var newList models.TodoList
	if err := c.ShouldBindJSON(&newList); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}
	newList.ID = mockdb.TodoListIDCounter
	newList.CreatedAt = time.Now()
	newList.UpdatedAt = time.Now()
	mockdb.TodoLists[newList.ID] = &newList
	mockdb.TodoListIDCounter++

	c.JSON(http.StatusCreated, newList)
}

func GetTodoLists(c *gin.Context) {
	var lists []models.TodoList
	for _, list := range mockdb.TodoLists {
		lists = append(lists, *list)
	}
	c.JSON(http.StatusOK, lists)
}

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

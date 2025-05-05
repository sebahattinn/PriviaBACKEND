package controllers

import (
	"net/http"
	"priviatodolist/models"
	"priviatodolist/services"
	"priviatodolist/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

// GetTodoItems godoc
// @Summary Get active items for a specific list
// @Description Retrieves non-deleted todo items for the specified list. Users can only access their own lists.
// @Tags TodoItems
// @Security BearerAuth
// @Produce json
// @Param id path int true "Todo List ID"
// @Success 200 {array} models.TodoItem "List of active todo items"
// @Failure 403 "You don't have permission to access this list"
// @Failure 404 "List not found"
// @Router /todolists/{id}/items [get]
func GetTodoItems(c *gin.Context) {
	listID, error := getIDParam(c)
	if error != nil {
		utils.HandleError(c, http.StatusBadRequest, error, "Invalid Todo List ID")
		return
	}
	userID, exists := getUserID(c)
	if !exists {
		utils.HandleError(c, http.StatusUnauthorized, nil, "User not authorized")
		return
	}
	items, err := services.GetItems(listID, userID)
	if err != nil {
		if err.Error() == "forbidden" {
			utils.HandleError(c, http.StatusForbidden, err, "You don't have permission to access this list")
			return
		} else if err.Error() == "list not found" {
			utils.HandleError(c, http.StatusNotFound, err, "Todo list not found")
			return
		}
		utils.HandleError(c, http.StatusInternalServerError, err, "Failed to retrieve items")
		return
	}
	c.JSON(http.StatusOK, items)
}

// AddTodoItem godoc
// @Summary Add a new todo item
// @Description Adds a new item to a specific todo list. Requires authenticated user access.
// @Tags TodoItems
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Todo List ID"
// @Param item body models.TodoItemCreate true "New Item Details"
// @Success 201 {object} models.TodoItem "Newly created todo item"
// @Failure 400 "Invalid request (e.g., malformed JSON)"
// @Failure 403 "No permission to access this list"
// @Failure 404 "Specified list not found"
// @Router /todolists/{id}/items [post]
func AddTodoItem(c *gin.Context) {
	listID, err := getIDParam(c)
	if err != nil {
		utils.HandleError(c, http.StatusBadRequest, err, "Invalid Todo List ID")
		return
	}
	userID, exists := getUserID(c)
	if !exists {
		utils.HandleError(c, http.StatusUnauthorized, nil, "User not authorized")
		return
	}
	var newItem models.TodoItem
	if err := c.ShouldBindJSON(&newItem); err != nil {
		utils.HandleError(c, http.StatusBadRequest, err, "Invalid request payload")
		return
	}
	item, err := services.AddItemToList(listID, userID, &newItem)
	if err != nil {
		if err.Error() == "unauthorized: list does not belong to user" {
			utils.HandleError(c, http.StatusForbidden, err, "You don't have permission to add items to this list")
			return
		} else if strings.Contains(err.Error(), "not found") {
			utils.HandleError(c, http.StatusNotFound, err, "Todo list not found")
			return
		}
		utils.HandleError(c, http.StatusInternalServerError, err, "Failed to add item")
		return
	}
	c.JSON(http.StatusCreated, item)
}

// UpdateItem godoc
// @Summary Update a todo item
// @Description Update an item in the list
// @Tags TodoItems
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Item ID"
// @Param item body models.TodoItemUpdate true "Updated Item"
// @Success 200 {object} models.TodoItem
// @Failure 400 "Invalid request payload"
// @Failure 404 "Item not found"
// @Router /items/{id} [put]
func UpdateTodoItem(c *gin.Context) {
	itemID, error := getIDParam(c)
	if error != nil {
		utils.HandleError(c, http.StatusBadRequest, error, "Invalid Todo Item ID")
		return
	}
	userID, exists := getUserID(c)
	if !exists {
		utils.HandleError(c, http.StatusUnauthorized, nil, "User not authorized")
		return
	}
	var updatedItem models.TodoItem
	if err := c.ShouldBindJSON(&updatedItem); err != nil {
		utils.HandleError(c, http.StatusBadRequest, err, "Invalid request payload")
		return
	}
	item, err := services.UpdateItem(itemID, userID, &updatedItem)
	if err != nil {
		utils.HandleError(c, http.StatusNotFound, err, "Item not found")
		return
	}
	c.JSON(http.StatusOK, item)
}

// DeleteTodoItem godoc
// @Summary Soft delete a todo item
// @Description Marks a specific todo item as deleted (soft delete). Not actually deleted, just fills the DeletedAt field.
// @Tags TodoItems
// @Security BearerAuth
// @Produce json
// @Param id path int true "Todo Item ID"
// @Success 200 {object} map[string]interface{} "Deletion successful"
// @Failure 400 "Item already deleted"
// @Failure 403 "No permission to delete this item"
// @Failure 404 "Todo item not found"
// @Router /items/{id} [delete]
func DeleteTodoItem(c *gin.Context) {
	itemID, error := getIDParam(c)
	if error != nil {
		utils.HandleError(c, http.StatusBadRequest, error, "Invalid Todo Item ID")
		return
	}
	userID, exists := getUserID(c)
	if !exists {
		utils.HandleError(c, http.StatusUnauthorized, nil, "User not authorized")
		return
	}
	if err := services.DeleteItem(itemID, userID); err != nil {
		utils.HandleError(c, http.StatusNotFound, err, "Item not found or already deleted")
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Item marked as deleted"})
}

// GetAllTodoItemsForAdmin godoc
// @Summary Get all items for a specific list (including deleted ones)
// @Description Retrieves all todo items for the specified list (including deleted ones). Only admin users can access.
// @Tags Admin
// @Security BearerAuth
// @Produce json
// @Param id path int true "Todo List ID"
// @Success 200 {array} models.TodoItem "All items in the list"
// @Failure 403 "Forbidden: Only admin users can access this endpoint"
// @Failure 404 "Todo list not found"
// @Router /admin/todolists/{id}/items [get]
func GetAllTodoItemsForAdmin(c *gin.Context) {
	listID, error := getIDParam(c)
	if error != nil {
		utils.HandleError(c, http.StatusBadRequest, error, "Invalid Todo List ID")
		return
	}
	items, err := services.GetAllItemsForAdmin(listID)
	if err != nil {
		if err.Error() == "list not found" {
			utils.HandleError(c, http.StatusNotFound, err, "Todo list not found")
		} else {
			utils.HandleError(c, http.StatusInternalServerError, err, "Failed to retrieve items")
		}
		return
	}
	c.JSON(http.StatusOK, items)
}

// Helper functions
func CalculateCompletion(list *models.TodoList) {
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

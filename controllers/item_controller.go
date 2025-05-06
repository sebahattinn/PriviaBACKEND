package controllers

import (
	"net/http"
	"priviatodolist/models"
	"priviatodolist/services"
	"priviatodolist/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

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

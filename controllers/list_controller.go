package controllers

import (
	"net/http"
	"priviatodolist/models"
	"priviatodolist/services"
	"priviatodolist/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

func getIDParam(c *gin.Context) (int, error) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return 0, err
	}
	return id, nil
}

func getUserID(c *gin.Context) (int, bool) {
	userID, exists := c.Get("userID")
	if !exists {
		return 0, false
	}
	return userID.(int), true
}

func CreateTodoList(c *gin.Context) {
	var newList models.TodoList
	if err := c.ShouldBindJSON(&newList); err != nil {
		utils.HandleError(c, http.StatusBadRequest, err, "Invalid request payload")
		return
	}
	userID, exists := getUserID(c)
	if !exists {
		utils.HandleError(c, http.StatusUnauthorized, nil, "User not authorized")
		return
	}
	list, err := services.CreateTodoList(userID, &newList)
	if err != nil {
		utils.HandleError(c, http.StatusBadRequest, err, err.Error())
		return
	}
	c.JSON(http.StatusCreated, list)
}

func GetTodoListsForAdmin(c *gin.Context) {
	lists, err := services.GetAllTodoListsForAdmin()
	if err != nil {
		utils.HandleError(c, http.StatusInternalServerError, err, "Failed to retrieve todo lists")
		return
	}
	c.JSON(http.StatusOK, lists)
}

func UpdateTodoList(c *gin.Context) {
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
	var updatedList models.TodoList
	if err := c.ShouldBindJSON(&updatedList); err != nil {
		utils.HandleError(c, http.StatusBadRequest, err, "Invalid request payload")
		return
	}
	list, err := services.UpdateTodoList(listID, userID, &updatedList)
	if err != nil {
		if err.Error() == "unauthorized" {
			utils.HandleError(c, http.StatusForbidden, err, "You are not allowed to update this list")
		} else {
			utils.HandleError(c, http.StatusNotFound, err, "Todo list not found")
		}
		return
	}
	c.JSON(http.StatusOK, list)
}

func DeleteTodoList(c *gin.Context) {
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
	if err := services.DeleteTodoList(listID, userID); err != nil {
		if err.Error() == "unauthorized" {
			utils.HandleError(c, http.StatusForbidden, err, "You are not allowed to delete this list")
		} else {
			utils.HandleError(c, http.StatusNotFound, err, "Todo list not found or already deleted")
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "List and all its items marked as deleted"})
}

func GetMyTodoLists(c *gin.Context) {
	userID, exists := getUserID(c)
	if !exists {
		utils.HandleError(c, http.StatusUnauthorized, nil, "User not authorized")
		return
	}
	lists, err := services.GetMyTodoLists(userID)
	if err != nil {
		utils.HandleError(c, http.StatusInternalServerError, err, "Failed to retrieve todo lists")
		return
	}
	c.JSON(http.StatusOK, lists)
}

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

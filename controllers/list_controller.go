package controllers

import (
	"net/http"
	"priviatodolist/models"
	"priviatodolist/services"
	"priviatodolist/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Helper functions
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

// CreateTodoList godoc
// @Summary Create a new todo list
// @Description Creates a new todo list for the authenticated user
// @Tags TodoLists
// @Accept json
// @Produce json
// @Param todoList body models.TodoListCreate true "Todo List details"
// @Success 201 {object} models.TodoList "Successfully created todo list"
// @Failure 400 "Invalid request payload"
// @Failure 401 "User not authorized"
// @Security BearerAuth
// @Router /todolists [post]
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

// GetTodoListsForAdmin godoc
// @Summary Retrieve all todo lists
// @Description Retrieves all todo lists including deleted ones (admin only)
// @Tags Admin
// @Produce json
// @Success 200 {array} models.TodoList "List of all todo lists"
// @Failure 401 "Unauthorized access"
// @Failure 500 "Internal server error"
// @Security BearerAuth
// @Router /admin/todolists [get]
func GetTodoListsForAdmin(c *gin.Context) {
	lists, err := services.GetAllTodoListsForAdmin()
	if err != nil {
		utils.HandleError(c, http.StatusInternalServerError, err, "Failed to retrieve todo lists")
		return
	}
	c.JSON(http.StatusOK, lists)
}

// UpdateTodoList godoc
// @Summary Update a todo list
// @Description Updates an existing todo list by ID
// @Tags TodoLists
// @Accept json
// @Produce json
// @Param id path int true "Todo List ID"
// @Param todoList body models.TodoListUpdate true "Updated todo list details"
// @Success 200 {object} models.TodoList "Successfully updated todo list"
// @Failure 400 "Invalid request payload"
// @Failure 401 "User not authorized"
// @Failure 403 "Forbidden access"
// @Failure 404 "Todo list not found"
// @Security BearerAuth
// @Router /todolists/{id} [put]
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

// DeleteTodoList godoc
// @Summary Delete a todo list
// @Description Soft deletes a todo list by setting DeletedAt timestamp
// @Tags TodoLists
// @Produce json
// @Param id path int true "Todo List ID"
// @Success 200 "List and all its items marked as deleted"
// @Failure 400 "Invalid Todo List ID"
// @Failure 401 "User not authorized"
// @Failure 403 "Forbidden access"
// @Failure 404 "Todo list not found or already deleted"
// @Security BearerAuth
// @Router /todolists/{id} [delete]
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

// GetMyTodoLists godoc
// @Summary Get user's todo lists
// @Description Retrieves all active todo lists for the authenticated user
// @Tags TodoLists
// @Produce json
// @Success 200 {array} models.TodoList "User's todo lists"
// @Failure 401 "User not authorized"
// @Failure 500 "Failed to retrieve todo lists"
// @Security BearerAuth
// @Router /todolists [get]
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

// CalculateListCompletion calculates the completion percentage of a todo list
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

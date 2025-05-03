package router

import (
	"priviatodolist/controllers"
	"priviatodolist/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Giriş endpoint'i (JWT'siz erişilebilir)
	r.POST("/login", controllers.Login)

	// JWT korumalı API grubu
	api := r.Group("/api")
	api.Use(middleware.JWTAuthMiddleware())

	// Todo listeleri işlemleri
	todo := api.Group("/todo")
	{
		todo.POST("/", controllers.CreateTodoList)
		todo.GET("/", controllers.GetTodoLists)
		todo.PUT("/:id", controllers.UpdateTodoList)
		todo.DELETE("/:id", controllers.DeleteTodoList)

		// Alt kaynak: Todo listesine ait item işlemleri
		todo.POST("/:id/items", controllers.AddItemToList)
		todo.PUT("/:listId/items/:itemId", controllers.UpdateItem)
		todo.DELETE("/:listId/items/:itemId", controllers.DeleteItem)
	}

	return r
}

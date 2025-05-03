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
	api.Use(middleware.JWTAuthMiddleware()) // JWT doğrulama middleware'ı ekle

	// Todo listeleri işlemleri
	todo := api.Group("/todo")
	{
		todo.POST("/", controllers.CreateTodoList)
		todo.GET("/", controllers.GetTodoLists)
		todo.PUT("/:todoId", controllers.UpdateTodoList)    // Parametre adı :id -> :todoId
		todo.DELETE("/:todoId", controllers.DeleteTodoList) // Parametre adı :id -> :todoId

		// Alt kaynak: Todo listesine ait item işlemleri
		todo.POST("/:todoId/items", controllers.AddItemToList)        // Parametre adı :id -> :todoId
		todo.PUT("/:todoId/items/:itemId", controllers.UpdateItem)    // Parametre adı :listId -> :todoId
		todo.DELETE("/:todoId/items/:itemId", controllers.DeleteItem) // Parametre adı :listId -> :todoId
	}

	// Admin sadece erişimi olan endpoint
	admin := api.Group("/admin")
	admin.Use(middleware.AdminOnly()) // Admin yetkilendirme middleware'ı ekle
	{
		admin.GET("/lists", controllers.GetTodoLists)
	}

	return r
}

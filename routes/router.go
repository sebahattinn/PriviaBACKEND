package routes

import (
	"priviatodolist/controllers"
	"priviatodolist/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API v1 grubu (JWT korumalı)
	api := r.Group("/api/v1")
	api.POST("/login", controllers.Login)
	api.Use(middleware.JWTAuthMiddleware())
	api.Use(middleware.GlobalErrorHandler())
	{
		api.GET("/todolists", controllers.GetMyTodoLists)
		api.POST("/todolists", controllers.CreateTodoList)
		api.GET("/todolists/:id/items", controllers.GetTodoItems)
		api.POST("/todolists/:id/items", controllers.AddTodoItem)
		api.PUT("/items/:id", controllers.UpdateTodoItem)
		api.DELETE("/items/:id", controllers.DeleteTodoItem)
		api.PUT("/todolists/:id", controllers.UpdateTodoList)
		api.DELETE("/todolists/:id", controllers.DeleteTodoList)

		adminOnly := api.Group("/admin")
		adminOnly.Use(middleware.AdminOnly())
		{
			adminOnly.GET("/todolists", controllers.GetTodoListsForAdmin)
			adminOnly.GET("/todolists/:id/items", controllers.GetAllTodoItemsForAdmin)
		}
	}

	return r
}

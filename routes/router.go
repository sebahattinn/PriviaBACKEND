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
	{
		// Kullanıcının kendi todolist işlemleri endpointlerin başına /api/v1 unutulmamalı.
		api.GET("/todolists", controllers.GetMyTodoLists)
		api.POST("/todolists", controllers.CreateTodoList)
		api.GET("/todolists/:id/items", controllers.GetItems)
		api.POST("/todolists/:id/items", controllers.AddItemToList)
		api.PUT("/items/:id", controllers.UpdateItem)
		api.DELETE("/items/:id", controllers.DeleteItem)
		api.PUT("/todolists/:id", controllers.UpdateTodoList)
		api.DELETE("/todolists/:id", controllers.DeleteTodoList)

		// Sadece admin olanlar tüm todolistleri görebilir
		adminOnly := api.Group("/admin")
		adminOnly.Use(middleware.AdminOnly())
		{
			adminOnly.GET("/todolists", controllers.GetTodoListsForAdmin)
		}
	}

	return r
}

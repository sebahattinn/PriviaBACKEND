package main

import (
	"priviatodolist/controllers"
	"priviatodolist/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Privia Todo List API
// @version         1.0
// @description     Bu API, yapılacaklar listelerini ve içindeki görevleri yönetmek için geliştirilmiştir.
// @termsOfService  http://swagger.io/terms/

// @contact.name   Sebahattin - Developer
// @contact.email  sebahattin@example.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api/v1
func main() {
	r := gin.Default()

	// Swagger bilgilerini ayarla
	docs.SwaggerInfo.BasePath = "/api/v1"

	// API Routes
	api := r.Group("/api/v1")
	{
		api.POST("/todolists", controllers.CreateTodoList)
		api.GET("/todolists", controllers.GetTodoLists)
		api.PUT("/todolists/:id", controllers.UpdateTodoList)
		api.DELETE("/todolists/:id", controllers.DeleteTodoList)

		api.POST("/todolists/:id/items", controllers.AddItemToList)
		api.PUT("/items/:id", controllers.UpdateItem)
		api.DELETE("/items/:id", controllers.DeleteItem)
	}

	// Swagger route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Sunucuyu başlat
	r.Run(":8080") // localhost:8080'de çalışır
}

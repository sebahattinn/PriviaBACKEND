package main

import (
	"priviatodolist/controllers"
	"priviatodolist/docs"
	"priviatodolist/middleware"

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

	// JWT doğrulama middleware'ı
	r.Use(middleware.JWTAuthMiddleware()) // Tüm API'lerde JWT doğrulaması yapılacak

	// API Routes
	api := r.Group("/api/v1")
	{
		// Admin kontrolü gereken rotalar
		adminRoutes := api.Group("/")
		adminRoutes.Use(middleware.AdminOnly()) // Sadece admin kullanıcılar erişebilecek
		{
			adminRoutes.POST("/todolists", controllers.CreateTodoList)
			adminRoutes.GET("/todolists", controllers.GetTodoLists)
			adminRoutes.PUT("/todolists/:id", controllers.UpdateTodoList)
			adminRoutes.DELETE("/todolists/:id", controllers.DeleteTodoList)
			adminRoutes.POST("/todolists/:id/items", controllers.AddItemToList)
			adminRoutes.PUT("/items/:id", controllers.UpdateItem)
			adminRoutes.DELETE("/items/:id", controllers.DeleteItem)
		}

		// Diğer kullanıcılar için erişim sağlanabilir API
		// Buradaki GET /todolists'i kaldırdık çünkü adminRoutes içinde zaten var
		api.GET("/todolists/:id/items", controllers.GetItems)
	}

	// Swagger route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Sunucuyu başlat
	r.Run(":8080") // localhost:8080'de çalışır
}

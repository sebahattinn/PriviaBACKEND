package router

import (
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.POST("/login", controllers.Login)

	authorized := r.Group("/api", middleware.AuthMiddleware())
	{
		authorized.POST("/todos", controllers.CreateTodo)
		authorized.GET("/todos", controllers.GetTodos)
		// Diğer CRUD işlemleri burada
	}

	return r
}

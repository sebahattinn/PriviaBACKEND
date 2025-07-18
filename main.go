package main

import (
	"priviatodolist/docs"
	"priviatodolist/routes"
)

// @title           Privia Todo List API
// @version         1.0
// @description     Bu API, yapılacaklar listelerini ve içindeki görevleri yönetmek için geliştirilmiştir.
// @termsOfService  http://swagger.io/terms/

// @contact.name   Sebahattin - Developer
// @contact.email  sebahattin@example.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8081
// @BasePath  /api/v1

// @SecurityDefinitions.apiKey BearerAuth
// @SecuritySchemeBearer
// @description Bearer Authentication kullanarak erişim
// @name Authorization
// @in header
// @type apiKey

func main() {
	docs.SwaggerInfo.BasePath = "/api/v1"
	r := routes.SetupRouter()
	r.Run(":8081")
}

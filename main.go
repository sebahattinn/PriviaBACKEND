package main

import (
	"priviatodolist/router"
)

func main() {
	r := router.SetupRouter()
	r.Run(":8080") // localhost:8080
}

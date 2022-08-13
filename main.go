package main

import (
	routes "goferm-gin-example/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.New()

	router.Use(gin.Logger())
	routes.TaskRouter(router)

	port := ":" + "8001"
	router.Run(port)
}

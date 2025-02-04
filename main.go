package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"users/src/routes"
)

func main() {
	mainRouter := gin.Default()
	replicaRouter := gin.Default()

	routes.mainProductsRoutes(mainRouter)
	routes.repliProductosRoutes(repliRouter)

	go func() {
		if err := mainRouter.Run(":8080"); err != nil {
			log.Fatalf("Error al iniciar el servidor principal: %v", err)
		}
	}()

	go func() {
		if err := replicaRouter.Run(":8081"); err != nil {
			log.Fatalf("Error al iniciar el servidor de replicación: %v", err)
		}
	}()

	select {}
}

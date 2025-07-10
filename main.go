package main

import (
	"os"
	"stock-dashboard/config"
	"stock-dashboard/db"
	"stock-dashboard/middleware"
	"stock-dashboard/routes"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	config.LoadEnv()

	db.Connect()

	server := gin.Default()

	server.Use(middleware.CorsMiddleware())

	routes.RegisterRoutes(server)

	port := os.Getenv("PORT")
	server.Run(":" + port)
}

package main

import (
	"fmt"
	"github.com/Dhanrajsinh002/go-url-shortener/routes"
	"github.com/Dhanrajsinh002/go-url-shortener/store"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.Use(cors.Default())

	// Register the routes to use them
	routes.UrlRoutes(r)

	// Initialize redis for storage
	store.InitializeStore()

	err := r.Run(":8000")
	if err != nil {
		panic(fmt.Sprintf("Failed to start server: %v", err))
	}
}
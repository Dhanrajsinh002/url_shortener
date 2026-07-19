package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Dhanrajsinh002/go-url-shortener/routes"
	"github.com/Dhanrajsinh002/go-url-shortener/store"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	databaseURL := os.Getenv("DATABASE_URL")
	redisAddr := os.Getenv("REDIS_ADDR")
	if databaseURL == "" {
		panic("DATABASE_URL is not set")
	}

	if redisAddr == "" {
		panic("REDIS_ADDR is not set")
	}

	r := gin.Default()
	r.Use(cors.Default())

	// Register the routes to use them
	routes.UrlRoutes(r)

	// Initialize redis for storage
	store.InitializeStore(databaseURL, redisAddr)

	err := r.Run(":8000")
	if err != nil {
		panic(fmt.Sprintf("Failed to start server: %v", err))
	}
}
package routes

import (
	"time"

	"github.com/Dhanrajsinh002/go-url-shortener/handler"
	"github.com/Dhanrajsinh002/go-url-shortener/middleware"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

func UrlRoutes(r *gin.Engine, redisClient *redis.Client) {
	r.GET("/", GetHome)
	r.POST("/create-short-url", CreateShortUrl)
	r.GET("/:shortUrl", GetShortUrl)
	r.GET("/stats/:shortUrl", GetStats)

	r.POST("/admin/register", 
		middleware.RateLimit(redisClient, "register", 5, time.Hour),
		func(c *gin.Context) { handler.RegisterUser(c) })

	r.POST("/admin/login", 
		middleware.RateLimit(redisClient, "login", 10, 15*time.Minute),
		func(c *gin.Context) { handler.Login(c) })

	adminGroup := r.Group("/admin")
	adminGroup.Use(middleware.AuthRequired())
	adminGroup.GET("/urls", func(c *gin.Context) { handler.ListUrls(c) })
	adminGroup.POST("/create-short-url", func(c *gin.Context) {handler.CreateUserShortUrl(c)})
}

func GetHome(c *gin.Context) {
	message := "Hey Go Url Shortener! with route OOPS!"
	c.JSON(200, message)
}

func CreateShortUrl(c *gin.Context) {
	handler.CreateShourtUrl(c)
}

func GetShortUrl(c *gin.Context) {
	handler.HandleShortUrlRedirect(c)
}

func GetStats(c *gin.Context) {
	handler.GetUrlStats(c)
}
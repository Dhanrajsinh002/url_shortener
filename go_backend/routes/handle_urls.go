package routes

import (
	"github.com/Dhanrajsinh002/go-url-shortener/handler"
	"github.com/gin-gonic/gin"
)

func UrlRoutes(r *gin.Engine) {
	r.GET("/", GetHome)
	r.POST("/create-short-url", CreateShortUrl)
	r.GET("/:shortUrl", GetShortUrl)
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
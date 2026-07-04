package handler

import (
	"net/http"
	"github.com/Dhanrajsinh002/go-url-shortener/shortener"
	"github.com/Dhanrajsinh002/go-url-shortener/store"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

// Request model definition
type UrlCreationReuqest struct {
	LongUrl string `json:"long_url" binding:"required"`
}

func CreateShourtUrl(c *gin.Context) {
	var creationReuqest UrlCreationReuqest
	if err := c.ShouldBindJSON(&creationReuqest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error" : err.Error()})
		return
	}

	shortUrl := shortener.GenerateShortLink(creationReuqest.LongUrl)
	if err := store.SaveUrlMapping(shortUrl, creationReuqest.LongUrl); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not save url"})
		return
	}

	host := "http://localhost:8000/"
	c.JSON(200, gin.H{
		"message" : "short url created successfully",
		"short_url" : host + shortUrl,
	})
}

func HandleShortUrlRedirect(c *gin.Context) {
	shortUrl := c.Param("shortUrl")

	initialUrl, err := store.RetrieveInitialUrl(shortUrl)
	if err == redis.Nil {
		// Key was not in Redis -> this short code does not exist.
		c.JSON(http.StatusNotFound, gin.H{"error": ""})
		return
	}
	if err != nil {
		// Any other error means Redis itself failed (e.g. down).
		c.JSON(http.StatusInternalServerError, gin.H{"error": ""})
		return
	}

	c.Redirect(http.StatusFound, initialUrl)
}
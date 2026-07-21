package handler

import (
	"net/http"
	"os"

	"github.com/Dhanrajsinh002/go-url-shortener/auth"
	"github.com/Dhanrajsinh002/go-url-shortener/shortener"
	"github.com/Dhanrajsinh002/go-url-shortener/store"
	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": formatValidationError(err)})
		return
	}

	user, err := store.GetUserByUsername(req.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
		return
	}

	if !auth.CheckPasswordHash(req.Password, user.PasswordHash) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
		return
	}

	token, err := auth.GenerateJWT(user.Username, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func ListUrls(c *gin.Context) {
	userID := c.MustGet("user_id").(int64)

	records, err := store.ListUrlsByUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch urls"})
		return
	}
	c.JSON(http.StatusOK, records)
}

func RegisterUser(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": formatValidationError(err)})
		return
	}

	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not hash password"})
		return
	}

	if err := store.CreateUser(req.Username, hash); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create account (username may already exist)"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "account created successfully"})
}

type UserCreateRequest struct {
	LongUrl string `json:"long_url" binding:"required"`
}

func CreateUserShortUrl(c *gin.Context) {
	userID := c.MustGet("user_id").(int64)

	var req UserCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": formatValidationError(err)})
		return
	}

	shortUrl := shortener.GenerateShortLink(req.LongUrl)
	if err := store.SaveUrlMappingForUser(shortUrl, req.LongUrl, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not save url"})
		return
	}

	baseUrl := os.Getenv("BASE_URL")
	if baseUrl == "" {
		baseUrl = "http://localhost:8000"
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "short url created successfully",
		"short_url": baseUrl + "/" + shortUrl,
	})
}

// RegisterAdmin is gated behind a setup token, not a public sign-up flow.
// Use it once to create your first admin, then treat ADMIN_SETUP_TOKEN as
// burned — rotate or remove it (see setup steps below).
type RegisterRequest struct {
	Username	string `json:"username" binding:"required,min=3,max=50"`
	Password	string `json:"password" binding:"required,min=8"`
}

func RegisterAdmin(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": formatValidationError(err)})
		return
	}

	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{ "error": "couls not hash password" })
		return
	}

	if err := store.CreateUser(req.Username, hash); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{ "error": "could not create admin user" })
		return
	}

	c.JSON(http.StatusCreated, gin.H{ "message": "admin created successfully" })
}
package store

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

var testStoreService = &StorageService{}

func init() {
	_ = godotenv.Load("../.env") // adjust path relative to the store package
	databaseURL := os.Getenv("DATABASE_URL")
	redisAddr := os.Getenv("REDIS_ADDR")
	
	if databaseURL == "" {
		panic("DATABASE_URL is not set")
	}
	if redisAddr == "" {
		panic("REDIS_ADDR is not set")
	}
	testStoreService = InitializeStore(databaseURL, redisAddr)
}

func TestStoreInit(t *testing.T) {
	assert.True(t, testStoreService.redisClient != nil)
}

func TestInsertionAndRetrieval(t *testing.T) {
	initialLink := "https://www.guru3d.com/news-story/spotted-ryzen-threadripper-pro-3995wx-with-8-channel-ddr4,2.html"
	shortURL := "Jsz4k57oAX"

	// Persist data mapping
	saveErr := SaveUrlMapping(shortURL, initialLink)
	assert.NoError(t, saveErr)

	// Retrieve initial url
	retrievedUrl, err := RetrieveInitialUrl(shortURL)

	assert.NoError(t, err)
	assert.Equal(t, initialLink, retrievedUrl)
}

func TestRetrieveNotFound(t *testing.T) {
	// A short code that was never saved should come back as ErrNotFound,
	// which is how the handler decides to return a 404. Previously this
	// checked redis.Nil directly, but a miss can now originate from either
	// Redis or Postgres, so the sentinel error abstracts that away.
	_, err := RetrieveInitialUrl("this-code-does-not-exist")

	assert.Equal(t, ErrNotFound, err)
}
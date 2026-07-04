package store

import (
	"github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"
	"testing"
)

var testStoreService = &StorageService{}

func init() {
	testStoreService = InitializeStore()
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
	// A short code that was never saved should come back as redis.Nil,
	// which is how the handler decides to return a 404.
	_, err := RetrieveInitialUrl("this-code-does-not-exist")

	assert.Equal(t, redis.Nil, err)
}
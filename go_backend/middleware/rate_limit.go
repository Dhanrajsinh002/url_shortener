package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

// RateLimit allows `limit` requests per `window` per client IP, keyed by
// route. Uses Redis INCR + EXPIRE so it works correctly even with
// multiple backend replicas sharing one Redis instance.
func RateLimit(redisClient *redis.Client, keyPrefix string, limit int64, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := fmt.Sprintf("ratelimit:%s:%s", keyPrefix, c.ClientIP())

		count, err := redisClient.Incr(key).Result()
		if err != nil {
			// Fail open - a Redis hookup should not block real users from logging in.
			c.Next()
			return
		}

		if count == 1 {
			redisClient.Expire(key, window)
		}

		if count > limit {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "too many attemps, please try again later",
			})
			return
		}
		c.Next()
	}
}
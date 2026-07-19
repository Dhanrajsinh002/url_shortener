package store

import (
	"context"
	"fmt"
	"time"
	"errors"

	"github.com/go-redis/redis"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Define the struct wrapper around raw Redis Client
type StorageService struct {
	redisClient *redis.Client
}

// Top level declaration for storeService, Redis and Postgres context 
var (
	storeService 	= &StorageService{}
	pgPool			*pgxpool.Pool
	ctx 			= context.Background()
)

const CacheDuration = 6 * time.Hour

// ErrNotFound is returned when a short code exists in neither Redis nor
// Postgres. Handlers check for this instead of redis.Nil now, since a
// miss can originate from either backend.
var ErrNotFound = errors.New("short url not found")

func InitializeStore(databaseurl string, redisAddr string) *StorageService {

	redisClient := redis.NewClient(&redis.Options{
		Addr		: redisAddr,
		Password	: "",
		DB			: 0,
	})

	pong, err := redisClient.Ping().Result()
	if err != nil {
		panic(fmt.Sprintf("Error init Redis: %v", err))
	}

	fmt.Printf("\nRedis started succesfully: pong mesg = {%s}", pong)
	storeService.redisClient = redisClient

	pool, err := pgxpool.New(ctx, databaseurl)
	if err != nil {
		panic(fmt.Sprintf("Error init PostgreSQL: %v", err))
	}
	if err := pool.Ping(ctx); err != nil {
		panic(fmt.Sprintf("Error ping PostgreSQL: %v", err))
	}

	fmt.Printf("\nPostgreSQL started succesfully: pool = {%v}", pool)
	pgPool = pool

	return storeService
}

// SaveUrlMapping writes to Postgres first (source of truth). If that
// succeeds, it also writes to Redis as a cache — a cache write failure is
// logged but does not fail the request, since Postgres already has it.
func SaveUrlMapping(shortUrl string, originalUrl string) error {

	_, pgErr := pgPool.Exec(ctx, `INSERT INTO urls (short_code, long_url) VALUES ($1, $2)`, shortUrl, originalUrl)
	if pgErr != nil {
		return pgErr
	}

	if rdsErr := storeService.redisClient.Set(shortUrl, originalUrl, CacheDuration).Err(); rdsErr != nil {
		fmt.Printf("Error writing to Redis cache: %v", rdsErr)
	}

	return nil
}

// RetrieveInitialUrl checks Redis first. On a cache miss it falls back to
// Postgres and backfills Redis for next time. Returns ErrNotFound only
// when the code exists in neither store.
func RetrieveInitialUrl(shortUrl string) (string, error) {

	val, err := storeService.redisClient.Get(shortUrl).Result()
	if err == nil {
		_ = IncrementClickCount(shortUrl)
		return val, nil
	}

	if err != redis.Nil {
		fmt.Printf("\nWarning reading from Redis cache failed for %s: %v", shortUrl, err)
	}

	var longUrl string
	err = pgPool.QueryRow(ctx, `SELECT long_url FROM urls WHERE short_code = $1`, shortUrl).Scan(&longUrl)
	if err != nil {
		return "", ErrNotFound
	}

	if err := storeService.redisClient.Set(shortUrl, longUrl, CacheDuration).Err(); err != nil {
		fmt.Printf("\nWarning writing to Redis cache failed for %s: %v", shortUrl, err)
	}

	_ = IncrementClickCount(shortUrl)
	return longUrl, nil
}

func IncrementClickCount(shortUrl string) error {
	_, err := pgPool.Exec(ctx, `UPDATE urls SET click_count = click_count + 1 WHERE short_code = $1`, shortUrl)
	return err
}

func GetClickCount(shortUrl string) (int64, error) {
	var count int64
	err := pgPool.QueryRow(ctx, `SELECT click_count FROM urls WHERE short_code = $1`, shortUrl).Scan(&count)

	if err != nil {
		return 0, ErrNotFound
	}

	return count, nil
}
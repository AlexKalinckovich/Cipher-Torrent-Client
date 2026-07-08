package config

import (
	"database/sql"
	"fmt"
	redisClient "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/redis"
	"github.com/redis/go-redis/v9"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func InitializeDatabase() *sql.DB {
	dsn := buildDSN()
	conn, connErr := sql.Open("mysql", dsn)
	HandleInitError(connErr)
	configureConnectionPool(conn)
	return conn
}

func InitializeRedisClient() *redis.Client {
	redisAddr := GetEnvOrDefault("REDIS_ADDR", "localhost:6379")
	redisPassword := GetEnvOrDefault("REDIS_PASSWORD", "")
	client, err := redisClient.NewRedisClient(redisAddr, redisPassword)
	if err != nil {
		panic("failed to initialize redis client: " + err.Error())
	}
	return client
}

func HandleInitError(err error) {
	if err != nil {
		panic(err)
	}
}

func buildDSN() string {
	return fmt.Sprintf("root:%s@tcp(%s:%s)/%s?parseTime=true",
		os.Getenv("DB_ROOT_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_DATABASE"),
	)
}

func configureConnectionPool(conn *sql.DB) {
	conn.SetMaxOpenConns(10)
	conn.SetMaxIdleConns(5)
	conn.SetConnMaxLifetime(15 * time.Minute)
}

func GetRequiredEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(key + " environment variable is not set")
	}
	return value
}

func GetEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

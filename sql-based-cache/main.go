package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type Cache struct {
	db *sql.DB
}

func NewCache(dsn string) (*Cache, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	return &Cache{db: db}, nil
}

func (c *Cache) Set(key, value string, expiration time.Duration) error {
	log.Printf("Setting cache for key %v with expiration %v", key, expiration)

	expirationTime := time.Now().Add(expiration)
	query := `REPLACE INTO cache (cache_key, cache_value, expiration_time) VALUES (?, ?, ?)`
	_, err := c.db.Exec(query, key, value, expirationTime)
	return err
}

func (c *Cache) Get(key string) (string, error) {
	var value string
	var expirationTimeStr string
	query := `SELECT cache_value, expiration_time FROM cache WHERE cache_key = ?`
	err := c.db.QueryRow(query, key).Scan(&value, &expirationTimeStr)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("cache miss")
		}
		return "", err
	}

	// Convert expiration_time to time.Time in UTC
	expirationTime, err := time.Parse("2006-01-02 15:04:05", expirationTimeStr)
	if err != nil {
		return "", fmt.Errorf("failed to parse expiration time: %v", err)
	}

	// Check if the cache has expired, now in UTC
	now := time.Now().UTC()
	if now.After(expirationTime) {
		err := c.Delete(key)
		if err != nil {
			log.Printf("Failed to delete expired cache: %v", err)
		}
		return "", fmt.Errorf("cache miss (now: %v, expiration: %v)", now, expirationTime)
	}

	return value, nil
}

func (c *Cache) Delete(key string) error {
	query := `DELETE FROM cache WHERE cache_key = ?`
	_, err := c.db.Exec(query, key)
	return err
}

func main() {
	// Replace with your own DSN
	dsn := "root:root_password@tcp(127.0.0.1:3306)/my_database"

	cache, err := NewCache(dsn)
	if err != nil {
		log.Fatal("Error creating cache:", err)
	}

	r := gin.Default()

	// Set cache value
	r.POST("/cache", func(c *gin.Context) {
		var requestBody struct {
			Key        string        `json:"key" binding:"required"`
			Value      string        `json:"value" binding:"required"`
			Expiration time.Duration `json:"expiration" binding:"required"`
		}

		if err := c.ShouldBindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := cache.Set(requestBody.Key, requestBody.Value, requestBody.Expiration)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set cache"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "cache set successfully"})
	})

	// Get cache value
	r.GET("/cache/:key", func(c *gin.Context) {
		key := c.Param("key")

		value, err := cache.Get(key)
		if err != nil {
			log.Printf("Failed to get cache for key %v: %v", key, err)
			c.JSON(http.StatusNotFound, gin.H{"error": "Cache miss or error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"key": key, "value": value})
	})

	r.Run(":8080") // Listen and serve on 0.0.0.0:8080
}

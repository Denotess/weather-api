package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"weather-api/internal/handlers"
	"weather-api/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func main() {
	if err := godotenv.Load(); err != nil {
		if err := godotenv.Load("cmd/api/.env"); err != nil {
			fmt.Println(err)
		}
	}
	addr := strings.TrimSpace(os.Getenv("ADDRESS"))
	if addr == "" {
		addr = "localhost:6379"
	}
	db := 0
	if raw := strings.TrimSpace(os.Getenv("DATABASE")); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil {
			fmt.Println("invalid DATABASE, using 0")
		} else {
			db = parsed
		}
	}
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: os.Getenv("PASSWORD"),
		DB:       db,
	})
	defer client.Close()
	cache := models.NewRedisCache(client)

	router := gin.Default()
	router.POST("/weather", handlers.Weather(cache))

	router.Run()
}

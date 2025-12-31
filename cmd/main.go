package main

import (
	"fmt"
	"weather-api/internal/handlers"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println(err)
	}
	router := gin.Default()
	router.POST("/weather", handlers.Weather)

	router.Run()
}

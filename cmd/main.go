package main

import (
	"fmt"
	"weather-api/internal/helpers"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println(err)
	}
	helpers.GetWeatherData()
}

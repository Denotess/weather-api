package main

import (
	"context"
	"fmt"
	"weather-api/internal/helpers"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println(err)
	}
	data, err := helpers.GetWeatherData(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%.1f %s\n", data.CurrentConditions.Temp, data.CurrentConditions.Conditions)
}

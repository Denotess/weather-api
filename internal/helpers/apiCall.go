package helpers

import (
	"fmt"
	"os"
)

func GetWeatherData() {
	key := os.Getenv("KEY")
	url := os.Getenv("URL")
	callUrl := fmt.Sprintf(url, key)
	fmt.Println(callUrl)
}

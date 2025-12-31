package handlers

import (
	"context"
	"log"
	"net/http"
	"weather-api/internal/helpers"
	"weather-api/internal/models"

	"github.com/gin-gonic/gin"
)

func Weather(ctx *gin.Context) {
	var query models.UserQuery
	if err := ctx.ShouldBindJSON(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "incorrect json body"})
		log.Println("incorrect json body")
		return
	}
	weatherData, err := helpers.GetWeatherData(context.Background(), query.Location)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Println(err.Error())
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"temp": weatherData.CurrentConditions.Temp, "conditions": weatherData.CurrentConditions.Conditions})

}

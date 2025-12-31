package handlers

import (
	"errors"
	"log"
	"net/http"
	"time"
	"weather-api/internal/helpers"
	"weather-api/internal/models"

	"github.com/gin-gonic/gin"
)

const cacheTTL = 30 * time.Minute

func Weather(cache models.Cache) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var query models.UserQuery
		if err := ctx.ShouldBindJSON(&query); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "incorrect json body"})
			log.Println("incorrect json body")
			return
		}
		normalizedLocation, err := helpers.NormalizeLocation(query.Location)
		if err != nil || normalizedLocation == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "location not set"})
			log.Println("location not set")
			return
		}

		cacheKey := "weather:" + normalizedLocation
		if cache != nil {
			cached, ok, err := cache.Get(ctx.Request.Context(), cacheKey)
			if err != nil {
				log.Println(err.Error())
			} else if ok {
				ctx.JSON(http.StatusOK, gin.H{"temp": cached.Temp, "conditions": cached.Conditions})
				return
			}
		}

		weatherData, err := helpers.GetWeatherData(ctx.Request.Context(), normalizedLocation)
		if err != nil {
			if errors.Is(err, helpers.ErrLocationNotSet) {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				log.Println(err.Error())
				return

			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			log.Println(err.Error())
			return
		}

		response := models.CachedWeather{
			Temp:       weatherData.CurrentConditions.Temp,
			Conditions: weatherData.CurrentConditions.Conditions,
		}
		if cache != nil {
			if err := cache.Set(ctx.Request.Context(), cacheKey, response, cacheTTL); err != nil {
				log.Println(err.Error())
			}
		}

		ctx.JSON(http.StatusOK, gin.H{"temp": response.Temp, "conditions": response.Conditions})
	}
}

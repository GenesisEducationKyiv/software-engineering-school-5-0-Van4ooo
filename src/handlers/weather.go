package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-Van4ooo/src/config"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-Van4ooo/src/services"
)

type WeatherHandler struct {
	WeatherService services.WeatherService
}

func NewWeatherHandler(cfg *config.AppConfig) *WeatherHandler {
	svc := services.NewOpenWeatherService(
		cfg.WeatherAPI.Key,
		cfg.WeatherAPI.BaseURL,
	)
	return &WeatherHandler{WeatherService: svc}
}

func (h *WeatherHandler) GetWeather(c *gin.Context) {
	city := c.Query("city")
	if city == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "city parameter is required"})
		return
	}

	weather, err := h.WeatherService.GetWeather(city)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"weather": weather})
}

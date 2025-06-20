package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-Van4ooo/src/models"
)

type WeatherService interface {
	GetByCity(city string) (*models.Weather, error)
}

type WeatherHandler struct {
	service WeatherService
}

func NewWeatherHandler(svc WeatherService) *WeatherHandler {
	return &WeatherHandler{service: svc}
}

func (h *WeatherHandler) GetWeather(c *gin.Context) {
	city := c.Query("city")
	if city == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "city parameter is required"})
		return
	}

	weather, err := h.service.GetByCity(city)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"weather": weather})
}

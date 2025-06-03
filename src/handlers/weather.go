package handlers

import (
	"github.com/Van4ooo/genesis_case_task/src/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetWeather(c *gin.Context) {
	city := c.Query("city")
	if city == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "City is required"})
		return
	}

	data, err := services.FetchCurrentWeather(city)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}

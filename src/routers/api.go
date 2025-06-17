package routers

import (
	"github.com/gin-gonic/gin"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-Van4ooo/src/handlers"
)

type APIRoutes struct{}

func (a *APIRoutes) Setup(r *gin.Engine) {
	api := r.Group("/api")
	{
		api.GET("/weather", handlers.GetWeather)
		api.POST("/subscribe", handlers.Subscribe)
		api.GET("/confirm/:token", handlers.Confirm)
		api.GET("/unsubscribe/:token", handlers.Unsubscribe)
	}
}

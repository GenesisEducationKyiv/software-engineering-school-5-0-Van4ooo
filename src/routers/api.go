package routers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-Van4ooo/src/config"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-Van4ooo/src/handlers"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-Van4ooo/src/repositories"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-Van4ooo/src/services"
)

type APIRoutes struct {
	WeatherHandler      *handlers.WeatherHandler
	SubscriptionHandler *handlers.SubscriptionHandler
}

func NewAPIRoutes(cfg config.AppSettings, db *gorm.DB) *APIRoutes {
	weatherHandler := handlers.NewWeatherHandler(
		services.NewOpenWeatherService(cfg.GetWeatherAPI()))

	subHandler := handlers.NewSubscriptionHandler(
		services.NewSubscriptionService(repositories.NewGormSubscriptionStorage(db)),
		services.NewSMTPEmailSender(cfg.GetSMTP()))

	return &APIRoutes{
		WeatherHandler:      weatherHandler,
		SubscriptionHandler: subHandler,
	}
}

func (a *APIRoutes) Setup(r *gin.Engine) {
	api := r.Group("/api")
	{
		api.GET("/weather", a.WeatherHandler.GetWeather)
		api.POST("/subscribe", a.SubscriptionHandler.Subscribe)
		api.GET("/confirm/:token", a.SubscriptionHandler.Confirm)
		api.DELETE("/unsubscribe/:token", a.SubscriptionHandler.Unsubscribe)
	}
}

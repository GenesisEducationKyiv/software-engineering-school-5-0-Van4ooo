package config

import (
	"github.com/Van4ooo/genesis_case_task/src/handlers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"os"
	"time"
)

type Config struct {
	DatabaseURL   string
	WeatherAPIKey string
}

func Load() Config {
	db := os.Getenv("DATABASE_URL")
	key := os.Getenv("WEATHER_API_KEY")

	if db == "" || key == "" {
		log.Fatal("DATABASE_URL and WEATHER_API_KEY must be set!!")
	}

	return Config{
		DatabaseURL:   db,
		WeatherAPIKey: key,
	}
}

func SetupAPI(r *gin.Engine) {
	api := r.Group("/api")

	api.GET("/weather", handlers.GetWeather)
	api.POST("/subscribe", handlers.Subscribe)
	api.GET("/confirm/:token", handlers.Confirm)
	api.GET("/unsubscribe/:token", handlers.Unsubscribe)
}

func SetupSwagger(r *gin.Engine) {
	r.StaticFile("/docs/swagger.yaml", "./docs/swagger.yaml")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(
		swaggerFiles.Handler,
		ginSwagger.URL("/docs/swagger.yaml"),
	))
}

func SetupStaticPages(r *gin.Engine) {
	r.GET("/subscribe", handlers.RenderSubscribePage)
	r.Static("/static/", "static/")
}

func SetupCors(r *gin.Engine) {
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
}

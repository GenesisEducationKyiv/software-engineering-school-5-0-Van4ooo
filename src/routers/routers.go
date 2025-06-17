package routers

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-Van4ooo/src/handlers"
)

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

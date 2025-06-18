package routers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-Van4ooo/src/config"
)

type Router interface {
	Setup(r *gin.Engine)
}

func SetupRoutes(r *gin.Engine, cfg *config.AppConfig, db *gorm.DB) {
	apiRoutes := NewAPIRoutes(cfg, db)
	subHandler := apiRoutes.SubscriptionHandler
	webRoutes := NewWebRoutes(subHandler)

	routers := []Router{
		&CORS{},
		apiRoutes,
		webRoutes,
	}
	for _, rt := range routers {
		rt.Setup(r)
	}
}

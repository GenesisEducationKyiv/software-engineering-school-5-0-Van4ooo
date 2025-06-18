package routers

import (
	"github.com/gin-gonic/gin"
	swagFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-Van4ooo/src/handlers"
)

type WebRoutes struct {
	SubscriptionHandler *handlers.SubscriptionHandler
}

func NewWebRoutes(subHandler *handlers.SubscriptionHandler) *WebRoutes {
	return &WebRoutes{SubscriptionHandler: subHandler}
}

func (w *WebRoutes) Setup(r *gin.Engine) {
	r.StaticFile("/docs/swagger.yaml", "./docs/swagger.yaml")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swagFiles.Handler,
		ginSwagger.URL("/docs/swagger.yaml"),
	))
	r.GET("/subscribe", w.SubscriptionHandler.RenderSubscribePage)
	r.Static("/static/", "static/")
}

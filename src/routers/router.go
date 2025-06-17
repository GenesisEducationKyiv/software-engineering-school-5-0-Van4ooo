package routers

import "github.com/gin-gonic/gin"

type Router interface {
	Setup(r *gin.Engine)
}

func SetupRoutes(r *gin.Engine) {
	routers := []Router{
		&CORS{},
		&APIRoutes{},
		&WebRoutes{},
	}
	for _, router := range routers {
		router.Setup(r)
	}
}

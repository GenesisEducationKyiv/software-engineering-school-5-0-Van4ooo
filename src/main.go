package main

import (
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-Van4ooo/src/config"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-Van4ooo/src/db"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-Van4ooo/src/services"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()
	db.Init(cfg.DatabaseURL)

	r := gin.Default()

	config.SetupCors(r)
	config.SetupStaticPages(r)
	config.SetupAPI(r)
	config.SetupSwagger(r)

	go services.StartScheduler()

	r.Run(":8080")
}

package main

import (
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-Van4ooo/src/config"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-Van4ooo/src/db"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-Van4ooo/src/services"
	"github.com/gin-gonic/gin"
	"log"
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

	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}

package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-Van4ooo/src/config"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-Van4ooo/src/db"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-Van4ooo/src/routers"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-Van4ooo/src/services"
)

func main() {
	cfg, err := config.Config()
	if err != nil {
		log.Fatal(err)
	}

	db.Init(cfg.DB)

	r := gin.Default()
	routers.SetupRoutes(r)

	go services.StartScheduler()

	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}

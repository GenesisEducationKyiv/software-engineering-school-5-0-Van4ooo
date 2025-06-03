package main

import (
	"github.com/Van4ooo/genesis_case_task/src/config"
	"github.com/Van4ooo/genesis_case_task/src/db"
	"github.com/Van4ooo/genesis_case_task/src/services"
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

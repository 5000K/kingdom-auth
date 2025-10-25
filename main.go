package main

import (
	"github.com/5000K/kingdom-auth/config"
	"github.com/5000K/kingdom-auth/db"
	"github.com/5000K/kingdom-auth/service"
)

func main() {
	cfg, err := config.Get()

	if err != nil {
		println(err.Error())
		return
	}

	driver, err := db.NewDriver(cfg)

	if err != nil {
		println(err.Error())
		return
	}

	printBanner()

	srv := service.NewService(cfg, driver)
	srv.Run()
}

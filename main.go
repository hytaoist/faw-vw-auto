package main

import (
	"embed"
	"fmt"

	"github.com/hytaoist/faw-vw-auto/config"
	"github.com/hytaoist/faw-vw-auto/delivery/http"
	"github.com/hytaoist/faw-vw-auto/domain"
	"github.com/hytaoist/faw-vw-auto/infrastructure/database"
	"github.com/hytaoist/faw-vw-auto/internal/log"
)

//go:embed assets
var assetsFS embed.FS

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Println("加载配置出错:", err)
		return
	}

	http.SetPushServerURL(cfg.BarkPushServerURL)

	log.Setup()
	db := database.NewPsql()

	faw := http.NewFAW(db)
	faw.LoadWebAPIConfig(cfg)
	faw.LoadAppConfig(cfg)
	faw.BackgroundRunning()
	// faw.Running()

	// Service
	use := domain.NewUsecase(db)
	svr := http.NewServer(use, assetsFS)
	svr.Start()

	fmt.Println(`
	╔══╦══╦╦═╦╗░░╔╗░╔╦╦═╦╗░░╔══╦╦╦══╦═╗
	║═╦╣╔╗║║║║╠══╣╚╦╝║║║║╠══╣╔╗║║╠╗╔╣║║
	║╔╝║╠╣║║║║╠══╩╗║╔╣║║║╠══╣╠╣║║║║║║║║
	╚╝░╚╝╚╩═╩═╝░░░╚═╝╚═╩═╝░░╚╝╚╩═╝╚╝╚═╝

	`)
}

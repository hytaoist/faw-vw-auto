package main

import (
	"embed"
	"fmt"
	"os"

	"github.com/hytaoist/faw-vw-auto/config"
	"github.com/hytaoist/faw-vw-auto/delivery/http"
	"github.com/hytaoist/faw-vw-auto/domain"
	"github.com/hytaoist/faw-vw-auto/infrastructure/database"
	"github.com/hytaoist/faw-vw-auto/internal/log"
)

//go:embed assets
var assetsFS embed.FS

func main() {
	fmt.Println(`
	╔══╦══╦╦═╦╗░░╔╗░╔╦╦═╦╗░░╔══╦╦╦══╦═╗
	║═╦╣╔╗║║║║╠══╣╚╦╝║║║║╠══╣╔╗║║╠╗╔╣║║
	║╔╝║╠╣║║║║╠══╩╗║╔╣║║║╠══╣╠╣║║║║║║║║
	╚╝░╚╝╚╩═╩═╝░░░╚═╝╚═╩═╝░░╚╝╚╩═╝╚╝╚═╝
	
	`)

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

	// If running as a scheduler/task or want a single-run invocation, set FAW_RUN_ONCE=1
	// In that mode, execute the job once and exit to avoid conflicts with external schedulers.
	if os.Getenv("FAW_RUN_ONCE") == "1" {
		faw.Running()
		return
	}

	// Default behavior: start internal cron background runner and the HTTP server.
	faw.BackgroundRunning()

	// Service
	use := domain.NewUsecase(db)
	svr := http.NewServer(use, assetsFS)
	svr.Start()

}

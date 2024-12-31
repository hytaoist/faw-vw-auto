package main

import (
	"embed"
	"fmt"

	"github.com/hytaoist/autosignin/config"
	"github.com/hytaoist/autosignin/delivery/http"
	"github.com/hytaoist/autosignin/domain"
	"github.com/hytaoist/autosignin/infrastructure/database"
	"github.com/hytaoist/autosignin/internal/log"
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
	db := database.NewPsql(cfg.Schema)

	faw := http.NewFAW(db)
	faw.LoadConfig(cfg)
	// faw.BackgroundRunning()
	faw.Running()

	// Service
	use := domain.NewUsecase(db)
	svr := http.NewServer(use, assetsFS)
	svr.Start()
}

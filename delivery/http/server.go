package http

import (
	"context"
	"embed"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	. "github.com/hytaoist/autosignin/domain"
	"github.com/hytaoist/autosignin/internal/log"
)

const (
	port = 9000
)

type server struct {
	http *http.Server
}

func NewServer(use Usecaser, assetsFS embed.FS) *server {
	html := newHTML(assetsFS)
	api := newAPI(use)
	r := newRouter(html, api)
	return &server{
		http: &http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: r,
		},
	}
}

func (s *server) Start() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	go func() {
		<-quit
		s.shutdown()
	}()
	log.Info(fmt.Sprintf("starting http server on port %d", port))
	err := s.http.ListenAndServe()
	if err != nil {
		log.Critical("http: server: listen and server")
		log.Critical(err)
	}
}

func (s *server) shutdown() {
	err := s.http.Shutdown(context.Background())
	if err != nil {
		log.Critical("http: server: shutdown")
		log.Critical(err)
	}
}

package server

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/sornick01/http_chat/chat"
	"github.com/sornick01/http_chat/chat/handlers"
	"github.com/sornick01/http_chat/chat/repository/localcache"
	"github.com/sornick01/http_chat/chat/usecase"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type App struct {
	httpServer *http.Server

	useCase chat.UseCase
}

func NewApp() (*App, error) {

	expiresAt, err := time.ParseDuration("15m")
	if err != nil {
		return nil, err
	}

	return &App{
		useCase: usecase.NewChat(
			localcache.NewLocalStorage(),
			"hash_salt",
			[]byte("signing_key"),
			expiresAt),
	}, nil
}

func (a *App) Run(port string) {
	mainRoute := chi.NewRouter()
	handlers.RegisterHTTPEndpoints(mainRoute, a.useCase)

	a.httpServer = &http.Server{
		Addr:    ":" + port,
		Handler: mainRoute,
	}

	go func() {
		if err := a.httpServer.ListenAndServe(); err != nil {
			if err == http.ErrServerClosed {
				log.Println("Server closed")
				return
			}
			log.Fatalf("Failed to listen and serve: %+v", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	stopAppCh := make(chan struct{})
	go func() {
		log.Println("Captured signal: ", <-quit)
		log.Println("Gracefully shutting down server...")

		if err := a.httpServer.Shutdown(context.Background()); err != nil {
			log.Fatal("Can't shutdown main server: ", err.Error())
		}
		stopAppCh <- struct{}{}
	}()

	<-stopAppCh
}

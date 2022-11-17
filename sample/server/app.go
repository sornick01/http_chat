package server

import (
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"sample/chat"
	"sample/chat/handlers"
	"sample/chat/repository/localcache"
	"sample/chat/usecase"
	"time"
)

type App struct {
	httpServer *http.Server

	useCase chat.UseCase
}

func NewApp() (*App, error) {

	expiresAt, err := time.ParseDuration("5m")
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
	log.Fatal(a.httpServer.ListenAndServe())
}

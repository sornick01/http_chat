package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sornick01/http_chat/chat"
)

func RegisterHTTPEndpoints(rout *chi.Mux, useCase chat.UseCase) {
	h := NewHandler(useCase)

	rout.Use(middleware.RequestID)
	rout.Use(middleware.Logger)
	rout.Post("/signUp", h.SignUp)
	rout.Post("/signIn", h.SignIn)

	r := chi.NewRouter()
	authMiddleware := NewAuthMiddleware(useCase)

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(authMiddleware.Auth)
	r.Post("/send", h.SendMessage)
	r.Get("/read/private", h.ReadPrivateMessage)
	r.Get("/read/global", h.ReadGlobalMessages)
	rout.Mount("/api", r)
}

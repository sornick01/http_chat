package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/sornick01/http_chat/chat"
)

func RegisterHTTPEndpoints(rout *chi.Mux, useCase chat.UseCase) {
	h := NewHandler(useCase)

	rout.Post("/signUp", h.SignUp)
	rout.Post("/signIn", h.SignIn)

	r := chi.NewRouter()
	middleware := NewAuthMiddleware(useCase)

	r.Use(middleware.Auth)
	r.Post("/send", h.SendMessage)
	rout.Mount("/api", r)
}

package handlers

import (
	"encoding/json"
	"errors"
	chat2 "github.com/sornick01/http_chat/internal/chat"
	"github.com/sornick01/http_chat/internal/models"
	"io"
	"net/http"
)

type Handlers struct {
	useCase chat2.UseCase
}

type signInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type signInResponse struct {
	Token string `json:"token"`
}

type messageInput struct {
	Recipient string `json:"recipient"`
	Text      string `json:"text"`
}

func NewHandler(useCase chat2.UseCase) *Handlers {
	return &Handlers{
		useCase: useCase,
	}
}

func (h *Handlers) SignUp(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	inp, err := jsonToSignInput(body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := h.useCase.SignUp(r.Context(), inp.Username, inp.Password); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handlers) SignIn(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	inp, err := jsonToSignInput(body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	token, err := h.useCase.SignIn(r.Context(), inp.Username, inp.Password)
	if err != nil {
		if err == chat2.ErrUserNotFound {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := &signInResponse{Token: token}
	respJson, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(respJson)
}

func (h *Handlers) SendMessage(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	mes, err := jsonToMessageInput(body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user := r.Context().Value(chat2.CtxUserKey).(*models.User)

	if mes.Recipient == "" {
		err = h.useCase.AddGlobalMessage(r.Context(), user, mes.Text)
	} else {
		err = h.useCase.AddPrivateMessage(r.Context(), user, mes.Recipient, mes.Text)
	}
	if err != nil {
		status := http.StatusInternalServerError
		if err == chat2.ErrUserNotFound {
			status = http.StatusBadRequest
		}
		w.WriteHeader(status)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handlers) ReadPrivateMessage(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(chat2.CtxUserKey).(*models.User)

	messages, err := h.useCase.GetPrivateMessages(r.Context(), user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(messages)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func (h *Handlers) ReadGlobalMessages(w http.ResponseWriter, r *http.Request) {
	messages, err := h.useCase.GetGlobalMessages(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(messages)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func jsonToSignInput(body []byte) (*signInput, error) {
	inp := signInput{}

	if err := json.Unmarshal(body, &inp); err != nil {
		return nil, err
	}

	if inp.Username == "" {
		return nil, errors.New("empty username")
	}

	return &inp, nil
}

func jsonToMessageInput(body []byte) (*messageInput, error) {
	mes := messageInput{}

	if err := json.Unmarshal(body, &mes); err != nil {
		return nil, err
	}

	if mes.Text == "" {
		return nil, errors.New("empty message")
	}

	return &mes, nil
}

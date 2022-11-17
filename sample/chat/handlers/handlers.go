package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"sample/chat"
	"sample/models"
)

type Handlers struct {
	useCase chat.UseCase
}

type signInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type signInResponse struct {
	Token string `json:"token"`
}

type messageInput struct {
	recipient string `json:"recipient"`
	text      string `json:"text"`
}

func NewHandler(useCase chat.UseCase) *Handlers {
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
		if err == chat.ErrUserNotFound {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	resp := &signInResponse{Token: token}
	respJson, err := json.Marshal(resp)

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

	user := r.Context().Value(chat.CtxUserKey).(*models.User)

	if mes.recipient == "" {
		err = h.useCase.AddGlobalMessage(r.Context(), user, mes.text)
	} else {
		err = h.useCase.AddPrivateMessage(r.Context(), user, mes.recipient, mes.text)
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
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

	if mes.text == "" {
		return nil, errors.New("empty message")
	}

	return &mes, nil
}

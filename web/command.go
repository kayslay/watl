package web

import (
	"net/http"

	whatzapp "github.com/Rhymen/go-whatsapp"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

func (a *App) Command(w http.ResponseWriter, r *http.Request) {
	refCode := chi.URLParam(r, "code")
	h, _, err := a.getSession(refCode)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, map[string]string{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	phone := chi.URLParam(r, "phone")

	message := whatzapp.TextMessage{
		Info: whatzapp.MessageInfo{
			RemoteJid: phone + "@s.whatsapp.net",
		},
		Text: "#!" + chi.URLParam(r, "command"),
	}

	s, err = h.Command(message)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, map[string]string{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	render.JSON(w, r, map[string]string{
		"status":  "success",
		"message": "command successful. " + s,
	})
}

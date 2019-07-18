package web

import (
	"io"
	"net/http"

	"github.com/go-chi/chi"
)

func Router(w io.WriteCloser) http.Handler {

	r := chi.NewRouter()
	app := NewApp(w)

	r.Route("/", func(r chi.Router) {

		r.Get("/login", app.Login)
		r.Get("/logout/{code}", app.Logout)
	})

	return r
}

package web

import (
	"github/kayslay/watl/whatsapp"
	"io"
	"sync"
)

type App struct {
	session map[string]*whatsapp.Handler
	code    map[string]string
	Writer  io.WriteCloser
	sync.Mutex
}

func NewApp(w io.WriteCloser) *App {
	return &App{
		session: map[string]*whatsapp.Handler{},
		code:    map[string]string{},
		Writer:  w,
	}
}

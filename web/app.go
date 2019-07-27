package web

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github/kayslay/watl/store"
	"github/kayslay/watl/whatsapp"
	"io"
	"log"
	"sync"
	"time"

	whatzapp "github.com/Rhymen/go-whatsapp"
)

const (
	dataPath = "data"
)

type App struct {
	session map[string]*whatsapp.Handler
	code    map[string]string
	Writer  io.WriteCloser
	str     *store.MgoStore
	sync.Mutex
}

func NewApp(w io.WriteCloser) *App {
	a := &App{
		session: map[string]*whatsapp.Handler{},
		code:    map[string]string{},
		Writer:  w,
		str:     store.NewMgo(),
	}
	// load previous sessions
	a.initSessions()
	return a
}

func (a App) saveSession(clientID, uniqueCode string, session whatzapp.Session) error {
	var b bytes.Buffer
	encoder := gob.NewEncoder(&b)
	err := encoder.Encode(session)
	if err != nil {
		return err
	}

	return a.str.SaveSession(store.Session{
		ClientID:   clientID,
		UniqueCode: uniqueCode,
		Session:    b.Bytes(),
		CreatedAt:  time.Now(),
	})

}

func (a *App) initSessions() {
	log.Println("initializing previous sessions")
	ss, err := a.str.GetSessions()
	if err != nil {
		log.Println("error init session", err)
		return
	}

	for _, s := range ss {
		session := whatzapp.Session{}
		var b = bytes.NewBuffer(s.Session)
		decoder := gob.NewDecoder(b)
		err = decoder.Decode(&session)
		if err != nil {
			log.Println("error decoding session", err)
			continue
		}
		wac, h, err := a.conn()
		if err != nil {
			log.Println("error connecting session", err)
			continue
		}
		wac.RestoreWithSession(session)
		a.code[s.UniqueCode] = s.ClientID
		a.session[s.ClientID] = h
		// load contacts
		h.LoadContact()
		go func() {
			<-h.Close
			delete(a.session, s.ClientID)
			delete(a.code, s.UniqueCode)
			// delete session
			fmt.Println("logged out ", a.session)
		}()
	}
}

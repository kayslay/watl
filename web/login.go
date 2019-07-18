package web

import (
	"fmt"
	"github/kayslay/watl/store"
	"github/kayslay/watl/whatsapp"
	"log"
	"net/http"
	"time"

	whatzapp "github.com/Rhymen/go-whatsapp"
	"github.com/dchest/uniuri"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

func (a *App) Login(w http.ResponseWriter, r *http.Request) {
	// get the qr code
	uniqueCode := "WHT_" + uniuri.NewLen(8)

	wac, err := whatzapp.NewConn(time.Minute)
	if err != nil {
		fmt.Println("could not connect at the moment " + err.Error())
		return
	}

	str := store.NewMgo()

	h, err := whatsapp.NewHandler(wac, a.Writer, str)
	if err != nil {
		log.Fatal(err)
	}
	//Add handler
	wac.AddHandler(h)
	qr := make(chan string)

	go func() {
		session, err := whatsapp.QrLogin(wac, qr)
		if err != nil {
			fmt.Println("could not login at the moment " + err.Error())
			return
		}
		// TODO save the session somewhere
		// AWS S3
		a.Lock()
		defer a.Unlock()
		// set the session
		a.session[session.ClientId] = h
		// set the handler id
		h.SetClientID(session.ClientId)
		// set the session refCode
		a.code[uniqueCode] = session.ClientId
	}()

	render.JSON(w, r, map[string]string{
		"qr":      <-qr,
		"refCode": uniqueCode,
		"message": "click link after scan",
		"link":    fmt.Sprintf("http:localhost:8000/%s/status", uniqueCode),
	})
}

func (a *App) Logout(w http.ResponseWriter, r *http.Request) {
	// get the qr code
	refCode := chi.URLParam(r, "code")
	a.Lock()
	defer a.Unlock()

	sessID, ok := a.code[refCode]
	if !ok {
		render.JSON(w, r, map[string]string{
			"status":  "success",
			"message": "ref code does not exist",
		})
		return
	}
	h, ok := a.session[sessID]
	if !ok {
		render.JSON(w, r, map[string]string{
			"status":  "success",
			"message": "could not find session",
		})
		return
	}

	h.Logout()
	delete(a.session, sessID)
	delete(a.code, refCode)
	fmt.Println(a.session, a.code)
	render.JSON(w, r, map[string]string{
		"status":  "success",
		"message": "logout successful",
	})

}

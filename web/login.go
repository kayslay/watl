package web

import (
	"errors"
	"fmt"
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

	wac, h, err := a.conn()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, map[string]string{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	qr := make(chan string)

	go func() {
		sess, err := whatsapp.QrLogin(wac, qr)
		if err != nil {
			fmt.Println("could not login at the moment " + err.Error())
			return
		}
		userID := h.GetInfo().Wid
		// TODO save the session somewhere
		a.saveSession(userID, uniqueCode, sess)
		// AWS S3
		a.Lock()
		defer a.Unlock()
		// map  handler to userID
		a.session[userID] = h
		// map userID to uniqueCode
		a.code[uniqueCode] = userID
		// set handler.ID ot uniqueCode
		h.ID = uniqueCode
		// setup handler
		h.Setup()
		go func() {
			<-h.Close
			a.Lock()
			defer a.Unlock()
			delete(a.session, userID)
			delete(a.code, uniqueCode)
			// delete session
			a.str.DeleteSession(userID)
			fmt.Println("logged out ", a.session)
		}()
	}()

	qrcode := <-qr
	log.Println(qrcode)
	if qrcode == "" {
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, map[string]string{
			"status":  "error",
			"message": "could not get qr code. check network connection",
		})
		return
	}

	render.JSON(w, r, map[string]string{
		"qr":      qrcode,
		"refCode": uniqueCode,
		"message": "click link after scan",
		"link":    fmt.Sprintf("http:localhost:8000/%s/", uniqueCode),
	})
}

func (a *App) Logout(w http.ResponseWriter, r *http.Request) {
	// get the qr code
	refCode := chi.URLParam(r, "code")
	h, sessID, err := a.getSession(refCode)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, map[string]string{
			"status":  "error",
			"message": err.Error(),
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

func (a *App) getSession(refCode string) (*whatsapp.Handler, string, error) {
	a.Lock()
	defer a.Unlock()

	userID, ok := a.code[refCode]
	if !ok {
		return nil, "", errors.New("ref code does not exist")
	}
	h, ok := a.session[userID]
	if !ok {
		return nil, "", errors.New("session expired")

	}
	return h, userID, nil
}

func (a App) conn() (*whatzapp.Conn, *whatsapp.Handler, error) {
	wac, err := whatzapp.NewConn(time.Minute)
	if err != nil {
		return nil, nil, err
	}

	str := a.str

	h, err := whatsapp.NewHandler(wac, a.Writer, str)
	if err != nil {
		return nil, nil, err
	}

	//Add handler
	wac.AddHandler(h)
	return wac, h, err
}

package whatsapp

import (
	"encoding/gob"
	"fmt"
	"os"
	"path"

	qrcodeTerminal "github.com/Baozisoftware/qrcode-terminal-go"

	whatzapp "github.com/Rhymen/go-whatsapp"
)

func Login(wac *whatzapp.Conn) error {
	//load saved session
	session, err := readSession()
	if err == nil {
		//restore session
		session, err = wac.RestoreWithSession(session)
		if err != nil {
			fmt.Errorf("restoring failed: %v\n", err)
			session, err = QrLogin(wac)
			return err
		}
	} else {
		//no saved session -> regular login
		session, err = QrLogin(wac)
		return err
	}

	//save session
	err = WriteSession(session)
	if err != nil {
		return fmt.Errorf("error saving session: %v\n", err)
	}
	return nil
}

func QrLogin(wac *whatzapp.Conn, ch ...chan string) (whatzapp.Session, error) {
	qr := make(chan string)
	go func() {
		terminal := qrcodeTerminal.New()
		qrcode := <-qr
		terminal.Get(qrcode).Print()
		if len(ch) > 0 {
			ch[0] <- qrcode
		}
	}()
	session, err := wac.Login(qr)
	if err != nil {
		return whatzapp.Session{}, fmt.Errorf("error during login: %v\n", err)
	}
	return session, nil
}

func readSession() (whatzapp.Session, error) {
	fmt.Println("reading session")
	session := whatzapp.Session{}
	file, err := os.Open(os.TempDir() + "/whatsappSession.gob")
	if err != nil {
		return session, err
	}
	defer file.Close()
	decoder := gob.NewDecoder(file)
	err = decoder.Decode(&session)
	if err != nil {
		return session, err
	}
	return session, nil
}

func WriteSession(session whatzapp.Session) error {
	file, err := os.Create(path.Join(os.TempDir(), "/"+session.ClientId+"_whatsappSession.gob"))
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := gob.NewEncoder(file)
	err = encoder.Encode(session)
	if err != nil {
		return err
	}
	return nil
}

func (h *Handler) Logout() error {
	h.c.RemoveHandler(h)
	return h.c.Logout()
}

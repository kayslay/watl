package main

import (
	"fmt"
	"github/kayslay/watl/config"
	"github/kayslay/watl/store"
	"github/kayslay/watl/whatsapp"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	whatzapp "github.com/Rhymen/go-whatsapp"
)

func init() {
	godotenv.Load("cmd/cli/.env")
}

func main() {

	wac, err := whatzapp.NewConn(50 * time.Second)
	if err != nil {
		log.Fatalf("error creating connection: %v\n", err)
	}
	w, err := os.OpenFile("data/messages.json", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	defer w.Close()
	if err != nil {
		log.Fatalf("error creating file: %v\n", err)
	}

	config.New(config.NewMongoConnect)
	str := store.NewMgo()

	h, err := whatsapp.NewHandler(wac, w, str)
	if err != nil {
		log.Fatal(err)
	}
	//Add handler
	wac.AddHandler(h)

	//login or restore
	if err := whatsapp.Login(wac); err != nil {
		log.Fatalf("error logging in: %v\n", err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	//Disconnect safe
	fmt.Println("Shutting down now.")
	session, err := wac.Disconnect()
	if err != nil {
		log.Fatalf("error disconnecting: %v\n", err)
	}

	if err := whatsapp.WriteSession(session); err != nil {
		log.Fatalf("error saving session: %v", err)
	}
}

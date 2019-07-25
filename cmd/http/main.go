package main

import (
	"fmt"
	"github/kayslay/watl/config"
	"github/kayslay/watl/web"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/middleware"

	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load("cmd/http/.env")
}

type nilWriteCloser struct {
}

func (n nilWriteCloser) Write(p []byte) (int, error) {
	return 0, nil
}
func (n nilWriteCloser) Close() error {
	return nil
}

func main() {
	var (
		w   io.WriteCloser
		err error
	)
	if os.Getenv("ENV") == "production" {
		w, err = nilWriteCloser{}, nil
	} else {
		w, err = os.OpenFile("data/messages.json", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	}
	defer w.Close()
	if err != nil {
		log.Fatalf("error creating file: %v\n", err)
	}

	config.New(config.NewMongoConnect)

	var r = chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)
	r.Route("/v1", func(r chi.Router) {
		r.Mount("/", web.Router(w))
	})

	r.Handle("/public/*", http.StripPrefix("/public", http.FileServer(http.Dir("./public"))))

	port := "8000"
	envPort := os.Getenv("PORT")
	if envPort != "" {
		port = envPort
	}
	log.Println("cooler loling")

	fmt.Println("server started on port ", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

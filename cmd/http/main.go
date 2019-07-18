package main

import (
	"fmt"
	"github/kayslay/watl/config"
	"github/kayslay/watl/web"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load("cmd/http/.env")
}

func main() {

	w, err := os.OpenFile("data/messages.json", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	defer w.Close()
	if err != nil {
		log.Fatalf("error creating file: %v\n", err)
	}

	config.New(config.NewMongoConnect)

	var r = chi.NewRouter()

	r.Route("/v1", func(r chi.Router) {
		r.Mount("/", web.Router(w))
	})
	port := "8000"
	envPort := os.Getenv("PORT")
	if envPort != "" {
		port = envPort
	}
	log.Println("cooler loling")

	fmt.Println("server started on port ", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

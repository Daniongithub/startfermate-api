package main

import (
	"log"
	"net/http"

	"github.com/rs/cors"

	"github.com/Daniongithub/startfermate-api/handlers"
)

const version = "3.1"

func main() {

	mux := http.NewServeMux()

	mux.HandleFunc("/fermata", handlers.Fermata)
	mux.HandleFunc("/bacino", handlers.Bacino)
	mux.HandleFunc("/versione", handlers.Version(version))

	handler := cors.Default().Handler(mux)

	log.Println("API attiva su :3005")

	log.Fatal(http.ListenAndServe(":3005", handler))
}

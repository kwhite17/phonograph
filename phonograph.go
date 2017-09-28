package main

import (
	"net/http"
	"os"

	"github.com/kwhite17/phonograph/server"
)

func main() {
	http.Handle("/search", server.SpotifySearchHandler{})
	http.Handle("/", http.FileServer(http.Dir("./static")))
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8080"
	}
	http.ListenAndServe(":"+port, nil)
}

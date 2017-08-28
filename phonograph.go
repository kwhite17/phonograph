package main

import (
	"net/http"

	"github.com/kwhite17/phonograph/server"
)

func main() {
	http.Handle("/search", server.SpotifySearchHandler{})
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.ListenAndServe(":8080", nil)
}

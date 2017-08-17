package server

import (
	"encoding/json"
	"net/http"

	"github.com/kwhite17/phonograph/client"
	"github.com/kwhite17/phonograph/search"
)

type spotifySearchHandler struct{}

func (sh spotifySearchHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	params := request.URL.Query()
	sCli := client.InitSpotifyClient()
	artist := sCli.FindArtist(params.Get("artist"))
	if search.IsNil(artist) {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write(nil)
	}
	candidate := sCli.FindArtist(params.Get("candidate"))
	if search.IsNil(candidate) {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write(nil)
	}
	artistNode := &search.ArtistNode{Value: artist}
	candidateNode := &search.ArtistNode{Value: candidate}
	result := search.BidirectionalBfs(artistNode, candidateNode, nil)
	jsonResult, err := json.Marshal(result)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write(nil)
	}
	writer.Write(jsonResult)
}

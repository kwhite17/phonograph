package server

import (
	"encoding/json"
	"net/http"

	"github.com/kwhite17/phonograph/client"
	"github.com/kwhite17/phonograph/model"
	"github.com/kwhite17/phonograph/search"
)

type SpotifySearchHandler struct{}

func (sh SpotifySearchHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	params := request.URL.Query()
	sCli := client.InitSpotifyClient()
	artist := sCli.FindArtist(params.Get("artist"))
	if model.IsNil(artist) {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write(nil)
	}
	candidate := sCli.FindArtist(params.Get("candidate"))
	if model.IsNil(candidate) {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write(nil)
	}
	artistNode := &model.ArtistNode{Value: artist}
	candidateNode := &model.ArtistNode{Value: candidate}
	result := search.BidirectionalBfs(artistNode, candidateNode, nil)
	jsonResult, err := json.Marshal(result)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write(nil)
	}
	writer.Write(jsonResult)
}

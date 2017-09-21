package server

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"

	"github.com/kwhite17/phonograph/client"
	"github.com/kwhite17/phonograph/model"
	"github.com/kwhite17/phonograph/search"
)

type SpotifySearchHandler struct{}

func (sh SpotifySearchHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	params := request.URL.Query()
	sCli := client.InitSpotifyClient()
	artist := sCli.FindArtist(url.QueryEscape(params.Get("artist")))
	log.Printf("Artist: %v\n", artist)
	if artist.ID == "" {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	candidate := sCli.FindArtist(url.QueryEscape(params.Get("collaborator")))
	log.Printf("Artist: %v\n", candidate)
	if candidate.ID == "" {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	artistNode := &model.ArtistNode{Value: artist}
	sCli.AddToCache(artist.ID, artistNode)
	candidateNode := &model.ArtistNode{Value: candidate}
	sCli.AddToCache(candidate.ID, candidateNode)
	nodeResult := search.BidirectionalBfs(artistNode, candidateNode, sCli)
	result := make([]model.Artist, 0)
	for i := 0; i < len(nodeResult); i++ {
		result = append(result, nodeResult[i].GetValue().(model.Artist))
	}
	log.Println(result)
	jsonResult, err := json.Marshal(result)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
	}
	writer.Write(jsonResult)
}

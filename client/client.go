package client

import (
	"encoding/json"
	"log"
	"net/http"

	"fmt"

	"github.com/bitly/go-simplejson"
	"github.com/kwhite17/phonograph/model"
)

type Client interface {
	Expand(value model.Node) []model.Node
}

type SpotifyClient struct {
	client   *http.Client
	baseUrl  string
	protocol string
	auth     string
}

//TODO: Get Dev Key for Spotify

func InitSpotifyClient() *SpotifyClient {
	return &SpotifyClient{client: &http.Client{}, baseUrl: "api.spotify.com/v1", protocol: "https"}
}

//TODO: Extract Parser functionality into testable methods

func (sClient *SpotifyClient) FindArtist(name string) model.Artist {
	endpoint := sClient.buildUrl("search", map[string]interface{}{"type": "artist", "q": "name"})
	resp, err := sClient.client.Get(endpoint)
	if err != nil {
		log.Println(fmt.Errorf("GetArtist Error: %v", err))
		return model.Artist{}
	}
	parser, err := simplejson.NewFromReader(resp.Body)
	if err != nil {
		log.Println(fmt.Errorf("GetArtist Error - Read JSON: %v", err))
		return model.Artist{}
	}
	return sClient.parseArtistFromJson(parser)
}

func (sClient *SpotifyClient) parseArtistFromJson(parser *simplejson.Json) model.Artist {
	artist := model.Artist{}
	parser = parser.GetPath("artists", "items").GetIndex(0)
	data, err := parser.Map()
	bytes, err := json.Marshal(data)
	if err != nil {
		log.Println(fmt.Errorf("GetArtist Error - Marshal:  %v", err))
		return model.Artist{}
	}
	err = json.Unmarshal(bytes, &artist)
	if err != nil {
		log.Println(fmt.Errorf("GetArtist Error - Unmarshal: %v", err))
		return model.Artist{}
	}
	return artist
}

func (sClient *SpotifyClient) getArtistAlbums(artist model.Artist) []string {
	albumEndpoint := sClient.buildUrl(fmt.Sprintf("artists/%s/albums", artist.ID), nil)
	albumResponse, err := sClient.client.Get(albumEndpoint)
	if err != nil {
		log.Println(fmt.Errorf("GetAssociatedArtists Error - Album Request Error: %v", err))
		return nil
	}
	parser, err := simplejson.NewFromReader(albumResponse.Body)
	if err != nil {
		log.Println(fmt.Errorf("GetAssociatedArtists - Albums Request Error: %v", err))
		return nil
	}
	albums, err := parser.GetPath("items").Array()
	albumIds := make([]string, 0)
	for _, v := range albums {
		mapV := v.(map[string]interface{})
		albumIds = append(albumIds, mapV["id"].(string))
	}
	return albumIds
}

func (sClient *SpotifyClient) getAssociatedArtists(artist model.Artist) []model.Artist {
	albumIds := sClient.getArtistAlbums(artist)
	albumsEndpoint := sClient.buildUrl("albums", map[string]interface{}{"ids": albumIds[0]})
	for i := 1; i < len(albumIds); i++ {
		albumsEndpoint = fmt.Sprintf("%s,%s", albumsEndpoint, albumIds[i])
	}
	albumsResponse, err := sClient.client.Get(albumsEndpoint)
	if err != nil {
		log.Println(fmt.Errorf("GetAssociatedArtists - Albums Request Error: %v", err))
		return nil
	}
	albumsParser, err := simplejson.NewFromReader(albumsResponse.Body)
	if err != nil {
		log.Println(fmt.Errorf("GetAssociatedArtists Error - Read JSON: %v", err))
		return nil
	}
	return sClient.parseAssociatedArtistsFromJson(albumsParser, artist.ID)
}

func (sClient *SpotifyClient) parseAssociatedArtistsFromJson(parser *simplejson.Json, sourceID string) []model.Artist {
	albums, err := parser.Get("albums").Array()
	if err != nil {
		log.Println(fmt.Errorf("GetAssociatedArtists - AlbumParser Error: %v", err))
		return nil
	}

	associatedArtists := make([]model.Artist, 0)
	for i := 0; i < len(albums); i++ {
		tracksParser := parser.Get("albums").GetIndex(i).Get("tracks").Get("items")
		tracks, err := tracksParser.Array()
		if err != nil {
			log.Println(fmt.Errorf("GetAssociatedArtists - TrackParser Error: %v", err))
			return associatedArtists
		}
		for j := 0; j < len(tracks); j++ {
			trackData, err := tracksParser.GetIndex(j).Map()
			trackBytes, err := json.Marshal(trackData)
			if err != nil {
				log.Println(fmt.Errorf("GetAssociatedArtists - TrackParser Error: %v", err))
				return associatedArtists
			}
			track := model.Track{}
			json.Unmarshal(trackBytes, &track)
			artistParser := tracksParser.GetIndex(j).Get("artists")
			artists, err := artistParser.Array()
			log.Println(len(artists))
			for k := 0; k < len(artists); k++ {
				artistData, err := artistParser.GetIndex(k).Map()
				artistBytes, err := json.Marshal(artistData)
				if err != nil {
					log.Println(fmt.Errorf("GetAssociatedArtists - Artist Parser Error: %v", err))
					return associatedArtists
				}
				associatedArtist := model.Artist{}
				json.Unmarshal(artistBytes, &associatedArtist)
				if associatedArtist.ID != sourceID {
					associatedArtist.RelatedSong = &track
					associatedArtists = append(associatedArtists, associatedArtist)
				}
			}
		}
	}
	return associatedArtists
}

func (sClient *SpotifyClient) buildUrl(endpoint string, params map[string]interface{}) string {
	resource := fmt.Sprintf("%s://%s/%s", sClient.protocol, sClient.baseUrl, endpoint)
	if params != nil {
		paramString := "?"
		for k, v := range params {
			if paramString == "?" {
				paramString = fmt.Sprintf("%s%s=%v", paramString, k, v)
			} else {
				paramString = fmt.Sprintf("%s&%s=%v", paramString, k, v)
			}
		}
		resource = resource + paramString
	}
	return resource
}

func (sCli SpotifyClient) Expand(parent model.Node) []model.Node {
	artist := parent.GetValue().(model.Artist)
	associatedArtists := sCli.getAssociatedArtists(artist)
	result := make([]model.Node, len(associatedArtists))
	for k, v := range associatedArtists {
		result[k] = &model.ArtistNode{Value: v}
	}
	return result
}

type GraphClient struct {
	Graph map[interface{}][]interface{}
}

func (gCli GraphClient) Expand(parent model.Node) []model.Node {
	descendants := gCli.Graph[parent]
	nodes := make([]model.Node, len(descendants))
	for k, v := range descendants {
		nodes[k] = v.(*model.GenericNode)
	}
	return nodes
}

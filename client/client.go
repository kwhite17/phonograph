package client

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httputil"

	"fmt"

	"github.com/bitly/go-simplejson"
	"github.com/kwhite17/phonograph/model"
)

type Client interface {
	Expand(value model.Node) []model.Node
}

type SpotifyClient struct {
	client            *http.Client
	baseUrl           string
	protocol          string
	auth              string
	discoveredArtists map[string]model.Node
}

func InitSpotifyClient() *SpotifyClient {
	cli := &SpotifyClient{client: &http.Client{}, baseUrl: "api.spotify.com/v1", protocol: "https", discoveredArtists: make(map[string]model.Node)}
	req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", bytes.NewBufferString("grant_type=client_credentials"))
	if err != nil {
		reqBytes, _ := httputil.DumpRequestOut(req, true)
		log.Println(fmt.Errorf("InitClient Error - Error:\n %v,\n Request:\n %q", err, reqBytes))
		return cli
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := cli.client.Do(req)
	if err != nil || resp.StatusCode-200 >= 200 {
		reqBytes, _ := httputil.DumpRequest(req, true)
		respBytes, _ := httputil.DumpResponse(resp, true)
		log.Println(fmt.Errorf("InitClient Error - Token: %v,\n Request: %q,\n Response: %q", err, reqBytes, respBytes))
		return cli
	}
	authObject, err := simplejson.NewFromReader(resp.Body)
	auth, err := authObject.GetPath("access_token").String()
	if err != nil {
		log.Println(fmt.Errorf("InitClient Error - Parser: %v", err))
		return cli
	}
	cli.auth = auth
	return cli
}

func (sClient *SpotifyClient) FindArtist(name string) model.Artist {
	log.Println(fmt.Sprintf("Service Call: %s", name))
	req := sClient.buildRequest("search", map[string]interface{}{"type": "artist", "q": name})
	if req == nil {
		return model.Artist{}
	}
	resp, err := sClient.client.Do(req)
	if err != nil || resp.StatusCode > 399 {
		reqBytes, _ := httputil.DumpRequest(req, true)
		respBytes, _ := httputil.DumpResponse(resp, true)
		log.Println(fmt.Errorf("GetArtist Error: %v,\n Request: %q,\n Response: %q", err, reqBytes, respBytes))
		return model.Artist{}
	}
	parser, err := simplejson.NewFromReader(resp.Body)
	if err != nil {
		log.Println(fmt.Errorf("GetArtist Error - Read JSON: %v", err))
		return model.Artist{}
	}
	return sClient.parseArtistFromJson(parser)
}

func (sClient *SpotifyClient) AddToCache(id string, node model.Node) {
	sClient.discoveredArtists[id] = node
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
	log.Println(fmt.Sprintf("Service Call - Albums: %s", artist.Name))
	req := sClient.buildRequest(fmt.Sprintf("artists/%s/albums", artist.ID), nil)
	if req == nil {
		return nil
	}
	resp, err := sClient.client.Do(req)
	if err != nil || resp.StatusCode-200 >= 200 {
		reqBytes, _ := httputil.DumpRequest(req, true)
		respBytes, _ := httputil.DumpResponse(resp, true)
		log.Println(fmt.Errorf("InitClient Error - Token: %v,\n Request: %q,\n Response: %q", err, reqBytes, respBytes))
		return nil
	}
	parser, err := simplejson.NewFromReader(resp.Body)
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
	log.Println(fmt.Sprintf("Service Call - Associated Artists: %s", artist.Name))
	albumIds := sClient.getArtistAlbums(artist)
	albumString := albumIds[0]
	for i := 1; i < len(albumIds); i++ {
		albumString = fmt.Sprintf("%s,%s", albumString, albumIds[i])
	}
	req := sClient.buildRequest("albums", map[string]interface{}{"ids": albumString})
	if req == nil {
		return nil
	}
	resp, err := sClient.client.Do(req)
	if err != nil || resp.StatusCode-200 >= 200 {
		reqBytes, _ := httputil.DumpRequest(req, true)
		respBytes, _ := httputil.DumpResponse(resp, true)
		log.Println(fmt.Errorf("GetAssociatedArtists - Albums Request Error: %v,\n Request: %q,\n Response: %q", err, reqBytes, respBytes))
		return nil
	}
	albumsParser, err := simplejson.NewFromReader(resp.Body)
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

func (sClient *SpotifyClient) buildRequest(endpoint string, params map[string]interface{}) *http.Request {
	req, err := http.NewRequest("GET", sClient.buildUrl(endpoint, params), nil)
	if err != nil {
		log.Println(fmt.Errorf("BuildRequest - Error: %v", err))
		return nil
	}
	req.Header.Add("Authorization", "Bearer "+sClient.auth)
	return req
}

func (sClient *SpotifyClient) Expand(parent model.Node) []model.Node {
	artist := parent.GetValue().(model.Artist)
	associatedArtists := sClient.getAssociatedArtists(artist)
	result := make([]model.Node, len(associatedArtists))
	for k, v := range associatedArtists {
		node, ok := sClient.discoveredArtists[v.ID]
		if ok {
			result[k] = node
		} else {
			result[k] = &model.ArtistNode{Value: v}
			sClient.discoveredArtists[v.ID] = result[k]
		}
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

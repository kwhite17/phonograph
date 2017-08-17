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
	FindArtist(name string) model.Artist
}

type SpotifyClient struct {
	client   *http.Client
	baseUrl  string
	protocol string
	auth     string
}

func InitSpotifyClient() *SpotifyClient {
	return &SpotifyClient{client: &http.Client{}}
}

//TODO: Extract URL Building into testable method
//TODO: Extract Parser functionality into testable methods

func (sClient *SpotifyClient) FindArtist(name string) model.Artist {
	endpoint := fmt.Sprintf("%s://%s/search?type=artist&q=%s", sClient.protocol, sClient.baseUrl, name)
	resp, err := sClient.client.Get(endpoint)
	if err != nil {
		log.Println(fmt.Errorf("GetArtist Error: %v", err))
		return model.Artist{}
	}
	artist := model.Artist{}
	parser, err := simplejson.NewFromReader(resp.Body)
	parser = parser.GetPath("artists", "items").GetIndex(0)
	data, err := parser.Bytes()
	if err != nil {
		log.Println(fmt.Errorf("GetArtist Error: %v", err))
		return model.Artist{}
	}
	json.Unmarshal(data, artist)
	return artist
}

func (sClient *SpotifyClient) getArtistAlbums(artist model.Artist) []string {
	albumEndpoint := fmt.Sprintf("%s://%s/artists/%s/albums", sClient.protocol, sClient.baseUrl, artist.ID)
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

func (sClient *SpotifyClient) GetAssociatedArtists(artist model.Artist) []model.Artist {
	albumIds := sClient.getArtistAlbums(artist)
	albumsEndpoint := fmt.Sprintf("%s://%s/albums/?ids=%s", sClient.protocol, sClient.baseUrl, albumIds[0])
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
		log.Println(fmt.Errorf("GetAssociatedArtists - Album Request Error: %v", err))
		return nil
	}

	albums, err := albumsParser.Get("albums").Array()
	if err != nil {
		log.Println(fmt.Errorf("GetAssociatedArtists - AlbumParser Error: %v", err))
		return nil
	}

	associatedArtists := make([]model.Artist, 0)
	for i := 0; i < len(albums); i++ {
		tracksParser := albumsParser.Get("albums").GetIndex(i).Get("tracks")
		tracks, err := tracksParser.Get("items").Array()
		if err != nil {
			log.Println(fmt.Errorf("GetAssociatedArtists - TrackParser Error: %v", err))
			return nil
		}
		for j := 0; j < len(tracks); j++ {
			trackJSON := tracksParser.Get("items").GetIndex(j)
			trackData, err := trackJSON.Bytes()
			if err != nil {
				log.Println(fmt.Errorf("GetAssociatedArtists - TrackParser Error: %v", err))
				return nil
			}
			track := model.Track{}
			json.Unmarshal(trackData, track)
			artistParser := tracksParser.Get("items").GetIndex(j).Get("artists")
			artists, err := artistParser.Array()
			for k := 0; k < len(artists); k++ {
				performer, err := artistParser.GetIndex(k).Bytes()
				if err != nil {
					log.Println(fmt.Errorf("GetAssociatedArtists - Artist Parser Error: %v", err))
					return nil
				}
				associatedArtist := model.Artist{}
				json.Unmarshal(performer, associatedArtist)
				if associatedArtist.ID != artist.ID {
					associatedArtist.RelatedSong = &track
					associatedArtists = append(associatedArtists, associatedArtist)
				}
			}
		}
	}
	return associatedArtists
}

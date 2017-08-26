package client

import (
	"encoding/json"
	"testing"

	"github.com/bitly/go-simplejson"
)

var sCli = InitSpotifyClient()

var artistJSON = map[string]interface{}{
	"artists": map[string]interface{}{
		"href": "https://api.spotify.com/v1/search?query=tania+bowra&offset=0&limit=20&type=artist",
		"items": []map[string]interface{}{{
			"href":       "https://api.spotify.com/v1/artists/08td7MxkoHQkXnWAYD8d6Q",
			"id":         "08td7MxkoHQkXnWAYD8d6Q",
			"name":       "Tania Bowra",
			"popularity": 0,
			"type":       "artist",
			"uri":        "spotify:artist:08td7MxkoHQkXnWAYD8d6Q",
		}},
		"limit":    20,
		"next":     nil,
		"offset":   0,
		"previous": nil,
		"total":    1,
	},
}

var albumsJSON = map[string]interface{}{
	"albums": []map[string]interface{}{{
		"album_type": "album",
		"artists": []map[string]interface{}{{
			"external_urls": map[string]interface{}{
				"spotify": "https://open.spotify.com/artist/53A0W3U0s8diEn9RhXQhVz",
			},
			"href": "https://api.spotify.com/v1/artists/53A0W3U0s8diEn9RhXQhVz",
			"id":   "53A0W3U0s8diEn9RhXQhVz",
			"name": "Keane",
			"type": "artist",
			"uri":  "spotify:artist:53A0W3U0s8diEn9RhXQhVz",
		}},
		"available_markets": []string{"AD", "AR", "AT", "AU", "BE", "BG", "BO", "BR", "CH", "CL", "CO", "CR", "CY", "CZ", "DE", "DK", "DO", "EC", "EE", "ES", "FI", "FR", "GB", "GR", "GT", "HK", "HN", "HU", "IE", "IS", "IT", "LI", "LT", "LU", "LV", "MC", "MT", "MY", "NI", "NL", "NO", "NZ", "PA", "PE", "PH", "PL", "PT", "PY", "RO", "SE", "SG", "SI", "SK", "SV", "TR", "TW", "UY"},
		"copyrights": []map[string]interface{}{{
			"text": "(C) 2013 Universal Island Records, a division of Universal Music Operations Limited",
			"type": "C",
		}, {
			"text": "(P) 2013 Universal Island Records, a division of Universal Music Operations Limited",
			"type": "P",
		}},
		"external_ids": map[string]interface{}{
			"upc": "00602537518357",
		},
		"external_urls": map[string]interface{}{
			"spotify": "https://open.spotify.com/album/41MnTivkwTO3UUJ8DrqEJJ",
		},
		"href":                   "https://api.spotify.com/v1/albums/41MnTivkwTO3UUJ8DrqEJJ",
		"id":                     "test_id",
		"name":                   "The Best Of Keane (Deluxe Edition)",
		"popularity":             65,
		"release_date":           "2013-11-08",
		"release_date_precision": "day",
		"tracks": map[string]interface{}{
			"href": "https://api.spotify.com/v1/albums/41MnTivkwTO3UUJ8DrqEJJ/tracks?offset=0&limit=50",
			"items": []map[string]interface{}{{
				"artists": []map[string]interface{}{{
					"external_urls": map[string]interface{}{
						"spotify": "https://open.spotify.com/artist/53A0W3U0s8diEn9RhXQhVz",
					},
					"href": "https://api.spotify.com/v1/artists/53A0W3U0s8diEn9RhXQhVz",
					"id":   "test_id",
					"name": "Keane",
					"type": "artist",
					"uri":  "spotify:artist:53A0W3U0s8diEn9RhXQhVz",
				}},
				"available_markets": []string{"AD", "AR", "AT", "AU", "BE", "BG", "BO", "BR", "CH", "CL", "CO", "CR", "CY", "CZ", "DE", "DK", "DO", "EC", "EE", "ES", "FI", "FR", "GB", "GR", "GT", "HK", "HN", "HU", "IE", "IS", "IT", "LI", "LT", "LU", "LV", "MC", "MT", "MY", "NI", "NL", "NO", "NZ", "PA", "PE", "PH", "PL", "PT", "PY", "RO", "SE", "SG", "SI", "SK", "SV", "TR", "TW", "UY"},
				"disc_number":       1,
				"duration_ms":       215986,
				"explicit":          false,
				"external_urls": map[string]interface{}{
					"spotify": "https://open.spotify.com/track/4r9PmSmbAOOWqaGWLf6M9Q",
				},
				"href":         "https://api.spotify.com/v1/tracks/4r9PmSmbAOOWqaGWLf6M9Q",
				"id":           "test_id",
				"name":         "Everybody's Changing",
				"preview_url":  "https://p.scdn.co/mp3-preview/641fd877ee0f42f3713d1649e20a9734cc64b8f9",
				"track_number": 1,
				"type":         "track",
				"uri":          "spotify:track:4r9PmSmbAOOWqaGWLf6M9Q",
			}, {
				"artists": []map[string]interface{}{{
					"external_urls": map[string]interface{}{
						"spotify": "https://open.spotify.com/artist/53A0W3U0s8diEn9RhXQhVz",
					},
					"href": "https://api.spotify.com/v1/artists/53A0W3U0s8diEn9RhXQhVz",
					"id":   "different_test_id",
					"name": "Keane",
					"type": "artist",
					"uri":  "spotify:artist:53A0W3U0s8diEn9RhXQhVz",
				}},
				"available_markets": []string{"AD", "AR", "AT", "AU", "BE", "BG", "BO", "BR", "CH", "CL", "CO", "CR", "CY", "CZ", "DE", "DK", "DO", "EC", "EE", "ES", "FI", "FR", "GB", "GR", "GT", "HK", "HN", "HU", "IE", "IS", "IT", "LI", "LT", "LU", "LV", "MC", "MT", "MY", "NI", "NL", "NO", "NZ", "PA", "PE", "PH", "PL", "PT", "PY", "RO", "SE", "SG", "SI", "SK", "SV", "TR", "TW", "UY"},
				"disc_number":       1,
				"duration_ms":       235880,
				"explicit":          false,
				"external_urls": map[string]interface{}{
					"spotify": "https://open.spotify.com/track/0HJQD8uqX2Bq5HVdLnd3ep",
				},
				"href":         "https://api.spotify.com/v1/tracks/0HJQD8uqX2Bq5HVdLnd3ep",
				"id":           "different_test_id",
				"name":         "Somewhere Only We Know",
				"preview_url":  "https://p.scdn.co/mp3-preview/e001676375ea2b4807cee2f98b51f2f3fe0d109b",
				"track_number": 2,
				"type":         "track",
				"uri":          "spotify:track:0HJQD8uqX2Bq5HVdLnd3ep",
			}},
			"limit":    50,
			"next":     nil,
			"offset":   0,
			"previous": nil,
			"total":    9,
		},
		"type": "album",
		"uri":  "spotify:album:6UXCm6bOO4gFlDQZV5yL37",
	}},
}

func TestBasicEndpointBuildUrl(t *testing.T) {
	expectedURL := "https://api.spotify.com/v1/test"
	actualURL := sCli.buildUrl("test", nil)
	if actualURL != expectedURL {
		t.Errorf("Expected: %s, Actual:%s", expectedURL, actualURL)
	}
}

func TestEmptyStringEndpointBuildUrl(t *testing.T) {
	expectedURL := "https://api.spotify.com/v1/"
	actualURL := sCli.buildUrl("", nil)
	if actualURL != expectedURL {
		t.Errorf("Expected: %s, Actual:%s", expectedURL, actualURL)
	}
}

func TestEndpointWithParamBuildUrl(t *testing.T) {
	expectedURL := "https://api.spotify.com/v1/test?param=value"
	actualURL := sCli.buildUrl("test", map[string]interface{}{"param": "value"})
	if actualURL != expectedURL {
		t.Errorf("Expected: %s, Actual:%s", expectedURL, actualURL)
	}
}

func TestEndpointWithMultipleParamsBuildUrl(t *testing.T) {
	expectedURL := "https://api.spotify.com/v1/test?param=value&type=code"
	actualURL := sCli.buildUrl("test", map[string]interface{}{"param": "value", "type": "code"})
	if actualURL != expectedURL {
		t.Errorf("Expected: %s, Actual:%s", expectedURL, actualURL)
	}
}

func TestParseArtistFromJsonBasic(t *testing.T) {
	jsonBytes, err := json.Marshal(artistJSON)
	parser, err := simplejson.NewJson(jsonBytes)
	if err != nil {
		t.Errorf("Test Initialization Error")
	}
	artist := sCli.parseArtistFromJson(parser)
	if artist.ID != "08td7MxkoHQkXnWAYD8d6Q" {
		t.Errorf("Expected: %s, Actual:%s", "08td7MxkoHQkXnWAYD8d6Q", artist.ID)
	}
	if artist.Name != "Tania Bowra" {
		t.Errorf("Expected: %s, Actual:%s", "Tania Bowra", artist.Name)
	}
	if artist.URL != "https://api.spotify.com/v1/artists/08td7MxkoHQkXnWAYD8d6Q" {
		t.Errorf("Expected: %s, Actual:%s", "https://api.spotify.com/v1/artists/08td7MxkoHQkXnWAYD8d6Q", artist.URL)
	}
}

func TestParseAssociatedArtistsFromJsonBasic(t *testing.T) {
	jsonBytes, err := json.Marshal(albumsJSON)
	sourceID := "test_id"
	parser, err := simplejson.NewJson(jsonBytes)
	if err != nil {
		t.Errorf("Test Initialization Error")
	}
	artists := sCli.parseAssociatedArtistsFromJson(parser, sourceID)
	if len(artists) != 1 {
		t.Errorf("Expected: %d, Actual:%d", 1, len(artists))
	}
	resultArtist := artists[0]
	if resultArtist.ID != "different_test_id" {
		t.Errorf("Expected: %s, Actual:%s", "different_test_id", resultArtist.ID)
	}

}

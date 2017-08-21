package client

import (
	"testing"
)

var sCli = InitSpotifyClient()

func TestBasicEndpointBuildUrl(t *testing.T) {
	expectedUrl := "https://api.spotify.com/v1/test"
	actualUrl := sCli.buildUrl("test", nil)
	if actualUrl != expectedUrl {
		t.Errorf("Expected: %s, Actual:%s", expectedUrl, actualUrl)
	}
}

func TestEmptyStringEndpointBuildUrl(t *testing.T) {
	expectedUrl := "https://api.spotify.com/v1/"
	actualUrl := sCli.buildUrl("", nil)
	if actualUrl != expectedUrl {
		t.Errorf("Expected: %s, Actual:%s", expectedUrl, actualUrl)
	}
}

func TestEndpointWithParamBuildUrl(t *testing.T) {
	expectedUrl := "https://api.spotify.com/v1/test?param=value"
	actualUrl := sCli.buildUrl("test", map[string]interface{}{"param": "value"})
	if actualUrl != expectedUrl {
		t.Errorf("Expected: %s, Actual:%s", expectedUrl, actualUrl)
	}
}

func TestEndpointWithMultipleParamsBuildUrl(t *testing.T) {
	expectedUrl := "https://api.spotify.com/v1/test?param=value&type=code"
	actualUrl := sCli.buildUrl("test", map[string]interface{}{"param": "value", "type": "code"})
	if actualUrl != expectedUrl {
		t.Errorf("Expected: %s, Actual:%s", expectedUrl, actualUrl)
		t.Fail()
	}
}

package model

type Artist struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	URL         string `json:"href"`
	RelatedSong *Track
}

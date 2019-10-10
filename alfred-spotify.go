package main

// run: alfred_workflow_data=workflow alfred_workflow_cache=/tmp/alfred alfred_workflow_bundleid=mk_testing go run alfred-deezer.go track mötley

import (
	"context"
	"log"
	"os"
	"strings"

	aw "github.com/deanishe/awgo"
	spotify "github.com/zmb3/spotify"
	"golang.org/x/oauth2/clientcredentials"
)

var wf *aw.Workflow
var client spotify.Client

func init() {
	wf = aw.New()

	config := &clientcredentials.Config{
		ClientID:     os.Getenv("SPOTIFY_ID"),
		ClientSecret: os.Getenv("SPOTIFY_SECRET"),
		TokenURL:     spotify.TokenURL,
	}
	token, err := config.Token(context.Background())
	if err != nil {
		log.Fatalf("couldn't get token: %v", err)
	}

	client = spotify.Authenticator{}.NewClient(token)
}

func main() {
	wf.Run(run)
}

func run() {
	contentType := os.Args[1]
	title := os.Args[2]

	switch contentType {
	case "track":
		runTracks(title)
	case "album":
		runAlbum(title)
	case "artist":
		runArtist(title)
	}

}

func getLocalURL(url string) string {
	return strings.Replace(url, "https://", "deezer://", -1)
}

func runAlbum(title string) {
	results, err := client.Search(title, spotify.SearchTypeAlbum)
	if err != nil {
		log.Fatal(err)
	}

	for _, album := range results.Albums.Albums {

		id := album.ID.String()
		url := "spotify:album:" + id

		wf.NewItem(album.Artists[0].Name + " - " + album.Name).
			Valid(true).
			Arg(url).
			Quicklook(url).
			UID("album" + id).
			NewModifier("cmd").
			Subtitle("Open in Spotify App").
			Arg(getLocalURL(url))
	}

	wf.SendFeedback()
}
func runArtist(title string) {
	// TODO implement me
}

func runTracks(title string) {

	results, err := client.Search(title, spotify.SearchTypeTrack)
	if err != nil {
		log.Fatal(err)
	}

	for _, track := range results.Tracks.Tracks {

		id := track.ID.String()
		trackURL := "spotify:track:" + id
		albumURL := "spotify:album:" + track.Album.ID.String()

		wf.NewItem(track.Artists[0].Name + " - " + track.Name).
			Subtitle(track.Album.Name).
			Valid(true).
			Arg(trackURL + " " + albumURL).
			Quicklook(trackURL).
			UID("album" + id).
			NewModifier("cmd").
			Subtitle("Open in Spotify App").
			Arg(getLocalURL(trackURL))
	}

	wf.SendFeedback()
}

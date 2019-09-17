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

// aw.Workflow is the main API
var wf *aw.Workflow
var client spotify.Client

func init() {
	// Create a new *Workflow using default configuration
	// (workflow settings are read from the environment variables
	// set by Alfred)
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
	// Wrap your entry point with Run() to catch and log panics and
	// show an error in Alfred instead of silently dying
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
		// var icon aw.Icon
		// icon.Value = album.Images

		id := album.ID.String()
		url := "spotify:album:" + id

		wf.NewItem(album.Artists[0].Name + " - " + album.Name).
			// Subtitle(album.Album.Title).
			Valid(true).
			// Icon(&icon).
			Arg(url).
			Quicklook(url).
			UID("album" + id).
			NewModifier("cmd").
			Subtitle("Open in Deezer App").
			Arg(getLocalURL(url))
	}

	// And send the results to Alfred
	wf.SendFeedback()
}
func runArtist(title string) {
	// TODO implement me
}

func runTracks(title string) {

}

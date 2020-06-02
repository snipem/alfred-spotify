package main

// run: make test

import (
	"context"
	"fmt"
	"log"
	"os"

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

// runAlbum runs the workflow for album searching
func runAlbum(title string) {
	results, err := client.Search(title, spotify.SearchTypeAlbum)
	if err != nil {
		log.Fatal(err)
	}

	for _, album := range results.Albums.Albums {

		id := album.ID.String()
		url := "spotify:album:" + id

		item := wf.NewItem(album.Artists[0].Name + " - " + album.Name).
			Valid(true).
			Arg(url).
			Quicklook(url).
			UID("album" + id)

		item.
			NewModifier("alt").
			Subtitle("Open in Spotify App").
			Arg(url)

		item.
			NewModifier("cmd").
			Subtitle("Open in Spotify Browser").
			Arg(getBrowserURL(album.ID, "album"))

	}

	wf.SendFeedback()
}

func getBrowserURL(id spotify.ID, spotifyType string) string {
	return fmt.Sprintf("https://open.spotify.com/%s/%s", spotifyType, id)
}

// runArtist runs the workflow for artist searching
func runArtist(title string) {
	// TODO implement me
}

// runTracks runs the workflow for track searching
func runTracks(title string) {

	results, err := client.Search(title, spotify.SearchTypeTrack)
	if err != nil {
		log.Fatal(err)
	}

	for _, track := range results.Tracks.Tracks {

		id := track.ID.String()
		trackURL := "spotify:track:" + id
		albumURL := "spotify:album:" + track.Album.ID.String()

		item := wf.NewItem(track.Artists[0].Name + " - " + track.Name).
			Subtitle(track.Album.Name).
			Valid(true).
			Arg(trackURL + " " + albumURL).
			Quicklook(trackURL).
			UID("album" + id)

		item.
			NewModifier("alt").
			Subtitle("Open in Spotify App").
			Arg(trackURL)

		item.
			NewModifier("cmd").
			Subtitle("Open in Spotify Browser").
			Arg(getBrowserURL(track.ID, "track"))
	}

	wf.SendFeedback()
}

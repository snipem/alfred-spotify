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
	case "playlist":
		runPlaylist(title)
	case "artist":
		runArtist(title)
    default:
        fmt.Fprintf(os.Stderr, "Unknown option %s\n", contentType)
        os.Exit(1)
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
		albumURI := "spotify:album:" + id

		item := wf.NewItem(fmt.Sprintf("%s - %s", album.Artists[0].Name, album.Name)).
			Valid(true).
			Arg(albumURI).
			Quicklook(album.Images[0].URL).
			UID(id)

		item.
			NewModifier("alt").
			Subtitle("Open in Spotify App").
			Arg(albumURI)

		item.
			NewModifier("cmd").
			Subtitle("Open in Spotify Browser").
			Arg(getBrowserURL(album.ID, "album"))

	}

	wf.SendFeedback()
}

// runPlaylist runs the workflow for playlist searching
func runPlaylist(title string) {
	results, err := client.Search(title, spotify.SearchTypePlaylist)
	if err != nil {
		log.Fatal(err)
	}

	for _, playlist := range results.Playlists.Playlists {

		id := playlist.ID.String()
		playlistURI := "spotify:playlist:" + id

		item := wf.NewItem(fmt.Sprintf("%s by %s", playlist.Name, playlist.Owner.DisplayName)).
			Valid(true).
			Arg(playlistURI).
			Quicklook(playlistURI).
			UID(id)

		item.
			NewModifier("alt").
			Subtitle("Open in Spotify App").
			Arg(playlistURI)

		item.
			NewModifier("cmd").
			Subtitle("Open in Spotify Browser").
			Arg(getBrowserURL(playlist.ID, "playlist"))

	}

	wf.SendFeedback()
}

func getBrowserURL(id spotify.ID, spotifyType string) string {
	return fmt.Sprintf("https://open.spotify.com/%s/%s", spotifyType, id)
}

// runArtist runs the workflow for artist searching
func runArtist(title string) {
	results, err := client.Search(title, spotify.SearchTypeArtist)
	if err != nil {
		log.Fatal(err)
	}

	for _, artist := range results.Artists.Artists {

		id := artist.ID.String()
		artistURI := "spotify:artist:" + id

		item := wf.NewItem(artist.Name).
			Valid(true).
			Arg(artistURI).
			Quicklook(artistURI).
			UID(id)

		item.
			NewModifier("alt").
			Subtitle("Open in Spotify App").
			Arg(artistURI)

		item.
			NewModifier("cmd").
			Subtitle("Open in Spotify Browser").
			Arg(getBrowserURL(artist.ID, "artist"))

	}

	wf.SendFeedback()
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

		item := wf.NewItem(fmt.Sprintf("%s - %s", track.Artists[0].Name, track.Name)).
			Subtitle(track.Album.Name).
			Valid(true).
			Arg(fmt.Sprintf("%s %s", trackURL, albumURL)).
			Quicklook(track.Album.Images[0].URL).
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

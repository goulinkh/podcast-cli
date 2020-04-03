package main

import (
	"fmt"

	audioplayer "github.com/goulinkh/podcast-cli/audio-player"
	"github.com/goulinkh/podcast-cli/podcasts"
)

var (
	cacheFolder = ".cache/"
)

func main() {
	podcast := podcasts.GetTop50Podcats()[2]
	latestEpisode := podcast.GetEpisodes(0, 1)[0]
	go func() {
		audioURL, _ := latestEpisode.AudioURL()
		fmt.Println(audioplayer.PlaySound(latestEpisode.Title, cacheFolder, audioURL))
	}()
	select {}
}

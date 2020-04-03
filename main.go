package main

import (
	"fmt"

	"github.com/goulinkh/podcast-cli/podcasts"
)

func main() {
	podcast := podcasts.GetTop50Podcats()[0]
	episodes := podcast.GetEpisodes(0, 20)
	fmt.Println(episodes)
}

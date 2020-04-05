package main

import (
	"fmt"
	"os"

	"github.com/akamensky/argparse"
	"github.com/goulinkh/podcast-cli/podcasts"
	"github.com/goulinkh/podcast-cli/ui"
)

func main() {
	parser := argparse.NewParser("podcast-cli", "CLI podcast player")
	podcastSearchQuery := parser.String("s", "search", &argparse.Options{Required: false, Help: "your podcast's name"})
	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
	}
	var podcastsList []*podcasts.Podcast
	if podcastSearchQuery == nil || *podcastSearchQuery == "" {
		podcastsList = podcasts.GetTop50Podcats()
	} else {
		podcastsList = podcasts.FindPodcasts(*podcastSearchQuery)
	}
	ui.NewUI(podcastsList)
	ui.Show()
}

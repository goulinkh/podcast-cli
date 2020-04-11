package main

import (
	"fmt"
	"log"
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
		fmt.Println("Error:", parser.Usage(err))
		return
	}
	var podcastsList []*podcasts.Podcast
	if podcastSearchQuery == nil || *podcastSearchQuery == "" {
		podcastsList, err = podcasts.GetTop50Podcats()
	} else {
		podcastsList, err = podcasts.FindPodcasts(*podcastSearchQuery)
		if len(podcastsList) == 0 {
			fmt.Println("Error: no podcasts found that satisfy the search query")
			return
		}
	}
	if err != nil {
		log.Fatal(err)
	}
	ui.NewUI(podcastsList)
	ui.Show()
}

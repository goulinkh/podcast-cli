package main

import (
	"log"
	"os"

	"github.com/akamensky/argparse"
	ui "github.com/gizak/termui/v3"

	itunesapi "github.com/goulinkh/podcast-cli/itunes-api"
	podcastcliui "github.com/goulinkh/podcast-cli/ui"
)

func main() {

	parser := argparse.NewParser("podcast-cli", "CLI podcast player")
	podcastSearchQuery := parser.String("s", "search", &argparse.Options{Required: false, Help: "your podcast's name"})
	err := parser.Parse(os.Args)
	if err != nil {
		log.Fatalln("Error:", parser.Usage(err))
		return
	}
	err = podcastcliui.InitUI()
	if err != nil {
		log.Fatal(err)
	}

	if podcastSearchQuery == nil || *podcastSearchQuery == "" {
		genres, err := itunesapi.GetGenres()
		if err != nil {
			log.Fatal(err)
		}
		genreWidget := &podcastcliui.GenresUI{
			Genres: genres,
		}
		genreWidget.InitComponents()
		podcastcliui.Show(genreWidget)
	} else {
		// make a search
		podcasts, err := itunesapi.FindPodcasts(*podcastSearchQuery)
		if err != nil {
			log.Fatalln("Error: Failed to search for podcasts")
		}
		podcastsWidget := &podcastcliui.PodcastsUI{Podcasts: podcasts}
		podcastsWidget.InitComponents()
		podcastcliui.Show(podcastsWidget)
	}

	uiEvents := ui.PollEvents()
	for {
		select {
		case e := <-uiEvents:
			cmd, err := podcastcliui.HandleKeyEvent(&e)
			if err != nil {
				log.Fatal(err)
			}
			if cmd == podcastcliui.Exit {
				return
			}
		}
	}
}

package main

import (
	"log"
	"os"

	"github.com/akamensky/argparse"
	ui "github.com/gizak/termui/v3"

	itunesapi "github.com/goulinkh/podcast-cli/itunes-api"
	"github.com/goulinkh/podcast-cli/rss"
	podcastcliui "github.com/goulinkh/podcast-cli/ui"
)

func main() {

	parser := argparse.NewParser("podcast-cli", "CLI podcast player")
	podcastSearchQuery := parser.String("s", "search", &argparse.Options{Required: false, Help: "your podcast's name"})
	rssUrl := parser.String("r", "rss", &argparse.Options{Required: false, Help: "custom podcast rss source"})
	offset := parser.Int("o", "offset", &argparse.Options{Required: false, Help: "play episode number"})
	err := parser.Parse(os.Args)
	if err != nil {
		log.Fatalln("Error:", parser.Usage(err))
		return
	}
	err = podcastcliui.InitUI()
	if err != nil {
		log.Fatal(err)
	}

	if podcastSearchQuery != nil && *podcastSearchQuery != "" {
		podcasts, err := itunesapi.FindPodcasts(*podcastSearchQuery)
		if err != nil {
			log.Fatalln("Error: Failed to search for podcasts")
		}
		podcastsWidget := &podcastcliui.PodcastsUI{Podcasts: podcasts}
		podcastsWidget.InitComponents()
		podcastcliui.Show(podcastsWidget)
	} else if rssUrl != nil && *rssUrl != "" {
		episodes, err := rss.ParseEpisodes(*rssUrl)
		if err != nil {
			log.Fatalln("Error: Failed to get episodes from the url: " + *rssUrl)
		}
		episodesWidget := &podcastcliui.EpisodesUI{Episodes: episodes}
		episodesWidget.InitComponents()
		podcastcliui.Show(episodesWidget)
		if offset != nil {
			episodesWidget.Play(*offset)
		}
	} else {
		genres, err := itunesapi.GetGenres()
		if err != nil {
			log.Fatal(err)
		}
		genreWidget := &podcastcliui.GenresUI{
			Genres: genres,
		}
		genreWidget.InitComponents()
		podcastcliui.Show(genreWidget)

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

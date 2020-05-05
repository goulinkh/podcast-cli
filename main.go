package main

import (
	"log"

	ui "github.com/gizak/termui/v3"

	itunesapi "github.com/goulinkh/podcast-cli/itunes-api"
	newui "github.com/goulinkh/podcast-cli/ui"
)

func main() {
	genres, err := itunesapi.GetGenres()
	if err != nil {
		log.Fatal(err)
	}
	err = newui.InitUI()
	if err != nil {
		log.Fatal(err)
	}
	genreWidget := &newui.GenresUI{
		Genres: genres,
	}
	genreWidget.InitComponents()
	newui.Show(genreWidget)

	uiEvents := ui.PollEvents()
	for {
		select {
		case e := <-uiEvents:
			cmd, err := newui.HandleKeyEvent(&e)
			if err != nil {
				log.Fatal(err)
			}
			if cmd == newui.Exit {
				return
			}
		}
	}
	// parser := argparse.NewParser("podcast-cli", "CLI podcast player")
	// podcastSearchQuery := parser.String("s", "search", &argparse.Options{Required: false, Help: "your podcast's name"})
	// err := parser.Parse(os.Args)
	// if err != nil {
	// 	fmt.Println("Error:", parser.Usage(err))
	// 	return
	// }
	// var podcastsList []*podcasts.Podcast
	// if podcastSearchQuery == nil || *podcastSearchQuery == "" {
	// 	podcastsList, err = podcasts.GetTop50Podcats()
	// } else {
	// 	podcastsList, err = podcasts.FindPodcasts(*podcastSearchQuery)
	// 	if len(podcastsList) == 0 {
	// 		fmt.Println("Error: no podcasts found that satisfy the search query")
	// 		return
	// 	}
	// }
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// ui.NewUI(podcastsList)
	// ui.Show()
}

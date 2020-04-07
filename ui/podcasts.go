package ui

import (
	"fmt"
	"strings"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/goulinkh/podcast-cli/podcasts"
)

var (
	podcastsListWidget   *widgets.List
	podcastDetailsWidget *widgets.Paragraph

	podcastsList []*podcasts.Podcast
)

func initPodcastsUI(podcasts []*podcasts.Podcast) {
	podcastsList = podcasts
	podcastsListWidget = widgets.NewList()
	podcastsListWidget.Title = "Podcasts"
	podcastsListWidget.TextStyle = ui.NewStyle(ui.ColorWhite)
	podcastsListWidget.SelectedRowStyle = ui.NewStyle(ui.ColorMagenta)
	podcastsListWidget.BorderStyle.Fg = ui.ColorMagenta

	podcastDetailsWidget = widgets.NewParagraph()
	podcastDetailsWidget.Title = "Details"
	podcastDetailsWidget.BorderStyle.Fg = ui.ColorMagenta
	podcastDetailsWidget.BorderLeft = false
	podcastDetailsWidget.BorderBottom = false

	updatePodcastsList()
}
func updatePodcastsList() {
	podcastsListWidget.Rows = make([]string, 0)
	for _, podcast := range podcastsList {
		podcastsListWidget.Rows = append(podcastsListWidget.Rows, podcast.Title)
	}
}

func onPodcastItemClick(index int) {
	episodesList = podcastsList[index].Episodes(0, 100)
	updateEpisodesList()
	currentListWidget = episodesListWidget
	episodesListWidget.SelectedRow = 0
	currentDetailsWidget = episodeDetailsWidget
	updateDetailsWidget()
	initGrid()
	frameUpdate()

}
func updatePodcastDetails() {
	currentPodcast := podcastsList[currentListWidget.SelectedRow]
	title := fmt.Sprintf("[Title](fg:magenta)        %s", currentPodcast.Title)
	description := fmt.Sprintf("[Description](fg:magenta)  %s", currentPodcast.Description)
	author := fmt.Sprintf("[Author](fg:magenta)       %s", currentPodcast.Author)
	podcastDetailsWidget.Text = strings.Join([]string{title, description, author}, "\n")
}

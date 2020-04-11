package ui

import (
	"fmt"
	"strings"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/goulinkh/podcast-cli/podcasts"
)

// TODO: make a column for podcast with podcast image instead of audio player
var (
	episodesListWidget   *widgets.List
	episodesList         []*podcasts.Episode
	episodeDetailsWidget *widgets.Paragraph
)

func initEpisodesUI() {
	//initEpisodesUI
	episodesListWidget = widgets.NewList()
	episodesListWidget.Title = "Episodes"
	episodesListWidget.SelectedRowStyle.Modifier = ui.ModifierClear
	episodesListWidget.TextStyle = ui.NewStyle(ui.ColorWhite)
	episodesListWidget.SelectedRowStyle = ui.NewStyle(ui.ColorMagenta)
	episodesListWidget.BorderStyle.Fg = ui.ColorMagenta

	episodeDetailsWidget = widgets.NewParagraph()
	episodeDetailsWidget.Title = "Details"
	episodeDetailsWidget.BorderStyle.Fg = ui.ColorMagenta
	episodeDetailsWidget.BorderLeft = false
	episodeDetailsWidget.BorderBottom = false
}

func updateEpisodesList() {
	episodesListWidget.Rows = make([]string, len(episodesList))
	for i, episode := range episodesList {
		episodesListWidget.Rows[i] = episode.Title
	}
}

func updateEpisodesDetails() {
	currentEpisode := episodesList[currentListWidget.SelectedRow]
	title := fmt.Sprintf("[Title](fg:magenta)        %s", currentEpisode.Title)
	description := fmt.Sprintf("[Description](fg:magenta)  %s", currentEpisode.Description)
	date := fmt.Sprintf("[Release Date](fg:magenta) %s", currentEpisode.ReleaseDate)
	episodeDetailsWidget.Text = strings.Join([]string{title, description, date}, "\n")
}

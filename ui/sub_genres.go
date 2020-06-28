package ui

import (
	"errors"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	itunesapi "github.com/goulinkh/podcast-cli/itunes-api"
)

type SubGenresUI struct {
	Genres     []*itunesapi.Genre
	gridWidget *ui.Grid
	listWidget *widgets.List
}

func (g *SubGenresUI) InitComponents() error {
	g.newGenresListWidget()
	err := g.newGridWidget()
	if err != nil {
		return err
	}
	return nil
}
func (g *SubGenresUI) MainUI() *ui.Grid {
	return g.gridWidget
}
func (g *SubGenresUI) HandleEvent(event *ui.Event) (Command, error) {
	switch event.ID {
	case "j", "<Down>":
		g.listWidget.ScrollDown()

	case "k", "<Up>":
		g.listWidget.ScrollUp()
	case "<Enter>":
		podcasts, err := g.Genres[g.listWidget.SelectedRow].GetPodcasts()
		if err != nil {
			return Nothing, err
		}
		podcastsUI := &PodcastsUI{Podcasts: podcasts}
		err = podcastsUI.InitComponents()
		if err != nil {
			return Nothing, err
		}
		Show(podcastsUI)
	}
	return Nothing, nil
}
func (g *SubGenresUI) newGenresListWidget() error {
	g.listWidget = widgets.NewList()
	g.listWidget.Title = "Select a Sub Genre"
	g.listWidget.TextStyle = ui.NewStyle(FgColor)
	g.listWidget.SelectedRowStyle.Fg = ui.ColorBlack
	g.listWidget.SelectedRowStyle.Bg = AccentColor
	g.listWidget.BorderStyle.Fg = AccentColor
	if g.Genres == nil {
		return errors.New("Missing Sub Genres array")
	}
	g.listWidget.Rows = make([]string, len(g.Genres))
	for i, genre := range g.Genres {
		g.listWidget.Rows[i] = genre.Text
	}
	return nil
}
func (g *SubGenresUI) newGridWidget() error {
	if g.listWidget == nil {
		return errors.New("Uninitialized sub genres list widget")
	}
	g.gridWidget = ui.NewGrid()
	termWidth, termHeight := ui.TerminalDimensions()
	g.gridWidget.SetRect(0, 0, termWidth, termHeight-1)
	placeholder := ui.NewBlock()
	placeholder.BorderBottom = false
	placeholder.BorderLeft = false
	placeholder.BorderStyle.Fg = AccentColor
	g.gridWidget.Set(
		ui.NewRow(1.0,
			ui.NewCol(1.0/2, g.listWidget),
			ui.NewCol(1.0/2,
				ui.NewRow(6.0/10, placeholder),
				ui.NewRow(4.0/10, audioPlayerWidget.MainUI()))))
	return nil
}
func (g *SubGenresUI) refreshComponents() {
	g.newGenresListWidget()
}

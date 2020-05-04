package newui

import (
	"errors"
	"fmt"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/goulinkh/podcast-cli/temp"
)

type GenresUI struct {
	Genres     []*temp.Genre
	gridWidget *ui.Grid
	listWidget *widgets.List
	inSubGenre bool
}

func (g *GenresUI) InitComponents() error {
	g.inSubGenre = false
	g.newGenresListWidget()
	err := g.newGridWidget()
	if err != nil {
		return err
	}
	return nil
}
func (g GenresUI) MainUI() *ui.Grid {
	return g.gridWidget
}
func (g *GenresUI) HandleEvent(event *ui.Event) (Command, error) {
	switch event.ID {
	case "j", "<Down>":
		g.listWidget.ScrollDown()

	case "k", "<Up>":
		g.listWidget.ScrollUp()
	case "<Enter>":
		if g.inSubGenre {
			podcasts, err := g.Genres[g.listWidget.SelectedRow].GetPodcasts()
			if err != nil {
				return Nothing, err
			}
			podcastsUI := &PodcastsUI{
				Podcasts: podcasts,
			}
			err = podcastsUI.InitComponents()
			fmt.Println(podcastsUI.MainUI())
			if err != nil {
				return Nothing, err
			}
			Show(podcastsUI)
		} else {
			g.inSubGenre = true
			g.Genres = g.Genres[g.listWidget.SelectedRow].SubGenre
			g.refreshComponents()
		}
	}
	return Nothing, nil
}
func (g *GenresUI) newGenresListWidget() error {
	g.listWidget = widgets.NewList()
	g.listWidget.Title = "Select a Genre"
	g.listWidget.TextStyle = ui.NewStyle(FgColor)
	g.listWidget.SelectedRowStyle = ui.NewStyle(AccentColor)
	g.listWidget.BorderStyle.Fg = AccentColor
	if g.Genres == nil {
		return errors.New("Missing Genres array")
	}
	g.listWidget.Rows = make([]string, len(g.Genres))
	for i, genre := range g.Genres {
		g.listWidget.Rows[i] = genre.Text
	}
	return nil
}
func (g *GenresUI) newGridWidget() error {
	if g.listWidget == nil {
		return errors.New("Uninitialized genres list widget")
	}
	g.gridWidget = ui.NewGrid()
	termWidth, termHeight := ui.TerminalDimensions()
	g.gridWidget.SetRect(0, 0, termWidth, termHeight-1)
	g.gridWidget.Set(
		ui.NewRow(1.0,
			ui.NewCol(1.0, g.listWidget)))
	return nil
}
func (g *GenresUI) refreshComponents() {
	g.listWidget.Rows = make([]string, len(g.Genres))
	for i, genre := range g.Genres {
		g.listWidget.Rows[i] = genre.Text
	}
	termWidth, termHeight := ui.TerminalDimensions()
	g.gridWidget.SetRect(0, 0, termWidth, termHeight-1)
	g.gridWidget.Set(
		ui.NewRow(1.0,
			ui.NewCol(1.0, g.listWidget)))
}

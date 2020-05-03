package newui

import (
	"errors"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/goulinkh/podcast-cli/temp"
)

type GenresUI struct {
	Genres           []*temp.Genre
	gridWidget       *ui.Grid
	genresListWidget *widgets.List
}

func (g *GenresUI) MainUI() *ui.Grid {
	return g.gridWidget
}
func (g *GenresUI) HandleEvent(event *ui.Event) {
	switch event.ID {
	case "j", "<Down>":
		g.genresListWidget.ScrollDown()

	case "k", "<Up>":
		g.genresListWidget.ScrollUp()
	}
}
func (g *GenresUI) InitComponents() error {
	g.newGenresListWidget()
	err := g.newGridWidget()
	if err != nil {
		return err
	}
	return nil
}
func (g *GenresUI) newGenresListWidget() error {
	g.genresListWidget = widgets.NewList()
	g.genresListWidget.Title = "Select a Genre"
	g.genresListWidget.TextStyle = ui.NewStyle(FgColor)
	g.genresListWidget.SelectedRowStyle = ui.NewStyle(AccentColor)
	g.genresListWidget.BorderStyle.Fg = AccentColor
	if g.Genres == nil {
		return errors.New("Missing Genres array")
	}
	g.genresListWidget.Rows = make([]string, len(g.Genres))
	for i, genre := range g.Genres {
		g.genresListWidget.Rows[i] = genre.Text
	}
	return nil
}
func (g *GenresUI) newGridWidget() error {
	if g.genresListWidget == nil {
		return errors.New("Uninitialized genres list widget")
	}
	g.gridWidget = ui.NewGrid()
	termWidth, termHeight := ui.TerminalDimensions()
	g.gridWidget.SetRect(0, 0, termWidth, termHeight-1)
	g.gridWidget.Set(
		ui.NewRow(1.0,
			ui.NewCol(1.0, g.genresListWidget)))
	return nil
}

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
}

func (g *GenresUI) InitComponents() error {
	g.newGenresListWidget()
	err := g.newGridWidget()
	if err != nil {
		return err
	}
	return nil
}
func (g *GenresUI) MainUI() *ui.Grid {
	return g.gridWidget
}
func (g *GenresUI) HandleEvent(event *ui.Event) (Command, error) {
	switch event.ID {
	case "j", "<Down>":
		g.listWidget.ScrollDown()

	case "k", "<Up>":
		g.listWidget.ScrollUp()
	case "<Enter>":
		subGenres := g.Genres[g.listWidget.SelectedRow].SubGenre
		var subGenreUI *SubGenresUI
		if subGenres == nil || len(subGenres) == 0 {
			subGenreUI = &SubGenresUI{Genres: []*temp.Genre{g.Genres[g.listWidget.SelectedRow]}}
		} else {
			subGenreUI = &SubGenresUI{Genres: g.Genres[g.listWidget.SelectedRow].SubGenre}
		}
		subGenreUI.InitComponents()
		Show(subGenreUI)
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
	fmt.Println(audioPlayerWidget)
	g.gridWidget.Set(
		ui.NewRow(1.0,
			ui.NewCol(1.0/2, g.listWidget),
			ui.NewCol(1.0/2, audioPlayerWidget.MainUI())))
	return nil
}

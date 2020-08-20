package ui

import (
	"errors"
	"fmt"
	"strings"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	itunesapi "github.com/goulinkh/podcast-cli/itunes-api"
)

type EpisodesUI struct {
	Episodes      []*itunesapi.Episode
	listWidget    *widgets.List
	detailsWidget *widgets.Paragraph
	gridWidget    *ui.Grid
}

func (e *EpisodesUI) InitComponents() error {
	e.initListWidget()
	e.initDetailsWidget()
	err := e.initGridWidget()
	return err
}

func (e *EpisodesUI) MainUI() *ui.Grid {
	return e.gridWidget
}

func (e *EpisodesUI) HandleEvent(event *ui.Event) (Command, error) {
	switch event.ID {
	case "j", "<Down>":
		e.listWidget.ScrollDown()
		e.updateDetailsWidget()
	case "k", "<Up>":
		e.listWidget.ScrollUp()
		e.updateDetailsWidget()

	case "<Enter>":
		audioPlayerWidget.Play(e.Episodes, e.listWidget.SelectedRow)
	}
	return Nothing, nil
}
func (e *EpisodesUI) Play(index int) {
	audioPlayerWidget.Play(e.Episodes, index)

}
func (e *EpisodesUI) initGridWidget() error {
	if e.listWidget == nil {
		return errors.New("List widget is not initialized")
	}
	if e.detailsWidget == nil {
		return errors.New("Details widget is not initialized")
	}
	e.gridWidget = ui.NewGrid()
	termWidth, termHeight := ui.TerminalDimensions()
	e.gridWidget.SetRect(0, 0, termWidth, termHeight-1)
	e.gridWidget.Set(
		ui.NewRow(1.0,
			ui.NewCol(1.0/2, e.listWidget),
			ui.NewCol(1.0/2,
				ui.NewRow(6.0/10, e.detailsWidget),
				ui.NewRow(4.0/10, audioPlayerWidget.MainUI()))))
	return nil
}

func (e *EpisodesUI) initListWidget() {
	e.listWidget = widgets.NewList()
	e.listWidget.Title = "Episodes"
	e.listWidget.SelectedRowStyle.Modifier = ui.ModifierClear
	e.listWidget.TextStyle.Fg = FgColor
	e.listWidget.SelectedRowStyle.Fg = ui.ColorBlack
	e.listWidget.SelectedRowStyle.Bg = AccentColor
	e.listWidget.BorderStyle.Fg = AccentColor
	e.listWidget.Rows = make([]string, len(e.Episodes))
	for i, episode := range e.Episodes {
		e.listWidget.Rows[i] = episode.Title
	}
}
func (e *EpisodesUI) initDetailsWidget() {
	e.detailsWidget = widgets.NewParagraph()
	e.detailsWidget.Title = "Details"
	e.detailsWidget.BorderStyle.Fg = AccentColor
	e.detailsWidget.BorderLeft = false
	e.detailsWidget.BorderBottom = false
	e.updateDetailsWidget()
}
func (e *EpisodesUI) updateDetailsWidget() {
	if e.Episodes == nil || len(e.Episodes) == 0 {
		return
	}
	currentEpisode := e.Episodes[e.listWidget.SelectedRow]
	title := fmt.Sprintf("[Title](fg:magenta)        %s", currentEpisode.Title)
	description := fmt.Sprintf("[Description](fg:magenta)  %s", currentEpisode.Description)
	date := fmt.Sprintf("[Release Date](fg:magenta) %s", currentEpisode.ReleaseDate)
	duration := fmt.Sprintf("[Duration](fg:magenta)     %d min", currentEpisode.DurationInMilliseconds/60000)
	e.detailsWidget.Text = strings.Join([]string{title, duration, date, description}, "\n")
}

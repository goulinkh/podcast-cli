package newui

import (
	"errors"
	"fmt"
	"strings"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/goulinkh/podcast-cli/temp"
)

type PodcastsUI struct {
	Podcasts      []*temp.Podcast
	listWidget    *widgets.List
	detailsWidget *widgets.Paragraph
	gridWidget    *ui.Grid
}

func (p *PodcastsUI) InitComponents() error {
	p.initListWidget()
	p.initDetailsWidget()
	err := p.initGridWidget()
	return err
}
func (p PodcastsUI) MainUI() *ui.Grid {
	return p.gridWidget
}
func (p *PodcastsUI) HandleEvent(event *ui.Event) (Command, error) {
	switch event.ID {
	case "j", "<Down>":
		p.listWidget.ScrollDown()
		p.updateDetailsWidget()
	case "k", "<Up>":
		p.listWidget.ScrollUp()
		p.updateDetailsWidget()

	case "<Enter>":
		episodes, err := p.Podcasts[p.listWidget.SelectedRow].GetEpisodes()
		if err != nil {
			return Nothing, err
		}
		episodesUI := &EpisodesUI{Episodes: episodes}
		err = episodesUI.InitComponents()
		if err != nil {
			return Nothing, err
		}
		Show(episodesUI)
	}
	return Nothing, nil
}
func (p *PodcastsUI) initGridWidget() error {
	if p.listWidget == nil {
		return errors.New("Uninitialized podcasts list widget")
	}
	if p.detailsWidget == nil {
		return errors.New("Uninitialized details widget")
	}
	p.gridWidget = ui.NewGrid()
	termWidth, termHeight := ui.TerminalDimensions()
	p.gridWidget.SetRect(0, 0, termWidth, termHeight-1)
	p.gridWidget.Set(
		ui.NewRow(1.0,
			ui.NewCol(1.0/2, p.listWidget),
			ui.NewCol(1.0/2, ui.NewRow(5.0/7, p.detailsWidget))))
	return nil
}
func (p *PodcastsUI) initListWidget() {
	p.listWidget = widgets.NewList()
	p.listWidget.Title = "Podcasts List"
	p.listWidget.TextStyle.Fg = FgColor
	p.listWidget.SelectedRowStyle.Fg = AccentColor
	p.listWidget.BorderStyle.Fg = AccentColor
	p.listWidget.Rows = make([]string, len(p.Podcasts))
	for i, podcast := range p.Podcasts {
		p.listWidget.Rows[i] = podcast.Title
	}
}
func (p *PodcastsUI) initDetailsWidget() {
	p.detailsWidget = widgets.NewParagraph()
	p.detailsWidget.Title = "Details"
	p.detailsWidget.BorderStyle.Fg = AccentColor
	p.detailsWidget.BorderLeft = false
	p.detailsWidget.BorderBottom = false
	p.updateDetailsWidget()
}
func (p *PodcastsUI) updateDetailsWidget() {
	currentPodcast := p.Podcasts[p.listWidget.SelectedRow]
	title := fmt.Sprintf("[Title](fg:magenta)        %s", currentPodcast.Title)
	description := fmt.Sprintf("[Description](fg:magenta)  %s", currentPodcast.Description)
	author := fmt.Sprintf("[Author](fg:magenta)       %s", currentPodcast.Author)
	p.detailsWidget.Text = strings.Join([]string{title, description, author}, "\n")
}

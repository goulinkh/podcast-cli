package ui

import (
	"fmt"
	"log"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	audioplayer "github.com/goulinkh/podcast-cli/audio-player"
	"github.com/goulinkh/podcast-cli/config"
	"github.com/goulinkh/podcast-cli/podcasts"
)

var (
	grid                  *ui.Grid
	currentListWidget     *widgets.List
	currentDetailsWidget  *widgets.Paragraph
	helpBarWidget         *widgets.Paragraph
	audioDuration         int
	audioDurationWidget   *widgets.Gauge
	nowPlayingWidget      *widgets.Paragraph
	currentPlayingEpisode *podcasts.Episode
)

func initGrid() {
	grid = ui.NewGrid()
	termWidth, termHeight := ui.TerminalDimensions()
	grid.SetRect(0, 0, termWidth, termHeight-1)
	grid.Set(
		ui.NewRow(1.0,
			ui.NewCol(1.0/2, currentListWidget),
			ui.NewCol(1.0/2, ui.NewRow(5.0/7, currentDetailsWidget),
				ui.NewRow(1.0/7, nowPlayingWidget),
				ui.NewRow(1.0/7, audioDurationWidget))))
	ui.Clear()
	ui.Render(grid, helpBarWidget)
}

func initHelpBar() {
	helpBarWidget = widgets.NewParagraph()
	helpBarWidget.Text = "[ Enter ](fg:black)[Select](fg:black,bg:green) " +
		"[ p, Space ](fg:black)[Play/Pause](fg:black,bg:green) " +
		"[Esc ](fg:black)[Back](fg:black,bg:green) " +
		"[Right ](fg:black)[+10s](fg:black,bg:green) " +
		"[Left ](fg:black)[-10s](fg:black,bg:green) " +
		"[ q ](fg:black)[Exit](fg:black,bg:green)" +
		"[ s ](fg:black)[SEARCH](fg:black,bg:green)"

	helpBarWidget.Border = false
	helpBarWidget.WrapText = true
	helpBarWidget.TextStyle = ui.Style{Modifier: ui.ModifierBold, Bg: ui.ColorWhite}
}

func NewUI(podcasts []*podcasts.Podcast) {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize the UI: %v", err)
	}

	initPodcastsUI(podcasts)
	initEpisodesUI()

	currentDetailsWidget = podcastDetailsWidget
	currentListWidget = podcastsListWidget
	updateDetailsWidget()
	initHelpBar()

	audioDurationWidget = widgets.NewGauge()
	audioDurationWidget.BorderLeft = false
	audioDurationWidget.BarColor = ui.ColorMagenta
	audioDurationWidget.BorderStyle.Fg = ui.ColorMagenta
	nowPlayingWidget = widgets.NewParagraph()
	nowPlayingWidget.BorderLeft = false
	nowPlayingWidget.BorderBottom = false
	nowPlayingWidget.TextStyle.Fg = ui.ColorMagenta
	nowPlayingWidget.Title = "NOW PLAYING"
	nowPlayingWidget.TitleStyle = ui.Style{Modifier: ui.ModifierBold}
	nowPlayingWidget.BorderStyle.Fg = ui.ColorMagenta
	termWidth, termHeight := ui.TerminalDimensions()
	helpBarWidget.SetRect(0, termHeight-1, termWidth, termHeight)
}

func Show() {

	defer ui.Close()
	initGrid()
	frameUpdate()

	uiEvents := ui.PollEvents()
	for {
		select {
		case <-time.After(time.Millisecond * 500):
			if audioplayer.MainCtrl != nil {
				if audioplayer.MainCtrl.Paused {
					audioDurationWidget.Title = "Stopped"
				} else {
					position := audioplayer.Position()
					audioDurationWidget.Title = "Running"
					audioDurationWidget.Label = fmt.Sprintf("%d:%d", position/60, position%60)
					if audioDuration > 0 {
						audioDurationWidget.Percent = (position * 100) / audioDuration
					}
				}
				rerender()
			}
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				return
			case "<Enter>":
				if currentListWidget == podcastsListWidget {
					onPodcastItemClick(currentListWidget.SelectedRow)
				} else if currentListWidget == episodesListWidget {
					episode := episodesList[episodesListWidget.SelectedRow]
					if episode != currentPlayingEpisode {
						currentPlayingEpisode = episode
						nowPlayingWidget.Text = currentPlayingEpisode.Title
						url, err := episode.AudioURL()
						audioplayer.MainCtrl = nil
						if err != nil {
							audioDurationWidget.Title = "Failed to fetch audio"
						} else {
							audioDurationWidget.Title = "Downloading audio ..."
							go func() {
								audioDuration, err = audioplayer.PlaySound(episode.Title, config.CachePath, url)
								if err != nil {
									audioDurationWidget.Title = "Unsupported audio content"
								}
								rerender()
							}()
						}
					}
				}
			case "p", "<Space>":
				if audioplayer.MainCtrl != nil {

					if audioplayer.MainCtrl.Paused {
						audioplayer.PauseSong(false)
						audioDurationWidget.Title = "Running"
					} else {
						audioplayer.PauseSong(true)
						audioDurationWidget.Title = "Stopped"
					}
				}
			case "<Right>":
				if audioplayer.MainCtrl != nil {
					position := audioplayer.Position() + 10

					if position < audioDuration {
						audioplayer.Seek(position)
					}
				}
			case "<Left>":
				if audioplayer.MainCtrl != nil {
					position := audioplayer.Position() - 10
					if position > 0 {
						audioplayer.Seek(position)
					}
				}
			case "<Resize>":
				payload := e.Payload.(ui.Resize)
				grid.SetRect(0, 0, payload.Width, payload.Height-1)
				helpBarWidget.SetRect(0, payload.Height-1, payload.Width, payload.Height)

			case "j", "<Down>":
				currentListWidget.ScrollDown()
				updateDetailsWidget()

			case "k", "<Up>":
				currentListWidget.ScrollUp()
				updateDetailsWidget()

			case "<C-d>":
				currentListWidget.ScrollHalfPageDown()
				updateDetailsWidget()

			case "<C-u>":
				currentListWidget.ScrollHalfPageUp()
				updateDetailsWidget()

			case "<C-f>":
				currentListWidget.ScrollPageDown()
				updateDetailsWidget()

			case "<C-b>":
				currentListWidget.ScrollPageUp()
				updateDetailsWidget()

			case "<Escape>", "<C-<Backspace>>", "<Backspace>":
				currentListWidget = podcastsListWidget
				currentDetailsWidget = podcastDetailsWidget
				updateDetailsWidget()
				frameUpdate()
			}
			rerender()
		}
	}
}

func rerender() {
	ui.Clear()
	ui.Render(grid, helpBarWidget)
}

func frameUpdate() {
	initGrid()
}

func updateDetailsWidget() {
	if currentListWidget == podcastsListWidget {
		updatePodcastDetails()
	} else if currentListWidget == episodesListWidget {
		updateEpisodesDetails()
	}
}

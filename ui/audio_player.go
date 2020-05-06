package ui

import (
	"fmt"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	audioplayer "github.com/goulinkh/podcast-cli/audio-player"
	itunesapi "github.com/goulinkh/podcast-cli/itunes-api"
)

type AudioPlayerWidget struct {
	nowPlaying          *itunesapi.Episode
	paused              bool
	audioPositionWidget *widgets.Gauge
	playerStatusWidget  *widgets.Paragraph
	grid                *ui.Grid
}

func (ap *AudioPlayerWidget) InitComponents() {
	ap.paused = true
	ap.initAudipPositionWidget()
	ap.initPlayerStatusWidget()
	ap.initGrid()
}

func (ap *AudioPlayerWidget) MainUI() *ui.Grid {
	return ap.grid
}
func (ap *AudioPlayerWidget) HandleEvent(e *ui.Event) (Command, error) {
	switch e.ID {
	case "p", "<Space>":
		ap.Pause()
	case "<Right>":
		if audioplayer.MainCtrl != nil && ap.nowPlaying != nil {
			position := audioplayer.Position() + 10

			if position < ap.nowPlaying.DurationInMilliseconds/1000 {
				audioplayer.Seek(position)
			}
		}
	case "<Left>":
		if audioplayer.MainCtrl != nil && ap.nowPlaying != nil {
			position := audioplayer.Position() - 10

			if position > 0 {
				audioplayer.Seek(position)
			}
		}
	}
	return Nothing, nil
}

func (ap *AudioPlayerWidget) Play(e *itunesapi.Episode) {
	if ap.nowPlaying != nil && ap.nowPlaying.Id == e.Id {
		return
	}
	if ap.nowPlaying == nil {
		go func() {
			for {
				select {
				case <-time.After(time.Millisecond * 100):
					if ap.paused {
						ap.audioPositionWidget.Title = "Paused"
					} else {
						ap.playerStatusWidget.Text = ap.nowPlaying.Title
						position := audioplayer.Position()
						ap.audioPositionWidget.Title = "Running"
						ap.audioPositionWidget.Label = fmt.Sprintf("%d:%d", position/60, position%60)
						audioDuration := e.DurationInMilliseconds / 1000
						if audioDuration > 0 {
							ap.audioPositionWidget.Percent = (position * 100) / audioDuration
						}
					}
					RefreshUI()
				}
			}
		}()
	}
	ap.playerStatusWidget.Title = "Downloading audio ..."
	RefreshUI()
	go func() {
		if err := audioplayer.PlaySound(e); err != nil {
			ap.playerStatusWidget.Title = "Failed to play audio"
			RefreshUI()
		}
		ap.nowPlaying = e
		ap.paused = false
		ap.playerStatusWidget.Title = "Now Playing"
	}()
	return
}

func (ap *AudioPlayerWidget) Pause() {
	ap.paused = !ap.paused
	audioplayer.PauseSong(ap.paused)
}

func (ap *AudioPlayerWidget) initAudipPositionWidget() {
	ap.audioPositionWidget = widgets.NewGauge()
	ap.audioPositionWidget.BorderLeft = false
	ap.audioPositionWidget.BarColor = AccentColor
	ap.audioPositionWidget.BorderStyle.Fg = AccentColor
}
func (ap *AudioPlayerWidget) initPlayerStatusWidget() {
	ap.playerStatusWidget = widgets.NewParagraph()
	ap.playerStatusWidget.BorderLeft = false
	ap.playerStatusWidget.BorderBottom = false
	ap.playerStatusWidget.TextStyle.Fg = AccentColor
	ap.playerStatusWidget.Title = "Now Playing"
	ap.playerStatusWidget.TitleStyle.Fg = FgColor
	ap.playerStatusWidget.BorderStyle.Fg = AccentColor
}
func (ap *AudioPlayerWidget) initGrid() {
	ap.grid = ui.NewGrid()
	ap.grid.Border = false
	ap.grid.Set(
		ui.NewRow(1.0/2, ap.playerStatusWidget),
		ui.NewRow(1.0/2, ap.audioPositionWidget),
	)
}

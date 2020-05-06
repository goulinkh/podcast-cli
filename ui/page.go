package ui

import (
	"errors"
	"fmt"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

var (
	FgColor           = ui.ColorWhite
	AccentColor       = ui.ColorMagenta
	pagesHistory      = make([]Page, 0)
	currentPage       Page
	helpBarWidget     *widgets.Paragraph
	audioPlayerWidget = &AudioPlayerWidget{}
)

type Page interface {
	MainUI() *ui.Grid
	HandleEvent(*ui.Event) (Command, error)
}

func InitUI() error {
	if err := ui.Init(); err != nil {
		errors.New(fmt.Sprintf("failed to initialize the UI: %v", err))
	}
	helpBarWidget = newHelpBarWidget()
	audioPlayerWidget.InitComponents()
	return nil
}
func Show(p Page) {
	if currentPage != nil {
		pagesHistory = append(pagesHistory, currentPage)
	}

	show(p)
}

func show(p Page) {
	currentPage = p
	RefreshUI()
}

func RefreshUI() {
	ui.Clear()
	ui.Render(currentPage.MainUI(), helpBarWidget)
}

func GoBack() {
	if len(pagesHistory) == 0 {
		return
	}
	previousPage := pagesHistory[len(pagesHistory)-1]
	pagesHistory = pagesHistory[:len(pagesHistory)-1]
	show(previousPage)
}
func HandleKeyEvent(e *ui.Event) (Command, error) {
	switch e.ID {
	case "q", "<C-c>":
		ui.Close()
		return Exit, nil

	case "<Escape>", "<C-<Backspace>>", "<Backspace>":
		GoBack()
		RefreshUI()
	case "<Resize>":
		payload := e.Payload.(ui.Resize)
		helpBarWidget.SetRect(0, payload.Height-1, payload.Width, payload.Height)
		currentPage.MainUI().SetRect(0, 0, payload.Width, payload.Height-1)
		RefreshUI()
	default:
		cmd, err := currentPage.HandleEvent(e)
		if err != nil {
			return cmd, err
		}
		cmd, err = audioPlayerWidget.HandleEvent(e)
		RefreshUI()
		return cmd, err
	}
	return Nothing, nil

}

func newHelpBarWidget() *widgets.Paragraph {
	helpBarWidget := widgets.NewParagraph()
	helpBarWidget.Text = "[ Enter ](fg:black)[Select](fg:black,bg:green) " +
		"[ p, Space ](fg:black)[Play/Pause](fg:black,bg:green) " +
		"[Esc ](fg:black)[Back](fg:black,bg:green) " +
		"[Right ](fg:black)[+10s](fg:black,bg:green) " +
		"[Left ](fg:black)[-10s](fg:black,bg:green) " +
		"[ q ](fg:black)[Exit](fg:black,bg:green)"
	helpBarWidget.Border = false
	helpBarWidget.WrapText = true
	helpBarWidget.TextStyle = ui.Style{Modifier: ui.ModifierBold, Bg: ui.ColorWhite}
	termWidth, termHeight := ui.TerminalDimensions()
	helpBarWidget.SetRect(0, termHeight-1, termWidth, termHeight)
	return helpBarWidget
}

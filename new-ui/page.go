package newui

import (
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

var (
	FgColor       = ui.ColorWhite
	AccentColor   = ui.ColorMagenta
	currentPage   *Page
	helpBarWidget *widgets.Paragraph
)

type Page interface {
	MainUI() *ui.Grid
	HandleEvent(*ui.Event)
}

func Show(p Page) {
	if helpBarWidget == nil {
		helpBarWidget = newHelpBarWidget()
	}
	ui.Clear()
	ui.Render(p.MainUI(), helpBarWidget)
	currentPage = &p
}
func HandleKeyEvent(e *ui.Event) Command {
	switch e.ID {
	case "q", "<C-c>":
		ui.Close()
		return Exit
	default:
		(*currentPage).HandleEvent(e)
		ui.Clear()
		ui.Render((*currentPage).MainUI(), helpBarWidget)
		return Nothing
	}

}
func newHelpBarWidget() *widgets.Paragraph {
	helpBarWidget := widgets.NewParagraph()
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
	termWidth, termHeight := ui.TerminalDimensions()
	helpBarWidget.SetRect(0, termHeight-1, termWidth, termHeight)
	return helpBarWidget
}

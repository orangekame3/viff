package pkg

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func BuildAppView() *tview.Application {
	return tview.NewApplication()
}

func BuildSidePane(text, title string) *tview.TextView {
	pane := tview.NewTextView().SetText(text)
	pane.SetDynamicColors(true)
	pane.SetTitle(fmt.Sprintf("File: %s", title)).SetBorder(true)
	return pane
}

func BuildInlinePane(text, title string) *tview.TextView {
	pane := tview.NewTextView().SetText(text)
	pane.SetDynamicColors(true)
	pane.SetTitle("Inline View").SetBorder(true)
	return pane
}

func BuildSideBySidePane(leftPane, rightPane *tview.TextView) *tview.Flex {
	flex := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(leftPane, 0, 1, false).
		AddItem(rightPane, 0, 1, false)
	flex.SetTitle("Side-by-Side View").SetBorder(true)
	return flex
}

func BuildPages(sideBySide *tview.Flex, inline *tview.TextView) *tview.Pages {
	pages := tview.NewPages()
	pages.AddPage("sydeBySide", sideBySide, true, true)
	pages.AddPage("inline", inline, true, false)
	return pages
}

func BuildHelpPane() *tview.Flex {
	help := tview.NewTextView().SetText("【help 】 <Space> change [side-by-side ⇄ inline], <Esc,q> quit viff, <h> hide help")
	legend := tview.NewTextView().SetText("【legend】 red → file1, green → file2")
	repo := tview.NewTextView().SetText("【check update】 → https://github.com/orangekame3/viff")

	textPane := tview.NewFlex().
		AddItem(help, 0, 1, false).
		AddItem(legend, 0, 1, false).
		AddItem(repo, 0, 1, false)
	return textPane
}

func BuildMainView(pages *tview.Pages, help *tview.Flex) *tview.Flex {
	main := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(pages, 0, 1, false).
		AddItem(help, 1, 0, false)
	main.SetTitle("viff").SetBorder(true)
	return main
}

func IsQuitKey(e *tcell.EventKey) bool {
	return e.Key() == tcell.KeyEscape || e.Rune() == 'q'
}

func IsChangeModeKey(e *tcell.EventKey) bool {
	return e.Rune() == ' '
}

func IsHelpKey(e *tcell.EventKey) bool {
	return e.Rune() == 'h'
}

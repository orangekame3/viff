// Package pkg provides the core functionality of the program.
package pkg

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// BuildSidePane returns a new side pane.
func BuildSidePane(text, title string) *tview.TextView {
	pane := tview.NewTextView().SetText(text)
	pane.SetDynamicColors(true)
	pane.SetTitle(fmt.Sprintf("File: %s", title)).SetBorder(true)
	return pane
}

// BuildInlinePane returns a new inline pane.
func BuildInlinePane(text, title string) *tview.TextView {
	pane := tview.NewTextView().SetText(text)
	pane.SetDynamicColors(true)
	pane.SetTitle("Inline View").SetBorder(true)
	return pane
}

// BuildSplitPane returns a new split pane.
func BuildSplitPane(leftPane, rightPane *tview.TextView) *tview.Flex {
	flex := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(leftPane, 0, 1, false).
		AddItem(rightPane, 0, 1, false)
	flex.SetTitle("Split View").SetBorder(true)
	return flex
}

// BuildPages returns a new pages.
func BuildPages(split *tview.Flex, inline *tview.TextView) *tview.Pages {
	pages := tview.NewPages()
	pages.AddPage("split", split, true, true)
	pages.AddPage("inline", inline, true, false)
	return pages
}

// BuildHelpPane returns a new help pane.
func BuildHelpPane() *tview.Flex {
	help := tview.NewTextView().SetText("[esc/q] quit, [space] change mode, [i] hide this info, [h] focus left, [l] focus right, [j] scroll down, [k] scroll up")
	textPane := tview.NewFlex().
		AddItem(help, 0, 1, false)
	return textPane
}

// BuildMainView returns a new main view.
func BuildMainView(pages *tview.Pages, help *tview.Flex) *tview.Flex {
	main := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(pages, 0, 1, false).
		AddItem(help, 1, 0, false)
	main.SetTitle("viff").SetBorder(true)
	return main
}

// IsQuitKey returns true if the key is a quit key.
func IsQuitKey(e *tcell.EventKey) bool {
	return e.Key() == tcell.KeyEscape || e.Rune() == 'q'
}

// IsChangeModeKey returns true if the key is a change mode key.
func IsChangeModeKey(e *tcell.EventKey) bool {
	return e.Rune() == ' '
}

// IsHelpKey returns true if the key is a help key.
func IsHelpKey(e *tcell.EventKey) bool {
	return e.Rune() == 'i'
}

// IsFocusLeftKey returns true if the key is a focus left key.
func IsFocusLeftKey(e *tcell.EventKey) bool {
	return e.Rune() == 'h'
}

// IsFocusRightKey returns true if the key is a focus right key.
func IsFocusRightKey(e *tcell.EventKey) bool {
	return e.Rune() == 'l'
}

// IsScrollDownKey returns true if the key is a scroll down key.
func IsScrollDownKey(e *tcell.EventKey) bool {
	return e.Rune() == 'j'
}

// IsScrollUpKey returns true if the key is a scroll up key.
func IsScrollUpKey(e *tcell.EventKey) bool {
	return e.Rune() == 'k'
}

// ScrollDown scrolls down the focused pane.
func ScrollDown(app *tview.Application) *tcell.EventKey {
	focusedPane, _ := app.GetFocus().(*tview.TextView)
	if focusedPane != nil {
		row, col := focusedPane.GetScrollOffset()
		focusedPane.ScrollTo(row+1, col)
		return nil
	}
	return nil
}

// ScrollUp scrolls up the focused pane.
func ScrollUp(app *tview.Application) *tcell.EventKey {
	focusedPane, _ := app.GetFocus().(*tview.TextView)
	if focusedPane != nil {
		row, col := focusedPane.GetScrollOffset()
		focusedPane.ScrollTo(row-1, col)
		return nil
	}
	return nil
}

// Package pkg provides the core functionality of the program.
package pkg

import (
	"fmt"

	catppuccin "github.com/catppuccin/go"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/spf13/viper"
)

// BuildSidePane returns a new side pane.
func BuildSidePane(text, title string) *tview.TextView {
	pane := tview.NewTextView().SetText(text)
	pane.SetDynamicColors(true)
	pane.SetTitle(fmt.Sprintf("File: %s", title)).SetTitleColor(tcell.GetColor(Flavour.Text().Hex)).SetBorder(true).SetBorderColor(tcell.GetColor(Flavour.Text().Hex))
	pane.SetBackgroundColor(tcell.GetColor(Flavour.Surface1().Hex))
	return pane
}

// BuildInlinePane returns a new inline pane.
func BuildInlinePane(text, title string) *tview.TextView {
	pane := tview.NewTextView().SetText(text)
	pane.SetDynamicColors(true)
	pane.SetTitle("Inline View").SetTitleColor(tcell.GetColor(Flavour.Text().Hex)).SetBorder(true).SetBorderColor(tcell.GetColor(Flavour.Text().Hex))
	pane.SetBackgroundColor(tcell.GetColor(Flavour.Surface1().Hex))
	return pane
}

// BuildSplitPane returns a new split pane.
func BuildSplitPane(leftPane, rightPane *tview.TextView) *tview.Flex {
	flex := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(leftPane, 0, 1, false).
		AddItem(rightPane, 0, 1, false)
	flex.SetTitle("Split View").SetTitleColor(tcell.GetColor(Flavour.Text().Hex)).SetBorder(true).SetBorderColor(tcell.GetColor(Flavour.Text().Hex))
	flex.SetBackgroundColor(tcell.GetColor(Flavour.Surface1().Hex))
	return flex
}

// BuildPages returns a new pages.
func BuildPages(split *tview.Flex, inline *tview.TextView) *tview.Pages {
	pages := tview.NewPages()
	pages.AddPage("split", split, true, true)
	pages.AddPage("inline", inline, true, false)
	pages.SetBackgroundColor(tcell.GetColor(Flavour.Surface1().Hex))
	return pages
}

// BuildHelpPane returns a new help pane.
func BuildHelpPane() *tview.Flex {
	help := tview.NewTextView().SetText("[esc/q] quit, [space] change mode, [i] hide this info, [h] focus left, [l] focus right, [j] scroll down, [k] scroll up").SetTextColor(tcell.GetColor(Flavour.Text().Hex))
	help.SetBackgroundColor(tcell.GetColor(Flavour.Surface1().Hex))
	textPane := tview.NewFlex().AddItem(help, 0, 1, false)
	textPane.SetBackgroundColor(tcell.GetColor(Flavour.Surface1().Hex))
	return textPane
}

// BuildMainView returns a new main view.
func BuildMainView(pages *tview.Pages, help *tview.Flex) *tview.Flex {
	main := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(pages, 0, 1, false).
		AddItem(help, 1, 0, false)
	main.SetTitle("viff").SetTitleColor(tcell.GetColor(Flavour.Text().Hex)).SetBorder(true).SetBorderColor(tcell.GetColor(Flavour.Text().Hex))
	main.SetBackgroundColor(tcell.GetColor(Flavour.Surface1().Hex))

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

var Flavour = catppuccin.Mocha

func init() {
	// viperを設定してconfig.tomlを読み込む

	viper.SetConfigName("config")      // config.toml
	viper.AddConfigPath("$HOME/.viff") // ユーザーのホームディレクトリ内の.viffディレクトリを追加
	err := viper.ReadInConfig()        // 設定ファイルを読み込む
	if err != nil {
		Flavour = catppuccin.Macchiato
	} else {
		theme := viper.GetString("theme") // 'theme'キーでテーマを取得
		Flavour = themeMap[theme]
	}
}

var themeMap = map[string]catppuccin.Flavour{
	"Macchiato": catppuccin.Macchiato,
	"Latte":     catppuccin.Latte,
}

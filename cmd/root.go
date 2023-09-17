/*
Copyright © 2023 Takafumi Miyanaga <miya.org.0309@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/filepicker"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gdamore/tcell/v2"
	"github.com/orangekame3/viff/pkg"
	"github.com/rivo/tview"
	"github.com/shibukawa/cdiff"
	"github.com/spf13/cobra"
)

type model struct {
	filepicker   filepicker.Model
	selectedFile string
	quitting     bool
	err          error
}

func (m *model) Init() tea.Cmd {
	return m.filepicker.Init()
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.filepicker, cmd = m.filepicker.Update(msg)

	// Did the user select a file?
	if didSelect, path := m.filepicker.DidSelectFile(msg); didSelect {
		// Get the path of the selected file.
		m.selectedFile = path
		m.quitting = true
		return m, tea.Quit
	}

	return m, cmd
}

func (m *model) View() string {
	if m.quitting {
		return ""
	}
	var s string
	if m.err != nil {
		s = m.filepicker.Styles.DisabledFile.Render(m.err.Error())
	} else if m.selectedFile == "" {
		s = "Pick a file:"
	} else {
		s = "Selected file: " + m.filepicker.Styles.Selected.Render(m.selectedFile)
	}
	return s + "\n\n" + m.filepicker.View() + "\n"
}

var rootCmd = &cobra.Command{
	Use:   "viff",
	Short: "A tool to display two files side by side in the terminal",
	Long:  `viff is a CLI tool that takes two file paths as arguments and displays the contents side by side in the terminal.`,
	Run: func(cmd *cobra.Command, args []string) {
		fp1 := filepicker.New()
		fp1.AllowedTypes = []string{".txt", ".go", ".md"}
		fp1.CurrentDirectory, _ = os.Getwd()

		m1 := model{
			filepicker: fp1,
		}

		p1 := tea.NewProgram(&m1)

		// 第1段階のfilepickerの実行
		if _, err := p1.Run(); err != nil {
			fmt.Printf("Error: %v", err)
			return
		}

		if m1.selectedFile == "" {
			fmt.Println("No file was selected.")
			return
		}

		// 第2段階のfilepickerの設定
		fp2 := filepicker.New()
		fp2.AllowedTypes = []string{".txt", ".go", ".md"}
		fp2.CurrentDirectory, _ = os.Getwd()

		m2 := model{
			filepicker: fp2,
		}

		p2 := tea.NewProgram(&m2)

		// 第2段階のfilepickerの実行
		if _, err := p2.Run(); err != nil {
			fmt.Printf("Error: %v", err)
			return
		}

		if m2.selectedFile == "" {
			fmt.Println("No file was selected.")
			return
		}
		file1 := m1.selectedFile
		file2 := m2.selectedFile

		// 選択したファイルを読み込む
		oldContent, err := os.ReadFile(file1)
		if err != nil {
			fmt.Println("failed to read file1: ", err)
			return
		}

		newContent, err := os.ReadFile(file2)
		if err != nil {
			fmt.Println("failed to read file2: ", err)
			return
		}

		diff := cdiff.Diff(string(oldContent), string(newContent), cdiff.LineByLine)
		oldText, newText := pkg.GenStringForSplit(diff)
		inlineText := pkg.GenStringForInline(diff)

		// Build View
		left := pkg.BuildSidePane(oldText, file1)
		right := pkg.BuildSidePane(newText, file2)
		split := pkg.BuildSplitPane(left, right)
		inline := pkg.BuildInlinePane(inlineText, "inline")
		pages := pkg.BuildPages(split, inline)
		help := pkg.BuildHelpPane()
		main := pkg.BuildMainView(pages, help)

		// Build App
		isInline := false
		isVisible := true
		app := tview.NewApplication()
		app.SetInputCapture(func(e *tcell.EventKey) *tcell.EventKey {
			if pkg.IsQuitKey(e) {
				app.Stop()
				return nil
			}
			// Change Mode
			if pkg.IsChangeModeKey(e) && isInline {
				pages.SwitchToPage("split")
				app.SetFocus(left)
				isInline = false
				return e
			}
			if pkg.IsChangeModeKey(e) && !isInline {
				pages.SwitchToPage("inline")
				app.SetFocus(inline)
				isInline = true

				return e
			}
			// Toggle Help
			if pkg.IsHelpKey(e) && isVisible {
				main.RemoveItem(help)
				isVisible = false
				return e
			}
			if pkg.IsHelpKey(e) && !isVisible {
				main.AddItem(help, 1, 0, false)
				isVisible = true
				return e
			}
			// Focus
			if pkg.IsFocusLeftKey(e) && !isInline {
				app.SetFocus(left)
				return nil
			}
			if pkg.IsFocusRightKey(e) && !isInline {
				app.SetFocus(right)
				return nil
			}
			// Scroll
			if pkg.IsScrollDownKey(e) {
				return pkg.ScrollDown(app)
			}
			if pkg.IsScrollUpKey(e) {
				return pkg.ScrollUp(app)
			}
			return e
		})
		if err := app.SetRoot(main, true).EnableMouse(true).Run(); err != nil {
			panic(err)
		}
	},
}

// Execute executes rootCmd
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
}

// SetVersionInfo sets version and date to rootCmd
func SetVersionInfo(version, date string) {
	rootCmd.Version = fmt.Sprintf("%s (Built on %s)", version, date)
}

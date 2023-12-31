//Package cmd provides the command line interface of the program.
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
	"path/filepath"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/kylelemons/godebug/diff"
	"github.com/orangekame3/diffy"
	"github.com/orangekame3/viff/pkg"
	"github.com/rivo/tview"
	"github.com/spf13/cobra"
)

// compareCmd represents the compare command
var compareCmd = &cobra.Command{
	Use:   "compare [file1] [file2]]",
	Short: "compaare two files",
	Long:  `viff is a CLI tool that takes two file paths as arguments and displays the contents side by side in the terminal.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			fmt.Println("requires two file paths as arguments")
			return
		}
		file1 := args[0]
		file2 := args[1]

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

		// Generate Diff
		aLines := strings.Split(string(oldContent), "\n")
		bLines := strings.Split(string(newContent), "\n")
		chunks := diff.DiffChunks(aLines, bLines)
		lines := diffy.FormatWithLineNumber(chunks)
		oldText, newText := pkg.GenStringForSplit(lines)
		inlineText := pkg.GenStringForInline(lines)

		// Build View
		left := pkg.BuildSidePane(oldText, filepath.Base(file1))
		right := pkg.BuildSidePane(newText, filepath.Base(file2))
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

func init() {
	rootCmd.AddCommand(compareCmd)
}

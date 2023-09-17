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
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/orangekame3/viff/pkg"
	"github.com/rivo/tview"
	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/spf13/cobra"
)

var charDiff bool

var rootCmd = &cobra.Command{
	Use:   "viff",
	Short: "A tool to display two files side by side in the terminal",
	Long:  `viff is a CLI tool that takes two file paths as arguments and displays the contents side by side in the terminal.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			fmt.Println("requires two file paths as arguments")
			return
		}

		file1Content, err := os.ReadFile(args[0])
		if err != nil {
			fmt.Println("failed to read file1: ", err)
			return
		}

		file2Content, err := os.ReadFile(args[1])
		if err != nil {
			fmt.Println("failed to read file2: ", err)
			return
		}

		dmp := diffmatchpatch.New()
		file1DiffContent := []string{}
		file2DiffContent := []string{}

		var diffs []diffmatchpatch.Diff
		if charDiff {
			diffs = dmp.DiffMain(string(file1Content), string(file2Content), false)
			for _, diff := range diffs {
				switch diff.Type {
				case diffmatchpatch.DiffEqual:
					file1DiffContent = append(file1DiffContent, diff.Text)
					file2DiffContent = append(file2DiffContent, diff.Text)
				case diffmatchpatch.DiffInsert:
					file2DiffContent = append(file2DiffContent, fmt.Sprintf("[green]%s[white]", diff.Text))
				case diffmatchpatch.DiffDelete:
					file1DiffContent = append(file1DiffContent, fmt.Sprintf("[red]%s[white]", diff.Text))
				}
			}
		} else {
			a, b, c := dmp.DiffLinesToChars(string(file1Content), string(file2Content))
			diffs = dmp.DiffMain(a, b, false)
			diffs = dmp.DiffCharsToLines(diffs, c)
			for _, diff := range diffs {
				switch diff.Type {
				case diffmatchpatch.DiffEqual:
					file1DiffContent = append(file1DiffContent, diff.Text)
					file2DiffContent = append(file2DiffContent, diff.Text)
				case diffmatchpatch.DiffInsert:
					file2DiffContent = append(file2DiffContent, fmt.Sprintf("[#000000:#00FF00:b] %s [-]", diff.Text))
				case diffmatchpatch.DiffDelete:
					file1DiffContent = append(file1DiffContent, fmt.Sprintf("[#000000:#FF0000:b] %s [-]", diff.Text))
				}
				fmt.Println(diff)
			}
		}
		file1WithLineNumbers := []string{}
		lineNumber := 1
		for _, line := range file1DiffContent {
			lines := strings.Split(strings.TrimSuffix(line, "\n"), "\n")
			fmt.Println(lines)
			for i, ln := range lines {
				fmt.Println(i, lineNumber, ln)
				if i < len(lines)-1 { // 最後の行は新しい行番号を割り当てない
					file1WithLineNumbers = append(file1WithLineNumbers, fmt.Sprintf("%d: %s\n", lineNumber, ln))
					lineNumber++
				} else {
					file1WithLineNumbers = append(file1WithLineNumbers, fmt.Sprintf("%d: %s", lineNumber, ln))
					lineNumber++
				}
			}
		}

		// file2DiffContentに行番号を付ける
		file2WithLineNumbers := []string{}
		lineNumber = 1
		for _, line := range file2DiffContent {
			lines := strings.Split(line, "\n")
			for i, ln := range lines {
				if i < len(lines)-1 {
					file2WithLineNumbers = append(file2WithLineNumbers, fmt.Sprintf("%d: %s\n", lineNumber, ln))
					lineNumber++
				} else {
					file2WithLineNumbers = append(file2WithLineNumbers, fmt.Sprintf("%d: %s", lineNumber, ln))
				}
			}
		}
		// Inline View
		var inlineContent []string
		for _, diff := range diffs {
			switch diff.Type {
			case diffmatchpatch.DiffEqual:
				lines := strings.Split(diff.Text, "\n")
				for _, line := range lines {
					if line != "" {
						inlineContent = append(inlineContent, line)
					}
				}
			case diffmatchpatch.DiffInsert:
				lines := strings.Split(diff.Text, "\n")
				for _, line := range lines {
					if line != "" {
						inlineContent = append(inlineContent, fmt.Sprintf("[#000000:#00FF00:b]%s[-]", line))
					}
				}
			case diffmatchpatch.DiffDelete:
				lines := strings.Split(diff.Text, "\n")
				for _, line := range lines {
					if line != "" {
						inlineContent = append(inlineContent, fmt.Sprintf("[#000000:#FF0000:b]%s[-]", line))
					}
				}
			}
		}
		// Build View
		left := pkg.BuildSidePane(strings.Join(file1WithLineNumbers, ""), args[0])
		right := pkg.BuildSidePane(strings.Join(file2WithLineNumbers, ""), args[1])
		sydeBySide := pkg.BuildSideBySidePane(left, right)
		inline := pkg.BuildInlinePane(strings.Join(inlineContent, ""), "inline")
		pages := pkg.BuildPages(sydeBySide, inline)
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
				pages.SwitchToPage("sydeBySide")
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
	rootCmd.Flags().BoolVar(&charDiff, "chardiff", false, "Use character-level diff")
}

// SetVersionInfo sets version and date to rootCmd
func SetVersionInfo(version, date string) {
	rootCmd.Version = fmt.Sprintf("%s (Built on %s)", version, date)
}
